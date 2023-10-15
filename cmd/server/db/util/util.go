package util

import (
	"fmt"
	"github.com/zcubbs/tlz/cmd/server/config"
	"github.com/zcubbs/tlz/cmd/server/logger"
)

var (
	log = logger.L()
)

func getPostgresConnectionString(dbCfg config.PostgresConfig) string {
	var sslMode string
	if dbCfg.SslMode {
		sslMode = "enable"
	} else {
		sslMode = "disable"
	}
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Username,
		dbCfg.Password,
		dbCfg.DbName,
		sslMode,
	)

	return dsn
}
