package mail

type SendgridMailer struct {
	ApiKey string
}

func (s SendgridMailer) SendMail(mail Mail) error {
	return nil
}
