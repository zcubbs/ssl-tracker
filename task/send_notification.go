package task

import (
	"context"
	"github.com/charmbracelet/log"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/pkg/mail"
	"strings"
)

type NotificationChannel string

const (
	NotificationChannelEmail NotificationChannel = "email"
	NotificationChannelSlack NotificationChannel = "slack"
	NotificationChannelTeams NotificationChannel = "teams"
	NotificationChannelSms   NotificationChannel = "sms"
)

func (t *Task) SendMailNotification(ctx context.Context, mailSender mail.Mailer) {
	var notifications []db.Notification

	notifications, err := t.store.GetNotificationsByChannel(ctx, (string)(NotificationChannelEmail))
	if err != nil {
		log.Error("Cannot get notifications from db", "error", err)
		return
	}

	for _, notification := range notifications {
		log.Info("Sending notification",
			"notification", notification.Subject,
			"to", notification.SendTo,
		)

		err = t.sendEmail(notification, mailSender)
		if err != nil {
			log.Error("Cannot send email notification", "error", err)
		}

		err = t.store.DeleteNotification(ctx, notification.ID)
		if err != nil {
			log.Error("Cannot delete notification", "error", err)
		}
	}

	if len(notifications) == 0 {
		log.Info("No notifications to send")
	}
}

func (t *Task) sendEmail(notification db.Notification, mailSender mail.Mailer) error {
	to := strings.Split(notification.SendTo, ",")
	err := mailSender.SendMail(mail.Mail{
		Subject: notification.Subject,
		To:      to,
		Content: notification.Message,
	})
	if err != nil {
		return err
	}
	return nil
}
