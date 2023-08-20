package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	mockdb "github.com/zcubbs/tlz/db/mock"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/pkg/password"
	"github.com/zcubbs/tlz/pkg/random"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := password.Check(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, pwd := randomUser(t)

	testCases := []struct {
		name          string
		body          createUserRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			body: createUserRequest{
				Username: user.Username,
				Password: pwd,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, pwd)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireBodyMatchUser(t, response.Body, user)
			},
		},
		{
			name: "InternalError",
			body: createUserRequest{
				Username: user.Username,
				Password: pwd,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
		{
			name: "DuplicateUsername",
			body: createUserRequest{
				Username: user.Username,
				Password: pwd,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrUniqueViolation)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusForbidden, response.StatusCode)
			},
		},
		{
			name: "InvalidUsername",
			body: createUserRequest{
				Username: "invalid-user#1",
				Password: pwd,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InvalidEmail",
			body: createUserRequest{
				Username: user.Username,
				Password: pwd,
				FullName: user.FullName,
				Email:    "invalid-email",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "TooShortPassword",
			body: createUserRequest{
				Username: user.Username,
				Password: "123",
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			app := fiber.New()
			path := "/api/users"
			app.Post(path, server.createUser)

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodPost, path, bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			tc.checkResponse(t, resp)
		})
	}
}

func TestLoginUser(t *testing.T) {
	user, pwd := randomUser(t)

	testCases := []struct {
		name          string
		body          loginUserRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			body: loginUserRequest{
				Username: user.Username,
				Password: pwd,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireBodyMatchLoginResponse(t, response.Body, user)
			},
		},
		{
			name: "UserNotFound",
			body: loginUserRequest{
				Username: "NotFound",
				Password: pwd,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "IncorrectPassword",
			body: loginUserRequest{
				Username: user.Username,
				Password: "incorrect",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InternalError",
			body: loginUserRequest{
				Username: user.Username,
				Password: pwd,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
		{
			name: "InvalidUsername",
			body: loginUserRequest{
				Username: "invalid username#",
				Password: pwd,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			app := fiber.New()
			path := "/api/users/login"
			app.Post(path, server.loginUser)

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodPost, path, bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			tc.checkResponse(t, resp)
		})
	}

}

func randomUser(t *testing.T) (user db.User, pwd string) {
	pwd = random.RandomString(6)
	hashedPassword, err := password.Hash(pwd)
	require.NoError(t, err)

	user = db.User{
		Username:       random.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       random.RandomOwner(),
		Email:          random.RandomEmail(),
	}
	return user, pwd
}

func requireBodyMatchUser(t *testing.T, body io.ReadCloser, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	requireUserMatch(t, gotUser, user)
}

func requireBodyMatchLoginResponse(t *testing.T, body io.ReadCloser, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var got loginUserResponse
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)

	gotUser := db.User{
		Username: got.User.Username,
		FullName: got.User.FullName,
		Email:    got.User.Email,
	}

	requireUserMatch(t, gotUser, user)
}

func requireUserMatch(t *testing.T, got, want db.User) {
	require.Equal(t, want.Username, got.Username)
	require.Equal(t, want.FullName, got.FullName)
	require.Equal(t, want.Email, got.Email)
	require.Empty(t, got.HashedPassword)
}
