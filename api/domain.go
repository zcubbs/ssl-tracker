package api

import (
	"database/sql"
	"errors"
	"fmt"
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

type CreateDomainRequest struct {
	Name string `json:"name"`
}

func (s *Server) CreateDomain(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	// Parse the request body
	var domainRequest CreateDomainRequest
	if err := c.BodyParser(&domainRequest); err != nil {
		log.Error("Cannot parse JSON",
			"error", err,
			"body", c.Body(),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	fmt.Println(domainRequest)

	// Validate the request body
	if !util.IsDomaineNameValid(domainRequest.Name) {
		log.Error("Invalid domain name", "domain", domainRequest.Name)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid domain name",
		})
	}

	// Check if the domain already exists
	_, err := s.store.GetDomain(c.Context(), domainRequest.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Error("Cannot get domain", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot get domain",
		})
	}

	// Add the domain to the database
	if _, err := s.store.InsertDomain(c.Context(), db.InsertDomainParams{
		Name: domainRequest.Name,
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

type GetDomainResponse struct {
	Name              string    `json:"name"`
	Status            Status    `json:"status"`
	Issuer            string    `json:"issuer"`
	CertificateExpiry time.Time `json:"certificate_expiry"`
	Until             string    `json:"until"`
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

	domainResponse := make([]GetDomainResponse, len(domains))
	for i, domain := range domains {
		u := "-"
		if domain.CertificateExpiry.Valid {
			u = util.TimeUntil(domain.CertificateExpiry.Time)
		}

		domainResponse[i] = GetDomainResponse{
			Name:              domain.Name,
			Status:            (Status)(domain.Status.String),
			Issuer:            domain.Issuer.String,
			CertificateExpiry: domain.CertificateExpiry.Time,
			Until:             u,
		}
	}

	// Respond with the domains
	return c.JSON(domainResponse)
}

type getDomainRequestParams struct {
	Name string `params:"name"`
}

func (s *Server) GetDomain(c *fiber.Ctx) error {
	// Get the domain name from the URL
	var req getDomainRequestParams
	if err := c.ParamsParser(&req); err != nil {
		log.Error("Cannot parse URL", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse URL",
		})
	}

	if !util.IsDomaineNameValid(req.Name) {
		log.Error("validation failed",
			"error", "invalid domain name",
			"domain", req.Name,
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid domain name",
		})
	}

	// Get the domain from the database
	domain, err := s.store.GetDomain(c.Context(), req.Name)
	if err != nil {
		log.Error("Cannot get domain", "error", err)
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Domain not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot get domain",
		})
	}

	// Respond with the domain
	return c.JSON(domain)
}
