package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const BROKER_URL string = "https://uat.sindbad.tech/"

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

func (client *Client) HealthCheck() {

	req, err := http.NewRequestWithContext(*client.ctx, http.MethodGet, BROKER_URL, nil)
	if err != nil {
		slog.Error("Something went wrong creating a request instance", "error", err)
		return
	}

	resp, err := client.http.Do(req)

	if err != nil {
		slog.Warn("Something went wrong reaching the broker", "err", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body stream", "error", err)
	}

	slog.Info("Broker Responded", "Status Code", resp.StatusCode, "Message", string(body))

}

func main() {

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
				go client.HealthCheck()

			case <-sig:
				slog.Info("Gracefully shutting down")
				return
			}
		}
	}()

	wg.Wait()

}
