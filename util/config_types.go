package util

import (
	"time"
)

type Config struct {
	Debug        bool               `mapstructure:"debug"`
	HttpServer   HttpServerConfig   `mapstructure:"http_server"`
	Database     DatabaseConfig     `mapstructure:"database"`
	Cron         CronConfig         `mapstructure:"cron"`
	Notification NotificationConfig `mapstructure:"notification"`
}

type HttpServerConfig struct {
	Port                int           `mapstructure:"port"`
	AllowOrigins        string        `mapstructure:"allow_origins"`
	AllowHeaders        string        `mapstructure:"allow_headers"`
	TZ                  string        `mapstructure:"tz"`
	EnablePrintRoutes   bool          `mapstructure:"enable_print_routes"`
	TokenSymmetricKey   string        `mapstructure:"token_symmetric_key"`
	AccessTokenDuration time.Duration `mapstructure:"access_token_duration"`
}

type DatabaseType string

const (
	Postgres DatabaseType = "postgres"
	Sqlite   DatabaseType = "sqlite"
)

type DatabaseConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres" json:"postgres"`
	Sqlite   SqliteConfig   `mapstructure:"sqlite" json:"sqlite"`
}

func (dc *DatabaseConfig) GetDatabaseType() DatabaseType {
	if dc.Postgres.Enabled {
		return Postgres
	} else if dc.Sqlite.Enabled {
		return Sqlite
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

type SqliteConfig struct {
	Enabled bool   `mapstructure:"enabled" json:"enabled"`
	DbName  string `mapstructure:"db_name" json:"db_name"`
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
	Mail MailConfig `mapstructure:"mail"`
}

type MailConfig struct {
	Smtp SmtpConfig `mapstructure:"smtp"`
}

type SmtpConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}
