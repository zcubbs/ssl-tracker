package main

import (
	"context"
	"embed"
	"github.com/charmbracelet/log"
	"github.com/zcubbs/tlz/db/migrations"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/gapi"
	"github.com/zcubbs/tlz/internal/util"
	"github.com/zcubbs/tlz/pkg/cron"
	"github.com/zcubbs/tlz/pkg/mail"
	"github.com/zcubbs/tlz/task"
	"github.com/zcubbs/tlz/worker"
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

	ctx := context.Background()
	// Migrate database
	err := migrations.Run(cfg.Database)
	if err != nil {
		log.Fatal("failed perform database migrations", "error", err)
	}

	// Connect to database with pgx pool
	conn, err := util.DbConnect(ctx, cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	// Initialize store
	store := db.NewSQLStore(conn)

	// Start cron jobs
	//startCronJobs(store, cfg.Cron)

	// Initialize mailer
	mail.Initialize(cfg.Notification.Mail)

	// Run task worker
	w := worker.New(cfg, store)
	go w.Run()

	// Create gRPC Server
	gs, err := gapi.NewServer(store, w.TaskDistributor, cfg,
		gapi.EmbedAssetsOpts{
			Dir:    swaggerDist,
			Path:   "/swagger/",
			Prefix: "docs/swagger",
		},
		gapi.EmbedAssetsOpts{
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
		go cron.StartCronJob(
			cfg.CheckCertificateValidity.CronPattern,
			t.CheckCertificateValidity,
		)
	}

	if cfg.SendMailNotification.Enabled {
		go cron.StartCronJob(
			cfg.SendMailNotification.CronPattern,
			t.SendMailNotification,
		)
	}
}
