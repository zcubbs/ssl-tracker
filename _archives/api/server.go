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
	"github.com/zcubbs/tlz/internal/util"
	"github.com/zcubbs/tlz/pkg/charmlogfiber"
	"github.com/zcubbs/tlz/pkg/token"
	"net/http"
)

type Server struct {
	store       db.Store
	app         *fiber.App
	tokenMaker  token.Maker
	cfg         util.Config
	staticEmbed *embed.FS
	validate    *XValidator
}

func NewServer(store db.Store, staticEmbed *embed.FS, cfg util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(cfg.Auth.TokenSymmetricKey)
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
	log.Info("starting HTTP server", "port", s.cfg.HttpServer.Port)
	return s.app.Listen(fmt.Sprintf(":%d", s.cfg.HttpServer.Port))
}

func (s *Server) ApplyDefaultConfig() {
	app := fiber.New(fiber.Config{
		EnablePrintRoutes:     s.cfg.HttpServer.EnablePrintRoutes,
		DisableStartupMessage: true,
	})

	s.app = app
	s.addValidator()
	s.applyMiddleware()
	s.addRoutes()
	s.embedStatic()
}

func (s *Server) addRoutes() {
	users := s.app.Group("/api/users")
	users.Post("/login", s.loginUser)
	users.Post("/", s.createUser)

	tokens := s.app.Group("/api/tokens")
	tokens.Post("/refresh", s.renewAccessToken)

	domains := s.app.Group("/api/domains")
	domains.Use(AuthMiddleware(s.tokenMaker))
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
		AllowOrigins: s.cfg.HttpServer.AllowOrigins,
		AllowHeaders: s.cfg.HttpServer.AllowHeaders,
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

func (s *Server) addValidator() {
	val := &XValidator{
		validator: validate,
	}

	err := val.validator.RegisterValidation("domain-name", validDomainName)
	if err != nil {
		log.Fatal("cannot register domain-name validator", "error", err)
	}

	s.validate = val
}
