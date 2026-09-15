package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"personal-http-server/health_checker/ssh"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	gomail "gopkg.in/mail.v2"
)

type HealthCheckErr string

const (
	UnResponsive  HealthCheckErr = "Unresponsive"
	BadStatusCode HealthCheckErr = "Bad Status Code"
)

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

type HealthCheckResult struct {
	Err     error
	Message string
	ErrType HealthCheckErr
}

type MailManager struct {
	dialer *gomail.Dialer
}

func (manager *MailManager) SendEmail(mail *Mail) error {

	message := gomail.NewMessage()
	message.SetHeader("From", mail.From)
	message.SetHeader("Subject", mail.Subject)
	message.SetHeader("To", mail.To)
	message.SetBody("plain/text", mail.Body)

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
	http   *http.Client
	ctx    *context.Context
	config *Config
}

func InitClient(timeout time.Duration, parentCtx *context.Context, config *Config) *Client {
	return &Client{
		http: &http.Client{
			Timeout: timeout,
		},
		ctx:    parentCtx,
		config: config,
	}

}

func (client *Client) HealthCheck(healthCheckChan chan HealthCheckResult) {

	timeoutCtx, cancel := context.WithTimeout(*client.ctx, time.Duration(time.Second*10))
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutCtx, http.MethodGet, client.config.BrokerUrl, nil)
	if err != nil {
		slog.Error("Something went wrong creating a request instance", "error", err)
		return
	}

	resp, err := client.http.Do(req)

	if err != nil {
		slog.Warn("Something went wrong reaching the broker", "err", err)

		healthCheckChan <- HealthCheckResult{
			Err:     err,
			Message: "Cannot reach broker server. Broker server is unreachable. Please immediately check the server through some monitoring service",
			ErrType: UnResponsive,
		}
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	//NOTE::Not a broker error. Might be an error in our side
	if err != nil {
		slog.Error("Error reading response body stream", "error", err)
		healthCheckChan <- HealthCheckResult{
			Err:     err,
			Message: fmt.Sprintf(""),
			ErrType: "",
		}
	}

	//NOTE:: This is commented for now since the broker dont have any /health endpoint for returning 200 when everything is working
	// if resp.StatusCode != http.StatusOK {
	// 	healthCheckChan <- HealthCheckResult{
	// 		Err:     fmt.Errorf("Broker return non 2xx status code"),
	// 		Message: "Broker return non 2xx status code",
	// 		ErrType: BadStatusCode,
	// 	}
	// }

	slog.Info("Broker Responded", "Status Code", resp.StatusCode, "Message", string(body))

}

func QueryDB(ctx context.Context) error {
	select {
	case <-time.After(3 * time.Second):
		fmt.Println("Query is done")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func main() {
	parentCtx, parentCancel := context.WithCancel(context.Background())
	testTimeoutCtx, testCancel := context.WithTimeout(context.Background(), time.Second*1)
	defer testCancel()
	defer parentCancel()

	if err := QueryDB(testTimeoutCtx); err != nil {
		fmt.Println("Timed out:", err) // context.DeadlineExceeded
	} else {
		fmt.Println("Nailed it!")
	}

	config, err := Load()

	if err != nil {
		slog.Error("Cannot load environment variables")
		return
	}

	mailManager := NewEmailManager(config)

	healthCheckChan := make(chan HealthCheckResult, 1)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	client := InitClient(time.Duration(time.Second*10), &parentCtx, config)
	var wg sync.WaitGroup

	wg.Go(func() {
		defer wg.Done()
		t := time.NewTicker(time.Second * 20)

		for {
			select {
			case <-t.C:
				slog.Info("Calling HealthCheck")
				go client.HealthCheck(healthCheckChan)

			case <-parentCtx.Done():
				slog.Info("Gracefully shutting down")
				return
			}
		}
	})

	//I dont know if this is appropriate
	wg.Go(func() {
		<-sig
		parentCancel()
	})

	wg.Go(func() {
		ssh.Connect(parentCtx)
	})

	wg.Go(func() {
		for {
			select {
			case val := <-healthCheckChan:
				switch val.ErrType {
				case BadStatusCode:
					if err := mailManager.SendEmail(NewMail(config.EmailFrom, config.EmailTo, string(val.ErrType), val.Message)); err != nil {
						slog.Error("Something went wrong sending email", "Error", err)
					}
				case UnResponsive:
					if err := mailManager.SendEmail(NewMail(config.EmailFrom, config.EmailTo, string(val.ErrType), val.Message)); err != nil {
						slog.Error("Something went wrong sending email", "Error", err)
					}
				}
			case <-parentCtx.Done():
				slog.Info("Gracefully shutting down in email sender routine")
				return
			}
		}

	})

	wg.Wait()

}
