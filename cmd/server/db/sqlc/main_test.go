package db

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/zcubbs/tlz/cmd/server/config"
	"github.com/zcubbs/tlz/cmd/server/db/migration"
	"github.com/zcubbs/tlz/cmd/server/db/util"
	"os"
	"testing"
)

var testStore Store

func TestMain(m *testing.M) {
	cfg := config.Bootstrap()
	ctx := context.Background()
	// Migrate database
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
