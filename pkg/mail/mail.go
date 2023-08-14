package mail

import (
	"github.com/charmbracelet/log"
	"github.com/zcubbs/tlz/internal/util"
)

var Sender Mailer

type Mailer interface {
	SendMail(mail Mail) error
}

type Mail struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

func Initialize(cfg util.MailConfig) {
	if !cfg.Smtp.Enabled {
		log.Warn("no mailer enabled, please enable either smtp, sendgrid or gomail")
		return
	}

	if cfg.Smtp.Enabled {
		Sender = GoMail{
			Host:     cfg.Smtp.Host,
			Port:     cfg.Smtp.Port,
			Username: cfg.Smtp.Username,
			Password: cfg.Smtp.Password,
			From:     cfg.Smtp.From,
		}
		log.Info("configured SMTP server")
		return
	}
}
