package config

const (
	ViperConfigName      = "config"
	ViperConfigType      = "yaml"
	ViperConfigEnvPrefix = "LU"
	DefaultDbName        = "tlz"
	Localhost            = "127.0.0.1"
	HttpPort             = 8000
	GrpcPort             = 9000
)

var (
	viperConfigPaths = [...]string{"./config"}

	defaults = map[string]interface{}{
		"debug":                           false,
		"http_server.port":                HttpPort,
		"http_server.allow_origins":       "http://localhost:8000,http://localhost:5173,http://127.0.0.1:5173,http://127.0.0.1:8000",
		"http_server.allow_headers":       "Origin, Content-Type, Accept",
		"http_server.tz":                  "UTC",
		"http_server.enable_print_routes": false,
		"http_server.read_header_timeout": "3s",
		"grpc_server.port":                GrpcPort,
		"grpc_server.enable_reflection":   true,
		"grpc_server.tls.enabled":         false,
		"grpc_server.tls.cert":            "",
		"grpc_server.tls.key":             "",
		"auth.token_symmetric_key":        "12345678901234567890123456789012",
		"auth.access_token_duration":      "30s",
		"auth.refresh_token_duration":     "15m",
		"database.auto_migration":         true,
		"database.postgres.enabled":       true,
		"database.postgres.host":          Localhost,
		"database.postgres.port":          5432,
		"database.postgres.username":      "postgres",
		"database.postgres.password":      "postgres",
		"database.postgres.db_name":       DefaultDbName,
		"database.postgres.ssl_mode":      false,
		"database.postgres.verbose":       false,
		"database.postgres.cert_pem":      "",
		"database.postgres.cert_key":      "",
		"database.postgres.max_conns":     10,
		"database.postgres.min_conns":     4,
	}

	allowedEnvVarKeys = []string{
		"debug",
		"init_admin_password",
		"http_server.port",
		"http_server.allow_origins",
		"http_server.allow_headers",
		"http_server.tz",
		"http_server.enable_print_routes",
		"http_server.read_header_timeout",
		"grpc_server.port",
		"grpc_server.enable_reflection",
		"grpc_server.tls.enabled",
		"grpc_server.tls.cert",
		"grpc_server.tls.key",
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
		"database.postgres.cert_pem",
		"database.postgres.cert_key",
		"database.postgres.max_conns",
		"database.postgres.min_conns",
	}
)
