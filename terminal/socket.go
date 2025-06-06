package terminal

import (
	"context"
	"log"
	"net/http"

	containertypes "github.com/docker/docker/api/types/container"

	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
)

func Socket(w http.ResponseWriter, r *http.Request) {
	containerID := r.PathValue("id")
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

	execResp, err := cli.ContainerExecCreate(ctx, containerID, containertypes.ExecOptions{
		Cmd:          []string{"/bin/bash"},
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	})
	if err != nil {
		log.Println("Exec create error:", err)
		return
	}

	hijackResp, err := cli.ContainerExecAttach(ctx, execResp.ID, containertypes.ExecStartOptions{Tty: true})
	if err != nil {
		log.Println("Exec attach error:", err)
		return
	}
	defer hijackResp.Close()

	// Stream container output to browser
	go func() {
		for {
			msgType, msg, err := wsConn.ReadMessage()
			if err != nil {
				return
			}
			if msgType == websocket.TextMessage || msgType == websocket.BinaryMessage {
				_, err := hijackResp.Conn.Write(msg)
				if err != nil {
					return
				}
			}
		}
	}()

	// Read from container and write to websocket
	buf := make([]byte, 1024)
	for {
		n, err := hijackResp.Reader.Read(buf)
		if err != nil {
			return
		}
		err = wsConn.WriteMessage(websocket.TextMessage, buf[:n])
		if err != nil {
			return
		}
	}
}
