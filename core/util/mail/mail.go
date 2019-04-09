package mail

import (
	"crypto/tls"
	"github.com/go-gomail/gomail"
)

type MailSender struct {
	Bcc      string
	BccName  string
	Host     string
	Port     int
	Email    string
	Password string
}

type MailMessage struct {
	From     string
	FromName string
	To       string
	ToName   string
	Subject  string
	Body     string
	MailSender
}

func (mm *MailMessage) Sent() error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", mm.From, mm.FromName)
	m.SetAddressHeader("To", mm.To, mm.ToName)
	m.SetHeader("Subject", mm.Subject)

	m.SetHeader("Bcc",
		m.FormatAddress(mm.Bcc, mm.BccName))

	m.SetBody("text/html", mm.Body)

	d := gomail.NewDialer(mm.Host, mm.Port, mm.Email, mm.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
