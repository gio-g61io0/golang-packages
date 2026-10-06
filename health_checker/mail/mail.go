package mail

import (
	"fmt"
	"net/smtp"
)

type Mailer struct{
	auth smtp.Auth
	addr string
}

func (m *Mailer) SendEmail(message string, sender string, recipient  []string) error{

	messageByte := []byte(message)
	err := smtp.SendMail(m.addr, m.auth, sender, recipient, messageByte)
	return err
}

func BuildMailMessage(to string, subject string, body string) []byte{
	byteMessage := []byte{}
	return fmt.Appendf(byteMessage, "To:%s\r\nSubject:%s\r\n\r\n%s\r\n",to,subject, body)
}
func NewMail(username string, password string, host string, port int) *Mailer{
	auth := smtp.PlainAuth("", username, password, host)
	return &Mailer{
		auth: auth,
		addr: fmt.Sprintf("%s:%d", host, port),
	}
}


