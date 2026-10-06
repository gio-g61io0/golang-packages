package mail

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
)

type Mailer struct{
	auth smtp.Auth
	addr string
}

func (m *Mailer) SendEmail(messageByte []byte, sender string, recipient  []string) error{
	err := smtp.SendMail(m.addr, m.auth, sender, recipient, messageByte)
	return err
}

func NewMail(username string, password string, host string, port string) *Mailer{
	auth := smtp.PlainAuth("", username, password, host)
	return &Mailer{
		auth: auth,
		addr: fmt.Sprintf("%s:%s", host, port),
	}
}

type Mail struct{
	To string
	Subject string
	Body string
}
func (m Mail)BuildMailMessage() []byte{
	byteMessage := []byte{}
	return fmt.Appendf(byteMessage, "To:%s\r\nSubject:%s\r\n\r\n%s\r\n",m.To,m.Subject, m.Body)
}

func SendAndListen(m *Mailer, mailChan chan Mail, ctx context.Context) error{
	for {
		select {
		case mail := <- mailChan:
			err := m.SendEmail(mail.BuildMailMessage(), "gio.gonzales@sindbad.tech", []string{mail.To})
			if err != nil{
				return err
			}
		case <- ctx.Done():
			slog.Info("Context work is done")
			return ctx.Err()
		}
	}
}


