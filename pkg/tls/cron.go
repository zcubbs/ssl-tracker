package tls

import (
	"context"
	"database/sql"
	"github.com/charmbracelet/log"
	"github.com/robfig/cron/v3"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/handler"
)

func StartCheckCertificateValidityCronJob(cronPattern string) {
	ctx := context.Background()

	if cronPattern == "" {
		log.Info("No cron pattern provided, not starting cron job")
		return
	}

	if cronPattern == "-" {
		log.Info("Checking certificate validity once")
		checkCertificateValidity(ctx)
		log.Info("Done checking certificate validity")
		return
	}

	c := cron.New(cron.WithSeconds()) // cron with second-level precision
	_, err := c.AddFunc(cronPattern, func() {
		checkCertificateValidity(ctx)
	})
	if err != nil {
		log.Fatalf("Cannot create cron job: %v", err)
	}

	log.Info("Starting cron job")
	c.Start()
}

func checkCertificateValidity(ctx context.Context) {
	// Get all domains from the database
	if db.Store == nil {
		log.Fatal("Database store not initialized")
	}

	domains, err := db.Store.GetDomains(ctx)
	if err != nil {
		log.Printf("Cannot get domains: %v", err)
		return
	}

	// Check the SSL certificate of each domain
	for _, domain := range domains {
		log.Info("Checking certificate of", domain.Name)
		status, err := CheckCertificate(domain.Name)
		if err != nil {
			log.Printf("Cannot check certificate of %s: %v", domain.Name, err)
			domain.CertificateExpiry.Valid = false
			domain.Status.String = (string)(handler.StatusUnknown)
			domain.Status.Valid = true
			if _, err := db.Store.UpdateDomain(ctx, db.UpdateDomainParams{
				Status:            sql.NullString{},
				CertificateExpiry: sql.NullTime{},
				Issuer:            sql.NullString{},
				Name:              "",
			}); err != nil {
				log.Printf("Cannot update domain: %v", err)
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
			Status:            sql.NullString{},
			CertificateExpiry: sql.NullTime{},
			Issuer:            sql.NullString{},
			Name:              "",
		}); err != nil {
			log.Printf("Cannot update domain: %v", err)
		}
	}
}
