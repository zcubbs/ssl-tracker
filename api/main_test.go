package api

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/internal/util"
	"github.com/zcubbs/tlz/pkg/random"
	"github.com/zcubbs/tlz/pkg/token"
	"github.com/zcubbs/tlz/worker"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		Auth: util.AuthConfig{
			TokenSymmetricKey:    random.RandomString(32),
			AccessTokenDuration:  time.Minute,
			RefreshTokenDuration: 5 * time.Minute,
		},
	}

	server, err := NewServer(store, taskDistributor, config)
	require.NoError(t, err)

	return server
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, userId uuid.UUID, duration time.Duration) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(username, userId, duration)
	require.NoError(t, err)

	bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{
			bearerToken,
		},
	}

	return metadata.NewIncomingContext(context.Background(), md)
}
