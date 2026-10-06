package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"personal-http-server/health_checker/mail"
	"personal-http-server/health_checker/ssh"
	"syscall"

	"github.com/joho/godotenv"
)

type Config struct {
	BrokerUrl    string
	EmailFrom    string
	EmailTo      string
	MailHost     string
	MailPort     string
	MailUsername string
	MailPassword string
}
func loadEnvConfig() (Config, error){
	err := godotenv.Load()

	if err != nil{
		return Config{}, err
	}

	
	return Config{
		MailHost: os.Getenv("MAIL_HOST"),
		MailPort: os.Getenv("MAIL_PORT"),
		MailUsername: os.Getenv("MAIL_USERNAME"),
		MailPassword: os.Getenv("MAIL_PASSWORD"),
	}, nil
}
func main() {
	config, err := loadEnvConfig()
	if err != nil{
		return 
	}

	ctx, cancel := context.WithCancel(context.Background())
	notifyCtx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, os.Interrupt)
	defer stop()

	err = run(notifyCtx, config)

	if err != nil {
		cancel()
	}

}

func run(ctx context.Context, cfg Config) error {
	
	containers := []string{"api_app", "tasks_grinder"}

	//Inits a mailer
	mailer := mail.NewMail(cfg.MailUsername, cfg.MailPassword, cfg.MailHost, cfg.MailPort)



	client, err := ssh.NewDockerCLient(ctx, mailer)
	if err != nil {
		slog.Error("Something went wrong initiating a docker client")
		return err
	}
	for _, container := range containers{
		_, err := client.InitContainerSupervision(ctx, container)
		if err != nil{
			slog.Error("Error initiating supervision", "Error", err)
		}

	}
	err = client.RunSupervision(ctx)

	if err != nil {
		return err
	}
	return nil
}
