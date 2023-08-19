package util

import "time"

type Config struct {
	Debug        bool               `mapstructure:"debug"`
	HttpServer   HttpServerConfig   `mapstructure:"http_server"`
	GrpcServer   GrpcServerConfig   `mapstructure:"grpc_server"`
	Auth         AuthConfig         `mapstructure:"auth"`
	Database     DatabaseConfig     `mapstructure:"database"`
	Redis        RedisConfig        `mapstructure:"redis"`
	Cron         CronConfig         `mapstructure:"cron"`
	Notification NotificationConfig `mapstructure:"notification"`
}

type HttpServerConfig struct {
	Port              int    `mapstructure:"port"`
	AllowOrigins      string `mapstructure:"allow_origins"`
	AllowHeaders      string `mapstructure:"allow_headers"`
	TZ                string `mapstructure:"tz"`
	EnablePrintRoutes bool   `mapstructure:"enable_print_routes"`
	// ReadHeaderTimeout is the amount of time allowed to read request headers. Default values: '3s'
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
}

type GrpcServerConfig struct {
	Port int `mapstructure:"port"`
}

type AuthConfig struct {
	TokenSymmetricKey    string        `mapstructure:"token_symmetric_key"`
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
}

type DatabaseType string

const (
	Postgres DatabaseType = "postgres"
)

type DatabaseConfig struct {
	AutoMigration bool           `mapstructure:"auto_migration" json:"auto_migration"`
	Postgres      PostgresConfig `mapstructure:"postgres" json:"postgres"`
}

func (dc *DatabaseConfig) GetDatabaseType() DatabaseType {
	if dc.Postgres.Enabled {
		return Postgres
	}
	return ""
}

type PostgresConfig struct {
	Enabled  bool   `mapstructure:"enabled" json:"enabled"`
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	DbName   string `mapstructure:"db_name" json:"db_name"`
	SslMode  bool   `mapstructure:"ssl_mode" json:"ssl_mode"`
	Verbose  bool   `mapstructure:"verbose" json:"verbose"`
	CertPem  string `mapstructure:"cert_pem"`
	CertKey  string `mapstructure:"cert_key"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type CronConfig struct {
	CheckCertificateValidity `mapstructure:"check_certificate_validity"`
	SendMailNotification     `mapstructure:"send_mail_notification"`
}

type CheckCertificateValidity struct {
	Enabled     bool   `mapstructure:"enabled"`
	CronPattern string `mapstructure:"cron_pattern"`
}

type SendMailNotification struct {
	Enabled     bool   `mapstructure:"enabled"`
	CronPattern string `mapstructure:"cron_pattern"`
}

type NotificationConfig struct {
	Mail          MailConfig `mapstructure:"mail"`
	ApiDomainName string     `mapstructure:"api_domain_name"`
}

type MailConfig struct {
	Smtp SmtpConfig `mapstructure:"smtp"`
}

type SmtpConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	FromAddress string `mapstructure:"from_address"`
	FromName    string `mapstructure:"from_name"`
}
