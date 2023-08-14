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
		"grpc_server.port":                             9000,
		"auth.token_symmetric_key":                     "12345678901234567890123456789012",
		"auth.access_token_duration":                   "15m",
		"auth.refresh_token_duration":                  "24h",
		"database.auto_migration":                      false,
		"database.postgres.enabled":                    true,
		"database.postgres.host":                       "localhost",
		"database.postgres.port":                       5432,
		"database.postgres.username":                   "postgres",
		"database.postgres.password":                   "postgres",
		"database.postgres.db_name":                    "tlz",
		"database.postgres.ssl_mode":                   false,
		"database.postgres.verbose":                    false,
		"cron.check_certificate_validity.enabled":      true,
		"cron.check_certificate_validity.cron_pattern": "*/10 * * * * *",
		"cron.send_mail_notification.enabled":          true,
		"cron.send_mail_notification.cron_pattern":     "*/10 * * * * *",
		"notification.mail.smtp.enabled":               true,
		"notification.mail.smtp.host":                  "localhost",
		"notification.mail.smtp.port":                  1025,
		"notification.mail.smtp.username":              "",
		"notification.mail.smtp.password":              "",
		"notification.mail.smtp.from":                  "no-reply@tlz",
	}

	allowedEnvVarKeys = []string{
		"debug",
		"http_server.port",
		"http_server.allow_origins",
		"http_server.allow_headers",
		"http_server.tz",
		"http_server.enable_print_routes",
		"grpc_server.port",
		"auth.token_symmetric_key",
		"auth.access_token_duration",
		"auth.refresh_token_duration",
		"database.auto_migration",
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
		"cron.send_mail_notification.enabled",
		"cron.send_mail_notification.cron_pattern",
		"notification.mail.smtp.enabled",
		"notification.mail.smtp.host",
		"notification.mail.smtp.port",
		"notification.mail.smtp.username",
		"notification.mail.smtp.password",
		"notification.mail.smtp.from",
	}
)
