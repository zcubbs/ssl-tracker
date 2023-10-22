package api

import (
	"context"
	"github.com/google/uuid"
	mockdb "github.com/zcubbs/tlz/cmd/server/db/mock"
	db "github.com/zcubbs/tlz/cmd/server/db/sqlc"
	mockwk "github.com/zcubbs/tlz/cmd/server/worker/mock"
	"github.com/zcubbs/x/random"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizeUser_Success(t *testing.T) {
	// Setup
	storeCtrl := gomock.NewController(t)
	defer storeCtrl.Finish()

	store := mockdb.NewMockStore(storeCtrl)
	worker := mockwk.NewMockTaskDistributor(storeCtrl)

	s := newTestServer(t, store, worker)

	username := random.String(32)
	userId := uuid.UUID{}
	duration := time.Minute

	createToken, _, err := s.tokenMaker.CreateToken(username, userId, duration)
	assert.Nil(t, err)
	assert.NotEmpty(t, createToken)

	md := metadata.New(map[string]string{
		string(authorizationHeader): "Bearer " + createToken,
	})

	ctx := metadata.NewIncomingContext(context.Background(), md)

	user := db.User{
		ID:       uuid.New(),
		Username: "test",
		Role:     "ROLE_USER",
	}

	store.EXPECT().
		GetUserByUsername(gomock.Any(), gomock.Any()).
		Times(1).
		Return(user, nil)

	// Call the function
	payload, err := s.requireUser(ctx)

	// Assertions
	assert.Nil(t, err)
	assert.NotNil(t, payload)
}
