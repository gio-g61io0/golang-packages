package ssh

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"golang.org/x/sync/errgroup"
)

var BADKEYWORDS = []string{"Error", ""}

const SEPARATOR = " "
const REPLACEMENT = "|"

type ParsedError struct {
	Day   int
	Month string
	Year  int
}

type DockerClient struct {
	mux        sync.Mutex
	client     *client.Client
	containers map[string]*ContainerSupervision
}

type ContainerSupervision struct {
	containerID                 string
	containerStreamReaderCloser io.ReadCloser //Holds the stream connection
	containerName               string
}

// I think we should only pass the docker client here
func NewDockerCLient(ctx context.Context) (*DockerClient, error) {
	helper, err := connhelper.GetConnectionHelper("ssh://sindbad_uat")

	if err != nil {
		slog.Warn("error Connecting ssh", "Error", err)
		return nil, err
	}
	cli, err := client.NewClientWithOpts(
		client.WithHost(helper.Host),
		client.WithDialContext(helper.Dialer),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	return &DockerClient{
		client: cli,
		containers: map[string]*ContainerSupervision{},
	}, nil
}

func (c *DockerClient) RunSupervision(ctx context.Context) error {

	g, gctx := errgroup.WithContext(ctx)
	for _, container := range c.containers {
		//Run the container supervision as a go
		fmt.Println("Container", container.containerName)
		g.Go(func() error { return container.supervise(gctx) })
	}
	return g.Wait()
}

func (c *DockerClient) InitContainerSupervision(ctx context.Context, containerName string) (*ContainerSupervision, error) {
	rc, err := c.client.ContainerLogs(ctx, containerName, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "1",
		Details:    false,
		Timestamps: true,
	})

	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("%d", time.Now().UnixMilli())
	container := &ContainerSupervision{
		containerID:                 id,
		containerStreamReaderCloser: rc,
		containerName:               containerName,
	}
	c.containers[id] = container
	return container, nil
}

/*
Runs a blocking readStream.
Closes the Container's stream when the parent context is cancelled
*/
func (cs *ContainerSupervision) supervise(ctx context.Context) error {

	watchDog := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			cs.containerStreamReaderCloser.Close()
		case <-watchDog: //Will exit the go routine to prevent go routine leakage
		}
	}()
	defer close(watchDog)

	pr, pw := io.Pipe()
	go func() {
		_, err := stdcopy.StdCopy(pw, pw, cs.containerStreamReaderCloser)
		pw.CloseWithError(err)
	}()
	defer pr.CloseWithError(context.Canceled)

	cs.readStream(pr)
	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}

func (cs *ContainerSupervision) readStream(reader *io.PipeReader) error {

	scanner := bufio.NewScanner(reader)

	if scanner.Err() != nil {
		return scanner.Err()
	}

	//Read and parse
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("Received a line %s\n", line)
		parsedError, err := ParseLine(line)
		if err != nil {
			slog.Error("Something went wrong parsing single line", "Error", err)
		}
		fmt.Printf("Parsed Error %+v\n", parsedError)
	}
	return nil

}

func RemoveDuplicateCharWithin(line string) string {
	seen := false
	result := ""
	replacement := ([]rune(REPLACEMENT))[0]
	var prev rune

	for _, char := range line {
		if char == replacement {
			if !seen && prev != replacement {
				seen = true
				result += string(char)
			} else {
				seen = false
			}
			prev = char
			continue
		}
		prev = char
		seen = false
		result += string(char)
	}
	return result
}

/*
Replaces a lot of spaces with the REPLACEMENT
It then removes the duplicate replacement in order to split the whole log line into a manageable strings
*/
func GetParts(line string, separator string) []string {
	line = strings.TrimSpace(line)
	replaced := strings.ReplaceAll(line, SEPARATOR, REPLACEMENT)

	cleanedLine := RemoveDuplicateCharWithin(replaced)
	parts := strings.Split(cleanedLine, separator)
	return parts

}

func ParseLine(line string) (*ParsedError, error) {
	parts := GetParts(line, REPLACEMENT)

	if len(parts) < 2 {
		return nil, fmt.Errorf("Invalid log line")
	}
	parsedTime, err := time.Parse(time.RFC3339, parts[0])

	if err != nil {
		slog.Error("Invalid time string", "Error", err)
		return nil, fmt.Errorf("Invalid time string")
	}

	loc, err := time.LoadLocation("Asia/Riyadh")
	if err != nil {
		return nil, fmt.Errorf("Invalid timezone")
	}

	ksaTime := parsedTime.In(loc)

	return &ParsedError{
		Year:  ksaTime.Year(),
		Day:   ksaTime.Day(),
		Month: ksaTime.Month().String(),
	}, nil

}
