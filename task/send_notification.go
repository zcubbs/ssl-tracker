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

func SendMailNotification(ctx context.Context) {
	var notifications []db.Notification

	notifications, err := db.Store.GetNotificationsByChannel(ctx, (string)(NotificationChannelEmail))
	if err != nil {
		log.Error("cannot get notifications from db", "error", err)
		return
	}

	for _, notification := range notifications {
		log.Info("sending notification", "notification", notification)

		err = sendEmail(notification)
		if err != nil {
			log.Error("cannot send email notification", "error", err)
		}

		err = db.Store.DeleteNotification(ctx, notification.ID)
		if err != nil {
			log.Error("cannot delete notification", "error", err)
		}
	}

	if len(notifications) == 0 {
		log.Info("no notifications to send")
	}
}

func sendEmail(notification db.Notification) error {
	to := strings.Split(notification.SendTo, ",")
	err := mail.Sender.SendMail(mail.Mail{
		From: "test",
		To:   to,
		Body: notification.Message,
	})
	if err != nil {
		return err
	}
	return nil
}
