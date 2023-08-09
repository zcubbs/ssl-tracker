package api

import (
	"database/sql"
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/util"
	"time"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusValid    Status = "valid"
	StatusUnknown  Status = "unknown"
	StatusExpired  Status = "expired"
	StatusExpiring Status = "expiring"
)

type DomainWrapper struct {
	Name              string    `json:"name"`
	Status            Status    `json:"status"`
	Issuer            string    `json:"issuer"`
	CertificateExpiry time.Time `json:"certificate_expiry"`
	Until             string    `json:"until"`
}

func (s *Server) AddDomain(c *fiber.Ctx) error {
	log.Info(c.Body())
	// Parse the request body
	var domain db.Domain
	if err := c.BodyParser(&domain); err != nil {
		log.Error("Cannot parse JSON", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	log.Info("Adding domain", "name", domain.Name)

	// Add the domain to the database
	if _, err := s.store.InsertDomain(c.Context(), db.InsertDomainParams{
		Name: domain.Name,
		Status: sql.NullString{
			String: (string)(StatusPending),
			Valid:  true,
		},
	}); err != nil {
		log.Error("Cannot add domain", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot add domain",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Domain added",
	})
}

func (s *Server) GetDomains(c *fiber.Ctx) error {
	// Get the domains from the database
	domains, err := s.store.GetDomains(c.Context())
	if err != nil {
		log.Error("Cannot get domains", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot get domains",
		})
	}

	domainWrappers := make([]DomainWrapper, len(domains))
	for i, domain := range domains {
		u := "-"
		if domain.CertificateExpiry.Valid {
			u = util.TimeUntil(domain.CertificateExpiry.Time)
		}

		domainWrappers[i] = DomainWrapper{
			Name:              domain.Name,
			Status:            (Status)(domain.Status.String),
			Issuer:            domain.Issuer.String,
			CertificateExpiry: domain.CertificateExpiry.Time,
			Until:             u,
		}
	}

	// Respond with the domains
	return c.JSON(domainWrappers)
}
