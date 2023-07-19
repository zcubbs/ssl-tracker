package main

import (
	"embed"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/zcubbs/tlz/db/migrations"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/handler"
	"github.com/zcubbs/tlz/pkg/cron"
	"github.com/zcubbs/tlz/task"
	"github.com/zcubbs/tlz/util"
	"net/http"
)

//go:embed web/dist/*
var f embed.FS

func main() {
	// Bootstrap configuration
	cfg := util.Bootstrap()

	// Connect to database
	db.Connect(cfg.Database)

	// Migrate database
	migrations.Migrate(cfg.Database)

	// Start cron jobs
	startCronJobs(cfg.Cron)

	app := fiber.New(fiber.Config{
		EnablePrintRoutes:     cfg.HttpServer.EnablePrintRoutes,
		DisableStartupMessage: true,
	})

	// Initialize default config
	app.Use(cors.New())
	app.Use(requestid.New())
	// Logging Request ID
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   cfg.HttpServer.TZ,
	}))

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.HttpServer.AllowOrigins,
		AllowHeaders: cfg.HttpServer.AllowHeaders,
	}))

	app.Post("/api/domains", handler.AddDomain)
	app.Get("/api/domains", handler.GetDomains)

	// Serve the frontend
	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(f),
		PathPrefix: "web/dist",
	}))

	log.Info("starting HTTP server", "port", cfg.HttpServer.Port)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.HttpServer.Port)))
}

func startCronJobs(cfg util.CronConfig) {
	if cfg.CheckCertificateValidity.Enabled {
		go cron.StartCronJob(
			cfg.CheckCertificateValidity.CronPattern,
			task.CheckCertificateValidity,
		)
	}
}
