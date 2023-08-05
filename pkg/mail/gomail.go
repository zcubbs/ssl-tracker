package mail

import "gopkg.in/gomail.v2"

type GoMail struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func (g GoMail) SendMail(mail Mail) error {
	m := gomail.NewMessage()
	m.SetHeader("From", g.From)
	m.SetHeader("To", "tlz.test@yopmail.com")
	m.SetHeader("Subject", "[TLZ] Certificate Status - "+mail.Subject)
	m.SetBody("text/html", mail.Body)

	d := gomail.NewDialer(
		g.Host,
		g.Port,
		g.Username,
		g.Password,
	)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	return nil
}
