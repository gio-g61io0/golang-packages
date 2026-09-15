package ssh

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var BADKEYWORDS = []string{"Error", ""}

const SEPARATOR = " "
const REPLACEMENT = "-"

type ParsedError struct {
	Day   int
	Month string
	Year  int
}

func Connect(ctx context.Context) error {
	helper, err := connhelper.GetConnectionHelper("ssh://sindbad_uat")

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
		Tail:       "1",
		Details:    true,
		Timestamps: true,
	})

	if err != nil {
		slog.Warn("Error finding container logs", "Error", err)
		return err
	}
	go func() {
		<-ctx.Done()
		cli.Close()
		rc.Close()
	}()

	defer rc.Close()

	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		line := scanner.Text()
		parsedError, err := ParseLine(line)
		if err != nil {
			slog.Error("An error occured while parsing a docker log line", "Error", err)
		}
		fmt.Printf("%v\n", parsedError)

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

func ParseLine(line string) (*ParsedError, error) {
	replaced := strings.ReplaceAll(line, SEPARATOR, REPLACEMENT)
	cleanedLine := RemoveDuplicateCharWithin(replaced)
	fmt.Println(cleanedLine)

	trimmedLine := strings.TrimSpace(line)

	if len(trimmedLine) == 0 {
		return nil, fmt.Errorf("Line does not contain anything except spaces")
	}

	fmt.Println(trimmedLine)
	parts := strings.Split(trimmedLine, SEPARATOR)

	fmt.Printf("%v\n", parts)
	if len(parts) < 2 {
		return nil, fmt.Errorf("Invalid log line")
	}

	parsedTime, err := time.Parse(time.RFC3339, parts[1])

	if err != nil {
		return nil, fmt.Errorf("Invalid log line")

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
