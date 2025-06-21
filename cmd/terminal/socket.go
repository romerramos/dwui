// DWUI (Docker Web UI)
// Copyright (C) 2025 Romer Ramos
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package terminal

import (
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
