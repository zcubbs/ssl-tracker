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
	"github.com/zcubbs/tlz/pkg/tls"
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
	// "*/5 * * * * *" means every 5 seconds
	// "-" means only run once
	go tls.StartCheckCertificateValidityCronJob("*/5 * * * * *")

	app := fiber.New(fiber.Config{
		EnablePrintRoutes:     false,
		DisableStartupMessage: true,
	})

	// Initialize default config
	app.Use(cors.New())
	app.Use(requestid.New())
	// Logging Request ID
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		//Format:     "[ip=${ip}]:${port} pid=${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "UTC",
	}))

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8000,http://localhost:5173,http://127.0.0.1:8000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Post("/api/domains", handler.AddDomain)
	app.Get("/api/domains", handler.GetDomains)

	// Serve the frontend
	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(f),
		PathPrefix: "web/dist",
	}))

	log.Infof("Starting HTTP server on port %d", cfg.HttpServer.Port)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.HttpServer.Port)))
}
