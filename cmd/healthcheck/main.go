package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"personal-http-server/health_checker/ssh"
	"syscall"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	notifyCtx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, os.Interrupt)
	defer stop()

	err := run(notifyCtx)

	if err != nil {
		cancel()
	}

}

func run(ctx context.Context) error {
	containers := []string{"api_app", "tasks_grinder"}
	client, err := ssh.NewDockerCLient(ctx)
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
