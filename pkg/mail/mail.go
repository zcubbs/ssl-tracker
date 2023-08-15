package mail

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/jordan-wright/email"
	"github.com/zcubbs/tlz/internal/util"
	"net/smtp"
)

type Mailer interface {
	SendMail(mail Mail) error
}

type Mail struct {
	To            []string `json:"to"`
	Cc            []string `json:"cc"`
	Bcc           []string `json:"bcc"`
	Subject       string   `json:"subject"`
	Content       string   `json:"body"`
	AttachedFiles []string `json:"attachedFiles"`
}

type DefaultSender struct {
	userName    string
	password    string
	fromName    string
	fromAddress string
	serverHost  string
	serverPort  int
}

func NewDefaultSender(cfg util.SmtpConfig) Mailer {
	return &DefaultSender{
		userName:    cfg.Username,
		password:    cfg.Password,
		fromName:    cfg.FromName,
		fromAddress: cfg.FromAddress,
		serverHost:  cfg.Host,
		serverPort:  cfg.Port,
	}
}

func (s DefaultSender) SendMail(mail Mail) error {
	log.Info("Sending email",
		"from", fmt.Sprintf("%s <%s>", s.fromName, s.fromAddress),
		"to", mail.To,
		"cc", mail.Cc,
		"bcc", mail.Bcc,
		"subject", mail.Subject,
		"content", mail.Content,
		"attachedFiles", mail.AttachedFiles,
	)

	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", s.fromName, s.fromAddress)
	e.Subject = mail.Subject
	e.HTML = []byte(mail.Content)
	e.To = mail.To
	e.Cc = mail.Cc
	e.Bcc = mail.Bcc
	for _, file := range mail.AttachedFiles {
		if _, err := e.AttachFile(file); err != nil {
			return fmt.Errorf("failed to attach file %s: %w", file, err)
		}
	}

	smtpGmailHost := fmt.Sprintf("%s:%d", s.serverHost, s.serverPort)
	smtpAuth := smtp.PlainAuth("", s.userName, s.password, s.serverHost)

	return e.Send(smtpGmailHost, smtpAuth)
}
