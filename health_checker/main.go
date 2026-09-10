package main

import (
	"context"
	"github.com/joho/godotenv"
	gomail "gopkg.in/mail.v2"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const BROKER_URL string = "https://uat.sindbad.tech/"
const EMAIL_FROM string = "test@example.com"
const EMAIL_TO string = "gio.gonzales@sindbad.tech"

type Config struct {
	BrokerUrl    string
	EmailFrom    string
	EmailTo      string
	MailHost     string
	MailPort     int
	MailUsername string
	MailPassword string
}
type Mail struct {
	From    string
	To      string
	Subject string
	Body    string
}

type MailManager struct {
	dialer *gomail.Dialer
}

func (manager *MailManager) SendEmail(mail *Mail) error {

	message := gomail.NewMessage()
	message.SetHeader("From", mail.From)
	message.SetHeader("Subject", mail.Subject)
	message.SetHeader("To", mail.To)

	if err := manager.dialer.DialAndSend(message); err != nil {
		slog.Error("Error sending email", "Error", err)
		return err
	}

	slog.Info("Email Sent Successfully", "To", mail.To)

	return nil

}

func NewEmailManager(config *Config) MailManager {

	return MailManager{
		dialer: gomail.NewDialer(config.MailHost, config.MailPort, config.MailUsername, config.MailPassword),
	}
}

func NewMail(from string, to string, subject string, body string) *Mail {
	return &Mail{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    body,
	}
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	intPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		return nil, err
	}

	return &Config{
		BrokerUrl:    os.Getenv("BROKER_URL"),
		EmailFrom:    os.Getenv("EMAIL_FROM"),
		EmailTo:      os.Getenv("EMAIL_TO"),
		MailHost:     os.Getenv("MAIL_HOST"),
		MailPort:     intPort,
		MailUsername: os.Getenv("MAIL_USERNAME"),
		MailPassword: os.Getenv("MAIL_PASSWORD"),
	}, nil
}

type Client struct {
	http *http.Client
	ctx  *context.Context
}

func InitClient(timeout time.Duration, parentCtx *context.Context) *Client {
	return &Client{
		http: &http.Client{
			Timeout: timeout,
		},
		ctx: parentCtx,
	}

}

func (client *Client) HealthCheck(healthCheckChan chan bool) {

	req, err := http.NewRequestWithContext(*client.ctx, http.MethodGet, BROKER_URL, nil)
	if err != nil {
		slog.Error("Something went wrong creating a request instance", "error", err)
		return
	}

	resp, err := client.http.Do(req)

	if err != nil {
		slog.Warn("Something went wrong reaching the broker", "err", err)
		healthCheckChan <- true
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		slog.Error("Error reading response body stream", "error", err)
		healthCheckChan <- true
	}

	slog.Info("Broker Responded", "Status Code", resp.StatusCode, "Message", string(body))

}

func main() {
	config, err := Load()

	if err != nil {
		slog.Error("Cannot load environment variables")
		return
	}

	mailManager := NewEmailManager(config)

	healthCheckChan := make(chan bool, 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	client := InitClient(time.Duration(time.Second*10), &ctx)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.NewTicker(time.Second * 20)

		for {
			select {
			case <-t.C:
				slog.Info("Calling HealthCheck")
				go client.HealthCheck(healthCheckChan)

			case <-sig:
				slog.Info("Gracefully shutting down")
				cancel()
				return
			}
		}
	}()

	wg.Go(func() {
		for {
			select {
			case val := <-healthCheckChan:
				if val {
					if err := mailManager.SendEmail(NewMail(config.EmailFrom, config.EmailTo, "Error Broker Area", "Test Body is here")); err != nil {
						slog.Error("Something went wrong sending email", "Error", err)
					}
				}
			case <-ctx.Done():
				slog.Info("Gracefully shutting down in email sender routine")
				return
			}
		}

	})

	wg.Wait()

}
