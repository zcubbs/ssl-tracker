package task

import (
	"context"
	"database/sql"
	"github.com/charmbracelet/log"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/handler"
	"github.com/zcubbs/tlz/pkg/tls"
)

func CheckCertificateValidity(ctx context.Context) {
	// Get all domains from the database
	if db.Store == nil {
		log.Fatal("Database store not initialized")
	}

	domains, err := db.Store.GetDomains(ctx)
	if err != nil {
		log.Error("Cannot get domains", "error", err)
		return
	}

	// Check the SSL certificate of each domain
	for _, domain := range domains {
		log.Info("Checking certificate", "domain", domain.Name)
		status, err := tls.CheckCertificate(domain.Name)
		if err != nil {
			log.Warn("Cannot check certificate", "domain", domain.Name, "error", err)
			domain.CertificateExpiry.Valid = false
			domain.Status.String = (string)(handler.StatusUnknown)
			domain.Status.Valid = true
			if _, err := db.Store.UpdateDomain(ctx, db.UpdateDomainParams{
				Status:            domain.Status,
				CertificateExpiry: domain.CertificateExpiry,
				Issuer:            sql.NullString{},
				Name:              domain.Name,
			}); err != nil {
				log.Error("Cannot update domain", "error", err)
			}
			continue
		}

		log.Info(status)

		// Update the domain in the database
		domain.CertificateExpiry.Time = status.ValidTo
		domain.CertificateExpiry.Valid = true
		domain.Status.String = (string)(handler.StatusValid)
		domain.Status.Valid = true
		domain.Issuer.String = status.Issuer
		domain.Issuer.Valid = true
		if _, err := db.Store.UpdateDomain(ctx, db.UpdateDomainParams{
			Status:            domain.Status,
			CertificateExpiry: domain.CertificateExpiry,
			Issuer:            domain.Issuer,
			Name:              domain.Name,
		}); err != nil {
			log.Error("Cannot update domain", "error", err.Error())
		}
	}
}
