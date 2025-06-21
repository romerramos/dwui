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
	"embed"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ShowPageData struct {
	ContainerID   string
	ContainerName string
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

		data := ShowPageData{
			ContainerID:   containerID,
			ContainerName: containerName,
		}

		tmpl := template.Must(template.ParseFS(templateFS, "cmd/terminal/show.gohtml"))

		tmpl.Execute(w, data)
	}
}
