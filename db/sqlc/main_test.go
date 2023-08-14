package db

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/zcubbs/tlz/internal/util"
	"os"
	"testing"
)

var testStore Store

func TestMain(m *testing.M) {
	config := util.Bootstrap()
	ctx := context.Background()
	conn, err := util.DbConnect(ctx, config.Database)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	testStore = NewSQLStore(conn)

	os.Exit(m.Run())
}
