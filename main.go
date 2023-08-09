package main

import (
	"embed"
	"github.com/charmbracelet/log"
	"github.com/zcubbs/tlz/api"
	"github.com/zcubbs/tlz/db/migrations"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/pkg/cron"
	"github.com/zcubbs/tlz/pkg/mail"
	"github.com/zcubbs/tlz/task"
	"github.com/zcubbs/tlz/util"
)

//go:embed web/dist/*
var f embed.FS

func main() {
	// Bootstrap configuration
	cfg := util.Bootstrap()

	// Initialize logger
	//logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	// Connect to database
	conn, err := util.Connect(cfg.Database)
	if err != nil {
		log.Fatal("cannot connect to database", "error", err)
	}

	// Initialize store
	store := db.NewStore(conn)

	// Migrate database
	err = migrations.Migrate(store, cfg.Database.GetDatabaseType())
	if err != nil {
		log.Fatal("cannot migrate database", "error", err)
	}

	// Start cron jobs
	startCronJobs(store, cfg.Cron)

	// Initialize mailer
	mail.Initialize(cfg.Notification.Mail)

	// Create Server
	s, err := api.NewServer(store, f, cfg.HttpServer)
	if err != nil {
		log.Fatal("cannot create server", "error", err)
	}
	err = s.Start()
	if err != nil {
		log.Fatal("cannot start server", "error", err)
	}
}

func startCronJobs(store *db.Store, cfg util.CronConfig) {
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
