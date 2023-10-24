package task

import (
	"bytes"
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/zcubbs/ssl-tracker/cmd/server/db/sqlc"
	"github.com/zcubbs/x/tls"
	"text/template"
	"time"
)

func (t *Task) CheckCertificateValidity(ctx context.Context) {
	domains, err := t.store.GetAllDomains(ctx)
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
			domain.Status.String = (string)(StatusUnknown)
			domain.Status.Valid = true
			if _, err := t.store.UpdateDomain(ctx, db.UpdateDomainParams{
				Status:            domain.Status,
				CertificateExpiry: domain.CertificateExpiry,
				Issuer:            pgtype.Text{},
				Name:              domain.Name,
			}); err != nil {
				log.Error("Cannot update domain", "error", err)
			}
			continue
		}

		newStatus := t.getCertStatus(status.ValidTo)

		err = t.checkNeedsNotification(ctx, domain, newStatus)
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
		if _, err := t.store.UpdateDomain(ctx, db.UpdateDomainParams{
			Status:            domain.Status,
			CertificateExpiry: domain.CertificateExpiry,
			Issuer:            domain.Issuer,
			Name:              domain.Name,
		}); err != nil {
			log.Error("Cannot update domain", "error", err.Error())
		}
	}
}

func (t *Task) getCertStatus(expiry time.Time) Status {
	if expiry.Before(time.Now()) {
		return StatusExpired
	}
	if expiry.Before(time.Now().AddDate(0, 0, 30)) {
		return StatusExpiring
	}
	return StatusValid
}

func (t *Task) checkNeedsNotification(ctx context.Context, domain db.Domain, newStatus Status) error {
	if domain.Status.Valid && domain.Status.String != (string)(newStatus) {
		log.Info("Status changed", "domain", domain.Name, "old", domain.Status.String, "new", newStatus)
		body, err := t.buildNotificationMessage(domain, newStatus)
		if err != nil {
			return err
		}

		if _, err := t.store.InsertNotification(ctx, db.InsertNotificationParams{
			Subject: "Domain " + domain.Name + " is now " + (string)(newStatus),
			Message: body,
			SendTo:  "tlz.test@yopmail.com",
			Channel: (string)(NotificationChannelEmail),
		}); err != nil {
			return err
		}
	}

	return nil
}

func (t *Task) buildNotificationMessage(domain db.Domain, newStatus Status) (string, error) {
	var body bytes.Buffer
	tmpl, err := template.ParseFiles("templates/expiry_notification.html")
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&body, struct {
		Name          string
		PhrasedStatus string
	}{
		Name:          domain.Name,
		PhrasedStatus: "is now <b>" + (string)(newStatus) + "</b>, was <b>" + domain.Status.String + "</b>.",
	})
	if err != nil {
		return "", err
	}

	return body.String(), nil
}
