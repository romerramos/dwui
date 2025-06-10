package logs

import (
	"bufio"
	"context"
	"log"
	"net/http"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/go-chi/chi/v5"

	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
)

func Socket(w http.ResponseWriter, r *http.Request) {
	var containerID = chi.URLParam(r, "containerID")
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		http.Error(w, "Docker client error", http.StatusInternalServerError)
		return
	}
	defer cli.Close()

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // for local dev, allow all origins
		},
	}
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}
	defer wsConn.Close()

	// Send some initial logs first, then follow new ones
	initialOptions := containertypes.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
		Timestamps: false,
		Tail:       "50", // Send last 50 lines initially for better context
	}

	initialOut, err := cli.ContainerLogs(ctx, containerID, initialOptions)
	if err == nil {
		// Send initial logs
		scanner := bufio.NewScanner(initialOut)
		for scanner.Scan() {
			line := scanner.Text()
			// Clean the line using the same logic as the service
			if len(line) > 8 && line[0] <= 2 {
				line = line[8:]
			}
			if line != "" {
				err = wsConn.WriteMessage(websocket.TextMessage, []byte(line))
				if err != nil {
					initialOut.Close()
					return
				}
			}
		}
		initialOut.Close()
	}

	// Now start following new logs
	followOptions := containertypes.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
		Since:      "1s", // Only get logs from the last 1 second onwards
	}

	followOut, err := cli.ContainerLogs(ctx, containerID, followOptions)
	if err != nil {
		log.Println("Error following logs:", err)
		return
	}
	defer followOut.Close()

	// Read from container logs line by line and write to websocket
	scanner := bufio.NewScanner(followOut)
	for scanner.Scan() {
		line := scanner.Text()
		// Clean the line using the same logic as the service
		if len(line) > 8 && line[0] <= 2 {
			line = line[8:]
		}
		if line != "" {
			err = wsConn.WriteMessage(websocket.TextMessage, []byte(line))
			if err != nil {
				return
			}
		}
	}
}
