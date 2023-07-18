package util

const (
	ViperConfigName      = "config"
	ViperConfigType      = "yaml"
	ViperConfigEnvPrefix = "TLZ"
	DefaultDbName        = "database"
)

var (
	viperConfigPaths = [...]string{"./config"}

	defaults = map[string]interface{}{
		"debug":                     false,
		"http_server.port":          8000,
		"database.postgres.enabled": false,
		"database.sqlite.enabled":   true,
		"database.sqlite.db_name":   DefaultDbName,
	}

	allowedEnvVarKeys = []string{
		"debug",
		"http_server.port",
		"database.postgres.enabled",
		"database.postgres.host",
		"database.postgres.port",
		"database.postgres.username",
		"database.postgres.password",
		"database.postgres.database",
		"database.postgres.ssl_mode",
		"database.postgres.verbose",
		"database.sqlite.enabled",
		"database.sqlite.db_name",
	}
)
