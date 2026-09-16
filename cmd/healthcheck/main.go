package main

import (
	"context"
	"log/slog"
	"personal-http-server/health_checker/ssh"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	err := run(ctx)

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
