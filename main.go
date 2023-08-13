package main

import (
	"context"
	"embed"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
	"github.com/zcubbs/tlz/api"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/pkg/cron"
	"github.com/zcubbs/tlz/pkg/mail"
	"github.com/zcubbs/tlz/pkg/util"
	"github.com/zcubbs/tlz/task"
)

//go:embed web/dist/*
var webDist embed.FS

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
	startCronJobs(store, cfg.Cron)

	// Initialize mailer
	mail.Initialize(cfg.Notification.Mail)

	// Create Server
	s, err := api.NewServer(store, &webDist, cfg.HttpServer)
	if err != nil {
		log.Fatal("cannot create server", "error", err)
	}

	err = s.Start()
	if err != nil {
		log.Fatal("cannot start server", "error", err)
	}
}

func CloseDbConn(conn *pgx.Conn, ctx context.Context) {
	func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.Error("cannot close database connection", "error", err)
		}
	}(conn, ctx)
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
