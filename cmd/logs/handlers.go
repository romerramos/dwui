package logs

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/dwui/cmd/containers"
	"github.com/go-chi/chi/v5"
)

type ShowPageData struct {
	Content       string
	ContainerID   string
	ContainerName string
}

func Show(templateFS embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var containerID = chi.URLParam(req, "containerID")
		var containerName = req.URL.Query().Get("name")

		if containerName == "" {
			containerName = containers.ShortenID(containerID)
		}

		data := ShowPageData{
			Content:       "",
			ContainerID:   containerID,
			ContainerName: containerName,
		}

		tmpl := template.Must(template.ParseFS(templateFS, "cmd/logs/show.gohtml"))
		tmpl.Execute(w, data)
	}
}
