package api

//
//import (
//	"bytes"
//	"database/sql"
//	"encoding/json"
//	"fmt"
//	"github.com/gofiber/fiber/v2"
//	"github.com/zcubbs/tlz/util"
//	"go.uber.org/mock/gomock"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"reflect"
//	"testing"
//
//	"github.com/stretchr/testify/require"
//	mockdb "github.com/zcubbs/tlz/db/mock"
//	db "github.com/zcubbs/tlz/db/sqlc"
//)
//
//type eqCreateUserParamsMatcher struct {
//	arg      db.CreateUserParams
//	password string
//}
//
//func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
//	arg, ok := x.(db.CreateUserParams)
//	if !ok {
//		return false
//	}
//
//	err := util.CheckPassword(e.password, arg.HashedPassword)
//	if err != nil {
//		return false
//	}
//
//	e.arg.HashedPassword = arg.HashedPassword
//	return reflect.DeepEqual(e.arg, arg)
//}
//
//func (e eqCreateUserParamsMatcher) String() string {
//	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
//}
//
//func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
//	return eqCreateUserParamsMatcher{arg, password}
//}
//
//func TestCreateUserAPI(t *testing.T) {
//	user, password := randomUser(t)
//
//	testCases := []struct {
//		name          string
//		body          fiber.Map
//		buildStubs    func(store *mockdb.MockStore)
//		checkResponse func(recoder *httptest.ResponseRecorder)
//	}{
//		{
//			name: "OK",
//			body: fiber.Map{
//				"username":  user.Username,
//				"password":  password,
//				"full_name": user.FullName,
//				"email":     user.Email,
//			},
//			buildStubs: func(store *mockdb.MockStore) {
//				arg := db.CreateUserParams{
//					Username: user.Username,
//					FullName: user.FullName,
//					Email:    user.Email,
//				}
//				store.EXPECT().
//					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
//					Times(1).
//					Return(user, nil)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusOK, recorder.Code)
//				requireBodyMatchUser(t, recorder.Body, user)
//			},
//		},
//		{
//			name: "InternalError",
//			body: fiber.Map{
//				"username":  user.Username,
//				"password":  password,
//				"full_name": user.FullName,
//				"email":     user.Email,
//			},
//			buildStubs: func(store *mockdb.MockStore) {
//				store.EXPECT().
//					CreateUser(gomock.Any(), gomock.Any()).
//					Times(1).
//					Return(db.User{}, sql.ErrConnDone)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusInternalServerError, recorder.Code)
//			},
//		},
//		{
//			name: "DuplicateUsername",
//			body: fiber.Map{
//				"username":  user.Username,
//				"password":  password,
//				"full_name": user.FullName,
//				"email":     user.Email,
//			},
//			buildStubs: func(store *mockdb.MockStore) {
//				store.EXPECT().
//					CreateUser(gomock.Any(), gomock.Any()).
//					Times(1).
//					Return(db.User{}, db.ErrUniqueViolation)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusForbidden, recorder.Code)
//			},
//		},
//		{
//			name: "InvalidUsername",
//			body: fiber.Map{
//				"username":  "invalid-user#1",
//				"password":  password,
//				"full_name": user.FullName,
//				"email":     user.Email,
//			},
//			buildStubs: func(store *mockdb.MockStore) {
//				store.EXPECT().
//					CreateUser(gomock.Any(), gomock.Any()).
//					Times(0)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusBadRequest, recorder.Code)
//			},
//		},
//		{
//			name: "InvalidEmail",
//			body: fiber.Map{
//				"username":  user.Username,
//				"password":  password,
//				"full_name": user.FullName,
//				"email":     "invalid-email",
//			},
//			buildStubs: func(store *mockdb.MockStore) {
//				store.EXPECT().
//					CreateUser(gomock.Any(), gomock.Any()).
//					Times(0)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusBadRequest, recorder.Code)
//			},
//		},
//		{
//			name: "TooShortPassword",
//			body: fiber.Map{
//				"username":  user.Username,
//				"password":  "123",
//				"full_name": user.FullName,
//				"email":     user.Email,
//			},
//			buildStubs: func(store *mockdb.MockStore) {
//				store.EXPECT().
//					CreateUser(gomock.Any(), gomock.Any()).
//					Times(0)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusBadRequest, recorder.Code)
//			},
//		},
//	}
//
//	ctrl := gomock.NewController(t)
//	store := mockdb.NewMockStore(ctrl)
//	server := newTestServer(t, store)
//
//	app := fiber.New()
//	path := "/api/users"
//	// Define route.
//	app.Get(path, func(c *fiber.Ctx) error {
//		return c.Status(fiber.StatusOK).SendString("OK!")
//	})
//	for i := range testCases {
//		tc := testCases[i]
//		tc.buildStubs(store)
//
//		t.Run(tc.name, func(t *testing.T) {
//
//			req := httptest.NewRequest(fiber.MethodGet, path, nil)
//
//			tc.setupAuth(t, req, server.tokenMaker)
//
//			resp, _ := app.Test(req, -1)
//			tc.checkResponse(t, resp)
//		})
//	}
//
//	for i := range testCases {
//		tc := testCases[i]
//
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			store := mockdb.NewMockStore(ctrl)
//			tc.buildStubs(store)
//
//			server := newTestServer(t, store)
//			recorder := httptest.NewRecorder()
//
//			// Marshal body data to JSON
//			data, err := json.Marshal(tc.body)
//			require.NoError(t, err)
//
//			url := "/users"
//			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
//			require.NoError(t, err)
//
//			server.router.ServeHTTP(recorder, request)
//			tc.checkResponse(recorder)
//		})
//	}
//}
//
//func TestLoginUserAPI(t *testing.T) {
//	user, password := randomUser(t)
//
//	testCases := []struct {
//		name          string
//		body          gin.H
//		buildStubs    func(store *mockdb.Store)
//		checkResponse func(recoder *httptest.ResponseRecorder)
//	}{
//		{
//			name: "OK",
//			body: fiber.Map{
//				"username": user.Username,
//				"password": password,
//			},
//			buildStubs: func(store *mockdb.Store) {
//				store.EXPECT().
//					GetUser(gomock.Any(), gomock.Eq(user.Username)).
//					Times(1).
//					Return(user, nil)
//				store.EXPECT().
//					CreateSession(gomock.Any(), gomock.Any()).
//					Times(1)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusOK, recorder.Code)
//			},
//		},
//		{
//			name: "UserNotFound",
//			body: fiber.Map{
//				"username": "NotFound",
//				"password": password,
//			},
//			buildStubs: func(store *mockdb.Store) {
//				store.EXPECT().
//					GetUser(gomock.Any(), gomock.Any()).
//					Times(1).
//					Return(db.User{}, db.ErrRecordNotFound)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusNotFound, recorder.Code)
//			},
//		},
//		{
//			name: "IncorrectPassword",
//			body: fiber.Map{
//				"username": user.Username,
//				"password": "incorrect",
//			},
//			buildStubs: func(store *mockdb.Store) {
//				store.EXPECT().
//					GetUser(gomock.Any(), gomock.Eq(user.Username)).
//					Times(1).
//					Return(user, nil)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusUnauthorized, recorder.Code)
//			},
//		},
//		{
//			name: "InternalError",
//			body: fiber.Map{
//				"username": user.Username,
//				"password": password,
//			},
//			buildStubs: func(store *mockdb.Store) {
//				store.EXPECT().
//					GetUser(gomock.Any(), gomock.Any()).
//					Times(1).
//					Return(db.User{}, sql.ErrConnDone)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusInternalServerError, recorder.Code)
//			},
//		},
//		{
//			name: "InvalidUsername",
//			body: fiber.Map{
//				"username": "invalid-user#1",
//				"password": password,
//			},
//			buildStubs: func(store *mockdb.Store) {
//				store.EXPECT().
//					GetUser(gomock.Any(), gomock.Any()).
//					Times(0)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusBadRequest, recorder.Code)
//			},
//		},
//	}
//
//	for i := range testCases {
//		tc := testCases[i]
//
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			store := mockdb.NewMockStore(ctrl)
//			tc.buildStubs(store)
//
//			server := newTestServer(t, store)
//			recorder := httptest.NewRecorder()
//
//			// Marshal body data to JSON
//			data, err := json.Marshal(tc.body)
//			require.NoError(t, err)
//
//			url := "/users/login"
//			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
//			require.NoError(t, err)
//
//			server.router.ServeHTTP(recorder, request)
//			tc.checkResponse(recorder)
//		})
//	}
//}
//
//func randomUser(t *testing.T) (user db.User, password string) {
//	password = util.RandomString(6)
//	hashedPassword, err := util.HashPassword(password)
//	require.NoError(t, err)
//
//	user = db.User{
//		Username:       util.RandomOwner(),
//		HashedPassword: hashedPassword,
//		FullName:       util.RandomOwner(),
//		Email:          util.RandomEmail(),
//	}
//	return
//}
//
//func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
//	data, err := io.ReadAll(body)
//	require.NoError(t, err)
//
//	var gotUser db.User
//	err = json.Unmarshal(data, &gotUser)
//
//	require.NoError(t, err)
//	require.Equal(t, user.Username, gotUser.Username)
//	require.Equal(t, user.FullName, gotUser.FullName)
//	require.Equal(t, user.Email, gotUser.Email)
//	require.Empty(t, gotUser.HashedPassword)
//}
