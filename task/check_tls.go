package task

import (
	"bytes"
	"context"
	"database/sql"
	"github.com/charmbracelet/log"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/handler"
	"github.com/zcubbs/tlz/pkg/tls"
	"html/template"
	"time"
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

		newStatus := getCertStatus(status.ValidTo)

		err = checkNeedsNotification(ctx, domain, newStatus)
		if err != nil {
			log.Error("failed to setup notification", "domain", domain.Name, "error", err)
		}

		// Update the domain in the database
		domain.CertificateExpiry.Time = status.ValidTo
		domain.CertificateExpiry.Valid = true
		domain.Status.String = (string)(newStatus)
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

func getCertStatus(expiry time.Time) handler.Status {
	if expiry.Before(time.Now()) {
		return handler.StatusExpired
	}
	if expiry.Before(time.Now().AddDate(0, 0, 30)) {
		return handler.StatusExpiring
	}
	return handler.StatusValid
}

func checkNeedsNotification(ctx context.Context, domain db.Domain, newStatus handler.Status) error {
	if domain.Status.Valid && domain.Status.String != (string)(newStatus) {
		log.Info("Status changed", "domain", domain.Name, "old", domain.Status.String, "new", newStatus)
		body, err := buildNotificationMessage(domain, newStatus)
		if err != nil {
			return err
		}

		if _, err := db.Store.InsertNotification(ctx, db.InsertNotificationParams{
			Message: body,
			SendTo:  "tlz.test@yopmail.com",
			Channel: (string)(NotificationChannelEmail),
		}); err != nil {
			return err
		}
	}

	return nil
}

func buildNotificationMessage(domain db.Domain, newStatus handler.Status) (string, error) {
	var body bytes.Buffer
	t, err := template.ParseFiles("templates/expiry_notification.html")
	if err != nil {
		return "", err
	}

	err = t.Execute(&body, struct {
		Name          string
		PhrasedStatus string
	}{
		Name:          domain.Name,
		PhrasedStatus: "is " + (string)(newStatus),
	})
	if err != nil {
		return "", err
	}

	return body.String(), nil
}
