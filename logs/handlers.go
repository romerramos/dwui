package logs

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ShowPageData struct {
	Content     string
	ContainerID string
}

func Show(templateFS embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var containerID = chi.URLParam(req, "containerID")

		data := ShowPageData{
			Content:     GetByContainer(containerID),
			ContainerID: containerID,
		}

		tmpl := template.Must(template.ParseFS(templateFS, "logs/show.html"))

		tmpl.Execute(w, data)
	}
}
