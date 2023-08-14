package api

import (
	"github.com/stretchr/testify/require"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/internal/util"
	"github.com/zcubbs/tlz/pkg/random"
	"os"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		Auth: util.AuthConfig{
			TokenSymmetricKey:    random.RandomString(32),
			AccessTokenDuration:  time.Minute,
			RefreshTokenDuration: 5 * time.Minute,
		},
	}

	server, err := NewServer(store, nil, config)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
