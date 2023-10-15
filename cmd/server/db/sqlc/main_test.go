package db

import (
	"context"
	"github.com/zcubbs/tlz/cmd/server/config"
	"github.com/zcubbs/tlz/cmd/server/db/migration"
	"github.com/zcubbs/tlz/cmd/server/db/util"
	"github.com/zcubbs/tlz/cmd/server/logger"
	"os"
	"testing"
)

var testStore Store

func TestMain(m *testing.M) {
	cfg := config.Bootstrap()
	ctx := context.Background()
	// Migrate database
	var (
		log = logger.L()
	)
	err := migration.Run(cfg.Database)
	if err != nil {
		log.Fatal("failed perform database migrations", "error", err)
	}
	conn, err := util.Connect(ctx, cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	testStore = NewSQLStore(conn)

	os.Exit(m.Run())
}
