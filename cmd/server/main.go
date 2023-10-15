package main

import (
	"context"
	"embed"
	"github.com/zcubbs/tlz/cmd/server/api"
	db "github.com/zcubbs/tlz/cmd/server/db/sqlc"
	"github.com/zcubbs/tlz/cmd/server/internal/logger"
	"github.com/zcubbs/tlz/cmd/server/internal/task"
	"github.com/zcubbs/tlz/cmd/server/internal/util"
	"github.com/zcubbs/tlz/worker"
	"github.com/zcubbs/x/cron"
	"github.com/zcubbs/x/mail"
)

var (
	Version = "0.0.0"
	Commit  = "none"
	Date    = "unknown"
)

//go:embed web/dist/*
var webDist embed.FS

//go:embed docs/swagger/*
var swaggerDist embed.FS

var cfg util.Config

func init() {
	// Bootstrap configuration
	cfg = util.Bootstrap()
}

func main() {
	// Initialize logger
	//logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	log := logger.GetLogger()

	ctx := context.Background()
	// Migrate database
	err := migration.Run(cfg.Database)
	if err != nil {
		log.Fatal("failed perform database migrations", "error", err)
	}

	// Connect to database with pgx pool
	conn, err := util.DbConnect(ctx, cfg.Database, log)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	// Initialize store
	store := db.NewSQLStore(conn)

	// Start cron jobs
	startCronJobs(store, cfg.Cron)

	// Initialize mailer
	mailer := mail.NewDefaultSender(mail.SmtpConfig{
		Username:    cfg.Notification.Mail.Smtp.Username,
		Password:    cfg.Notification.Mail.Smtp.Password,
		FromName:    cfg.Notification.Mail.Smtp.FromName,
		FromAddress: cfg.Notification.Mail.Smtp.FromAddress,
		Host:        cfg.Notification.Mail.Smtp.Host,
		Port:        cfg.Notification.Mail.Smtp.Port,
	})

	// Run task worker
	w := worker.New(cfg, store, mailer, worker.Attributes{
		ApiDomainName: cfg.Notification.ApiDomainName,
	})
	go w.Run()

	// Create gRPC Server
	gs, err := api.NewServer(store, w.TaskDistributor, cfg,
		api.EmbedAssetsOpts{
			Dir:    swaggerDist,
			Path:   "/swagger/",
			Prefix: "docs/swagger",
		},
		api.EmbedAssetsOpts{
			Dir:    webDist,
			Path:   "/",
			Prefix: "web/dist",
		},
	)
	if err != nil {
		log.Fatal("cannot create grpc server", "error", err)
	}

	// Start HTTP Gateway
	go gs.StartHttpGateway()

	// Start gRPC Server
	gs.StartGrpcServer()
}

func startCronJobs(store db.Store, cfg util.CronConfig) {
	t := task.New(store)
	if cfg.CheckCertificateValidity.Enabled {
		cj := cron.NewJob(
			"check_certificate_validity",
			cfg.CheckCertificateValidity.CronPattern,
			t.CheckCertificateValidity,
			cron.WithLogger(logger.GetLoggerWithName("cron.check_certificate_validity")),
		)

		cj.Start()
	}
}
