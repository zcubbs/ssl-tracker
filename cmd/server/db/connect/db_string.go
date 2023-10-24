package connect

import (
	"fmt"
	"github.com/zcubbs/ssl-tracker/cmd/server/config"
	"github.com/zcubbs/ssl-tracker/cmd/server/logger"
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
