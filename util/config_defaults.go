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
		"debug":                                        false,
		"http_server.port":                             8000,
		"http_server.allow_origins":                    "http://localhost:8000,http://localhost:5173,http://127.0.0.1:5173,http://127.0.0.1:8000",
		"http_server.allow_headers":                    "Origin, Content-Type, Accept",
		"http_server.tz":                               "UTC",
		"http_server.enable_print_routes":              false,
		"database.postgres.enabled":                    false,
		"database.sqlite.enabled":                      true,
		"database.sqlite.db_name":                      DefaultDbName,
		"cron.check_certificate_validity.enabled":      true,
		"cron.check_certificate_validity.cron_pattern": "*/10 * * * * *",
	}

	allowedEnvVarKeys = []string{
		"debug",
		"http_server.port",
		"http_server.allow_origins",
		"http_server.allow_headers",
		"http_server.tz",
		"http_server.enable_print_routes",
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
		"cron.check_certificate_validity.enabled",
		"cron.check_certificate_validity.cron_pattern",
	}
)
