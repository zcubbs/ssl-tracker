package main

import (
	"embed"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"log"
	"net/http"
)

var db Database

type Config struct {
	DatabaseType string `json:"database_type"`
}

//go:embed web/dist/*
var f embed.FS

func main() {
	// Load the configuration
	config := Config{
		DatabaseType: "sqlite",
	}
	// Initialize the correct database based on configuration
	if config.DatabaseType == "sqlite" {
		db = NewSQLiteDB("./database.sqlite")
	} else if config.DatabaseType == "postgres" {
		db = &PostgresDB{}
	}

	// Start cron jobs
	// "*/5 * * * * *" means every 5 seconds
	// "-" means only run once
	go StartCheckCertificateValidityCronJob("-")

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

	app.Post("/api/domains", addDomain)
	app.Get("/api/domains", getDomains)

	// Serve the frontend
	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(f),
		PathPrefix: "web/dist",
	}))

	log.Fatal(app.Listen(":8000"))
}
