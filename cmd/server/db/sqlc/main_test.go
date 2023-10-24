package db

import (
	"context"
	"github.com/zcubbs/ssl-tracker/cmd/server/config"
	dbConnect "github.com/zcubbs/ssl-tracker/cmd/server/db/connect"
	mig "github.com/zcubbs/ssl-tracker/cmd/server/db/migration"
	"github.com/zcubbs/ssl-tracker/cmd/server/logger"
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
	err := mig.Run(cfg.Database)
	if err != nil {
		log.Fatal("failed perform database migrations", "error", err)
	}
	conn, err := dbConnect.Connect(ctx, cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	testStore = NewSQLStore(conn)

	os.Exit(m.Run())
}
