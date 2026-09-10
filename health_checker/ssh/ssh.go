package ssh

import (
	"context"
	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"log/slog"
	"os"
	// "golang.org/x/tools/go/analysis/passes/nilfunc"
)

type SSH struct {
}

func Connect(ctx context.Context) error {
	helper, err := connhelper.GetConnectionHelper("ssh://sindbad_uat")
	// buffer := make([]byte, 2046)

	if err != nil {
		slog.Warn("error Connecting ssh", "Error", err)
		return err
	}
	cli, err := client.NewClientWithOpts(
		client.WithHost(helper.Host),
		client.WithDialContext(helper.Dialer),
		client.WithAPIVersionNegotiation(),
	)

	if err != nil {
		return err
	}
	defer cli.Close()
	rc, err := cli.ContainerLogs(ctx, "api_app", container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "100",
		Details:    true,
		Timestamps: true,
	})

	if err != nil {
		slog.Warn("Error finding container logs", "Error", err)
		return err
	}
	defer rc.Close()

	// for {
	// 	nRead, err := rc.Read(buffer)
	//
	// 	if err != nil {
	// 		slog.Error("Error reading container stream reader", "Error", err)
	// 		return err
	// 	}
	//
	// 	if nRead == 0 {
	// 		slog.Warn("Nothing to read", "Read Bytes", nRead)
	// 		break
	// 	}
	// }

	if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, rc); err != nil && err == io.EOF {
		slog.Info("End of File!!")
	}

	return nil

}
