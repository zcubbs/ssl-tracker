package mail

import (
	"github.com/stretchr/testify/require"
	"github.com/zcubbs/tlz/internal/util"
	"testing"
)

func TestSendEmailWithGmail(t *testing.T) {
	cfg := util.Bootstrap()

	sender := NewDefaultSender(cfg.Notification.Mail.Smtp)

	content := `
	<h1>Test Email</h1>
	<p>This is a test email</p>
	`

	mail := Mail{
		To:            []string{cfg.Notification.Mail.Smtp.FromAddress},
		Cc:            []string{},
		Bcc:           []string{},
		Subject:       "Test email",
		Content:       content,
		AttachedFiles: []string{"./testdata/attachement.md"},
	}

	err := sender.SendMail(mail)
	require.NoError(t, err)
}
