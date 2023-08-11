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
	store       db.Store
	app         *fiber.App
	tokenMaker  token.Maker
	cfg         util.HttpServerConfig
	staticEmbed *embed.FS
}

func NewServer(store db.Store, staticEmbed *embed.FS, cfg util.HttpServerConfig) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create new tokenMaker: %w", err)
	}
	s := &Server{
		store:       store,
		tokenMaker:  tokenMaker,
		cfg:         cfg,
		staticEmbed: staticEmbed,
	}

	s.ApplyDefaultConfig()

	return s, nil
}

func (s *Server) Start() error {
	log.Info("starting HTTP server", "port", s.cfg.Port)
	return s.app.Listen(fmt.Sprintf(":%d", s.cfg.Port))
}

func (s *Server) ApplyDefaultConfig() {
	app := fiber.New(fiber.Config{
		EnablePrintRoutes:     s.cfg.EnablePrintRoutes,
		DisableStartupMessage: true,
	})

	s.app = app
	s.applyMiddleware()
	s.addRoutes()
	s.embedStatic()
}

func (s *Server) addRoutes() {
	users := s.app.Group("/api/users")
	users.Use(AuthMiddleware(s.tokenMaker))
	users.Post("/login", s.loginUser)
	users.Post("/", s.createUser)

	domains := s.app.Group("/api/domains")
	domains.Post("/", s.CreateDomain)
	domains.Get("/", s.GetDomains)
}

func (s *Server) applyMiddleware() {
	// Initialize default config
	s.app.Use(cors.New())
	s.app.Use(requestid.New())
	// Logging Request ID
	s.app.Use(requestid.New())
	s.app.Use(charmlogfiber.New(log.Default()))

	// Or extend your config for customization
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: s.cfg.AllowOrigins,
		AllowHeaders: s.cfg.AllowHeaders,
	}))
}

func (s *Server) embedStatic() {
	if s.staticEmbed != nil {
		s.app.Use("/", filesystem.New(filesystem.Config{
			Root:       http.FS(s.staticEmbed),
			PathPrefix: "web/dist",
		}))
	}
}
