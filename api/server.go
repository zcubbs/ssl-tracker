package api

import (
	"embed"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/pkg/charmlogfiber"
	"github.com/zcubbs/tlz/token"
	"github.com/zcubbs/tlz/util"
	"net/http"
)

type Server struct {
	store      *db.Store
	app        *fiber.App
	tokenMaker token.Maker
	cfg        util.HttpServerConfig
}

func NewServer(store *db.Store, staticEmbed embed.FS, cfg util.HttpServerConfig) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create new tokenMaker: %w", err)
	}
	s := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		cfg:        cfg,
	}

	s.app = s.getDefaultFiberConfig(staticEmbed, cfg)

	return s, nil
}

func (s *Server) Start() error {
	log.Info("starting HTTP server", "port", s.cfg.Port)
	return s.app.Listen(fmt.Sprintf(":%d", s.cfg.Port))
}

func (s *Server) getDefaultFiberConfig(staticEmbed embed.FS, cfg util.HttpServerConfig) *fiber.App {
	app := fiber.New(fiber.Config{
		EnablePrintRoutes:     cfg.EnablePrintRoutes,
		DisableStartupMessage: true,
	})

	// Initialize default config
	app.Use(cors.New())
	app.Use(requestid.New())
	// Logging Request ID
	app.Use(requestid.New())
	app.Use(charmlogfiber.New(log.Default()))

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowOrigins,
		AllowHeaders: cfg.AllowHeaders,
	}))

	app.Post("/api/users", s.createUser)
	app.Post("/api/users/login", s.loginUser)

	app.Post("/api/domains", s.AddDomain)
	app.Get("/api/domains", s.GetDomains)

	// Serve the frontend
	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(staticEmbed),
		PathPrefix: "web/dist",
	}))

	return app
}
