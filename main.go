package main

import (
	"context"
	"embed"
	"github.com/charmbracelet/log"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/gapi"
	"github.com/zcubbs/tlz/pkg/cron"
	"github.com/zcubbs/tlz/pkg/mail"
	"github.com/zcubbs/tlz/pkg/util"
	"github.com/zcubbs/tlz/task"
)

//go:embed web/dist/*
var webDist embed.FS

//go:embed docs/swagger/*
var swaggerDist embed.FS

func main() {
	// Bootstrap configuration
	cfg := util.Bootstrap()

	// Initialize logger
	//logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	ctx := context.Background()
	// Connect to database
	conn, err := util.Connect(ctx, cfg.Database)
	if err != nil {
		log.Fatal("cannot connect to database", "error", err)
	}

	// Initialize store
	store := db.NewSQLStore(conn)

	// Start cron jobs
	//startCronJobs(store, cfg.Cron)

	// Initialize mailer
	mail.Initialize(cfg.Notification.Mail)

	// Create gRPC Server
	gs, err := gapi.NewServer(store, cfg,
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
	go gs.StartGateway()

	// Start gRPC Server
	gs.Start()

	//// Create Https Server
	//s, err := api.NewServer(store, &webDist, cfg)
	//if err != nil {
	//	log.Fatal("cannot create server", "error", err)
	//}
	//// Start Http Server
	//err = s.Start()
	//if err != nil {
	//	log.Fatal("cannot start server", "error", err)
	//}
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
