package db

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/zcubbs/tlz/db/migrations"
	"github.com/zcubbs/tlz/internal/logger"
	"github.com/zcubbs/tlz/internal/util"
	"os"
	"testing"
)

var testStore Store

func TestMain(m *testing.M) {
	config := util.Bootstrap()
	ctx := context.Background()
	// Migrate database
	err := migrations.Run(config.Database)
	if err != nil {
		log.Fatal("failed perform database migrations", "error", err)
	}
	conn, err := util.DbConnect(ctx, config.Database, logger.GetLogger())
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	testStore = NewSQLStore(conn)

	os.Exit(m.Run())
}
