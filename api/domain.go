package api

import (
	"errors"
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/pkg/token"
	"github.com/zcubbs/tlz/pkg/util"
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

const (
	InvalidDomainName = "Invalid domain name"
	CannotGetDomain   = "Cannot get domain(s)"
)

type CreateDomainRequest struct {
	Name string `json:"name" validate:"required,domain-name"`
}

type CreateDomainResponse struct {
	Name   string    `json:"name"`
	Status string    `json:"status"`
	ID     uuid.UUID `json:"id"`
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

	// Validate the request body
	if !util.IsDomaineNameValid(domainRequest.Name) {
		log.Error(InvalidDomainName, "domain", domainRequest.Name)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": InvalidDomainName,
		})
	}

	// Check if the domain already exists
	_, err := s.store.GetDomain(c.Context(), domainRequest.Name)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Error(
			CannotGetDomain,
			"domain", domainRequest.Name,
			"error", err,
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": CannotGetDomain,
		})
	}

	authPayload := c.Locals(authorizationPayloadKey)
	log.Info("authPayload", "authPayload", authPayload)

	var ownerId uuid.UUID
	if authPayload != nil {
		ownerId = authPayload.(*token.Payload).ID
	}

	var newDomain db.Domain
	// Add the domain to the database
	if newDomain, err = s.store.InsertDomain(c.Context(), db.InsertDomainParams{
		Name:  domainRequest.Name,
		Owner: ownerId,
		Status: pgtype.Text{
			String: (string)(StatusPending),
			Valid:  true,
		},
	}); err != nil {
		log.Error("Cannot add domain", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot add domain",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(CreateDomainResponse{
		Name:   newDomain.Name,
		ID:     newDomain.ID,
		Status: newDomain.Status.String,
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
	domains, err := s.store.GetAllDomains(c.Context())
	if err != nil {
		log.Error(CannotGetDomain, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": CannotGetDomain,
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
			"error", InvalidDomainName,
			"domain", req.Name,
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": InvalidDomainName,
		})
	}

	// Get the domain from the database
	domain, err := s.store.GetDomain(c.Context(), req.Name)
	if err != nil {
		log.Error(CannotGetDomain, "error", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Domain not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": CannotGetDomain,
		})
	}

	// Respond with the domain
	return c.JSON(domain)
}
