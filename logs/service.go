package logs

import (
	"context"
	"io"
	"regexp"
	"strings"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func GetByContainer(containerID string) string {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	defer cli.Close()

	options := containertypes.LogsOptions{ShowStdout: true}

	out, err := cli.ContainerLogs(ctx, containerID, options)
	if err != nil {
		panic(err)
	}

	defer out.Close()
	logBytes, err := io.ReadAll(out)
	if err != nil {
		panic(err)
	}

	logs := string(logBytes)
	return cleanLogOutput(logs)
}

func cleanLogOutput(logs string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	cleaned := ansiRegex.ReplaceAllString(logs, "")
	lines := strings.Split(cleaned, "\n")
	var cleanedLines []string

	for _, line := range lines {
		if len(line) > 8 {
			if line[0] <= 2 && len(line) > 8 {
				line = line[8:]
			}
		}
		if strings.TrimSpace(line) != "" {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return strings.Join(cleanedLines, "\n")
}
