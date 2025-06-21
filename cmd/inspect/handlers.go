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

package inspect

import (
	"context"
	"embed"
	"html/template"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/go-chi/chi/v5"
)

type ShowPageData struct {
	ContainerID     string
	ContainerName   string
	EnvironmentVars []EnvVar
	Ports           []Port
}

type EnvVar struct {
	Key   string
	Value string
}

type Port struct {
	ContainerPort string
	HostPort      string
	Type          string
}

func Show(templateFS embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var containerID = chi.URLParam(req, "containerID")
		var containerName = req.URL.Query().Get("name")

		// Fallback to shortened ID if name is not provided
		if containerName == "" {
			if len(containerID) > 12 {
				containerName = containerID[:12]
			} else {
				containerName = containerID
			}
		}

		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			http.Error(w, "Docker client error", http.StatusInternalServerError)
			return
		}
		defer cli.Close()

		// Inspect the container to get detailed information
		containerJSON, err := cli.ContainerInspect(ctx, containerID)
		if err != nil {
			http.Error(w, "Container inspect error", http.StatusInternalServerError)
			return
		}

		// Extract environment variables
		var envVars []EnvVar
		for _, env := range containerJSON.Config.Env {
			// Split by first '=' to separate key and value
			for i, char := range env {
				if char == '=' {
					key := env[:i]
					value := env[i+1:]
					envVars = append(envVars, EnvVar{Key: key, Value: value})
					break
				}
			}
		}

		// Extract port mappings
		var ports []Port
		if containerJSON.NetworkSettings != nil && containerJSON.NetworkSettings.Ports != nil {
			for containerPort, bindings := range containerJSON.NetworkSettings.Ports {
				portStr := string(containerPort)
				if bindings != nil && len(bindings) > 0 {
					for _, binding := range bindings {
						hostPort := binding.HostPort
						if hostPort == "" {
							hostPort = "Not mapped"
						}
						ports = append(ports, Port{
							ContainerPort: portStr,
							HostPort:      hostPort,
							Type:          "TCP/UDP",
						})
					}
				} else {
					// Port is exposed but not mapped
					ports = append(ports, Port{
						ContainerPort: portStr,
						HostPort:      "Not mapped",
						Type:          "TCP/UDP",
					})
				}
			}
		}

		data := ShowPageData{
			ContainerID:     containerID,
			ContainerName:   containerName,
			EnvironmentVars: envVars,
			Ports:           ports,
		}

		tmpl := template.Must(template.ParseFS(templateFS, "cmd/inspect/show.gohtml"))
		tmpl.Execute(w, data)
	}
}
