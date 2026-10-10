package ssh

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"personal-http-server/health_checker/constants"
	"personal-http-server/health_checker/mail"
	"personal-http-server/health_checker/utils"
	"strings"
	"sync"
	"time"

	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"golang.org/x/sync/errgroup"
)



type ErrResult struct{
	Date ParsedDate
	Container  *ContainerSupervision
	Keywords map[constants.ERRTYPE][]string
}


type ParsedDate struct {
	Day   int
	Month string
	Year  int
}

type DockerClient struct {
	mux        sync.Mutex
	client     *client.Client
	containers map[string]*ContainerSupervision
	mailChan chan mail.Mail
	mailer  *mail.Mailer
}

type ContainerSupervision struct {
	containerID                 string
	containerStreamReaderCloser io.ReadCloser //Holds the stream connection
	containerName               string
	mailChan chan mail.Mail
}

// I think we should only pass the docker client here
func NewDockerCLient(ctx context.Context, mailer *mail.Mailer) (*DockerClient, error) {
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
	mailChan := make(chan mail.Mail)
	

	return &DockerClient{
		client: cli,
		containers: map[string]*ContainerSupervision{},
		mailChan: mailChan,
		mailer: mailer,
	}, nil
}

func (c *DockerClient) RunSupervision(ctx context.Context) error {

	g, gctx := errgroup.WithContext(ctx)
	for _, container := range c.containers {
		//Run the container supervision as a go
		fmt.Println("Container", container.containerName)
		g.Go(func() error { return container.supervise(gctx, c.mailChan) })
	}
	
	g.Go(func() error {return mail.SendAndListen(c.mailer, c.mailChan, gctx)})
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
		mailChan: make(chan mail.Mail),
	}
	c.containers[id] = container
	return container, nil
}

/*
Runs a blocking readStream.
Closes the Container's stream when the parent context is cancelled
*/
func (cs *ContainerSupervision) supervise(ctx context.Context, mailerChan chan mail.Mail) error {

	watchDog := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			slog.Info("Context is done.. closing routine")
			cs.containerStreamReaderCloser.Close()
		case <-watchDog: //Will exit the go routine to prevent go routine leakage
		}
	}()
	defer close(watchDog)

	pr, pw := io.Pipe()


	//This go routine does copies data into pw from src
	go func() {
		_, err := stdcopy.StdCopy(pw, pw, cs.containerStreamReaderCloser)

		//Not an EOF error
		if err != nil{
			slog.Error("Std Copy returns non nil error", "Error", err)
		}
		pw.CloseWithError(err)
	}()

	defer pr.CloseWithError(context.Canceled)

	cs.readStream(pr, mailerChan)
	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}


//Reads the data into pipe sent from docker stream reader
func (cs *ContainerSupervision) readStream(reader *io.PipeReader, mailerChan chan mail.Mail) error {

	scanner := bufio.NewScanner(reader)

	if scanner.Err() != nil {
		return scanner.Err()
	}

	//Read and parse
	for scanner.Scan() {
		line := scanner.Text()
		parsedError, err := ParseLine(line)
		if err != nil {
			slog.Error("Something went wrong parsing single line", "Error", err)
		}

		//Send email here (This must be done through a channel I think or in a go routine)
		if parsedError != nil{

			formattedDetectedKeywords := utils.FormatMapToString(parsedError.Keywords)

			//This is a blocking operation since the channel is not buffered
			mailerChan<-mail.Mail{
				To: "gio.gonzales@carsu.edu.ph",
				Subject: "Log Observer Detected Brokerage Containers Anomalies",
				Body: fmt.Sprintf("An anomaly was detected for container %s\n", parsedError.Container.containerName) + "Detected keywords: \n" + fmt.Sprintf("%s", formattedDetectedKeywords),
			}
			
		}

		fmt.Printf("Parsed Error from container %s %+v\n", cs.containerName, parsedError)
	}
	return nil

}

func RemoveDuplicateCharWithin(line string) string {
	seen := false
	result := ""
	replacement := ([]rune(constants.REPLACEMENT))[0]
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
	replaced := strings.ReplaceAll(line, constants.SEPARATOR, constants.REPLACEMENT)

	cleanedLine := RemoveDuplicateCharWithin(replaced)
	parts := strings.Split(cleanedLine, separator)
	return parts

}


func ParseLine(line string) (*ErrResult, error) {
	parsedErrorKeywords := map[constants.ERRTYPE][]string{}
	
	parts := GetParts(line, constants.REPLACEMENT)

	if len(parts) < 2 {
		return nil, fmt.Errorf("Invalid log line")
	}

	joinedParts := strings.Join(parts[1:], "")
	
	
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


	//check if any keyword exists
	for _, badKeyword := range constants.BADKEYWORDS{
		for _, keyword := range badKeyword.Keywords{
			if strings.Contains(joinedParts, keyword){
				parsedErrorKeywords[badKeyword.ErrorType] = append(parsedErrorKeywords[badKeyword.ErrorType], keyword)
			}
		}
	}

	return &ErrResult{
		Date: ParsedDate{
			Year:  ksaTime.Year(),
			Day:   ksaTime.Day(),
			Month: ksaTime.Month().String(),
		},
		Keywords: parsedErrorKeywords,
	}, nil
}


