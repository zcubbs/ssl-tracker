package api

//import (
//	"context"
//	"database/sql"
//	"fmt"
//	"github.com/stretchr/testify/require"
//	"github.com/zcubbs/go-pkg/password"
//	"github.com/zcubbs/go-pkg/random"
//	mockdb "github.com/zcubbs/linkup/cmd/server/db/mock"
//	db "github.com/zcubbs/linkup/cmd/server/db/sqlc"
//	pb "github.com/zcubbs/linkup/pb"
//	"go.uber.org/mock/gomock"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/status"
//	"reflect"
//	"testing"
//)
//
//type createUserParamsMatcher struct {
//	arg      db.CreateUserParams
//	password string
//	user     db.User
//}
//
//func (e createUserParamsMatcher) Matches(x interface{}) bool {
//	actualArg, ok := x.(db.CreateUserParams)
//	if !ok {
//		return false
//	}
//
//	err := password.Check(e.password, actualArg.HashedPassword)
//	if err != nil {
//		return false
//	}
//
//	e.arg.HashedPassword = actualArg.HashedPassword
//	return reflect.DeepEqual(e.arg, actualArg)
//}
//
//func (e createUserParamsMatcher) String() string {
//	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
//}
//
//func CreateUserParamsMatcher(arg db.CreateUserParams, password string, user db.User) gomock.Matcher {
//	return createUserParamsMatcher{arg, password, user}
//}
//
//func randomUser(t *testing.T) (user db.User, pwd string) {
//	pwd = random.RandomString(15)
//	hashedPassword, err := password.Hash(pwd)
//	require.NoError(t, err)
//
//	return db.User{
//		Username:       random.RandomString(15),
//		HashedPassword: hashedPassword,
//		FullName:       random.RandomString(20),
//		Email:          random.RandomEmail(),
//		Role:           pb.Role_ROLE_USER.String(),
//	}, pwd
//}
//
//func TestCreateUser(t *testing.T) {
//	user, pwd := randomUser(t)
//
//	testCases := []struct {
//		name          string
//		req           *pb.CreateUserRequest
//		setup         func(t *testing.T, store *mockdb.MockStore)
//		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
//	}{
//		{
//			name: "OK",
//			req: &pb.CreateUserRequest{
//				Username: user.Username,
//				FullName: user.FullName,
//				Email:    user.Email,
//				Password: pwd,
//				Role:     pb.Role_ROLE_ADMIN,
//			},
//			setup: func(t *testing.T, store *mockdb.MockStore) {
//				arg := db.CreateUserParams{
//					Username:       user.Username,
//					FullName:       user.FullName,
//					Email:          user.Email,
//					HashedPassword: pwd,
//					Role:           pb.Role_ROLE_ADMIN.String(),
//				}
//				store.EXPECT().
//					CreateUser(gomock.Any(), CreateUserParamsMatcher(arg, pwd, user)).
//					Times(1).
//					Return(user, nil)
//			},
//			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
//				require.NoError(t, err)
//				require.NotNil(t, res)
//				createdUser := res.GetUser()
//				require.Equal(t, user.Username, createdUser.Username)
//				require.Equal(t, user.FullName, createdUser.FullName)
//				require.Equal(t, user.Email, createdUser.Email)
//			},
//		},
//		{
//			name: "InternalError",
//			req: &pb.CreateUserRequest{
//				Username: user.Username,
//				FullName: user.FullName,
//				Email:    user.Email,
//				Password: pwd,
//				Role:     pb.Role_ROLE_ADMIN,
//			},
//			setup: func(t *testing.T, store *mockdb.MockStore) {
//				store.EXPECT().
//					CreateUser(gomock.Any(), gomock.Any()).
//					Times(1).
//					Return(db.User{}, sql.ErrConnDone)
//			},
//			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
//				require.Error(t, err)
//				st, ok := status.FromError(err)
//				require.True(t, ok)
//				require.Equal(t, st.Code(), codes.Internal)
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			storeCtrl := gomock.NewController(t)
//			defer storeCtrl.Finish()
//
//			store := mockdb.NewMockStore(storeCtrl)
//			tc.setup(t, store)
//
//			server := newTestServer(t, store)
//			res, err := server.CreateUser(context.Background(), tc.req)
//			tc.checkResponse(t, res, err)
//		})
//	}
//}
