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
