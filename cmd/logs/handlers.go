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

		// Just render the template with empty content - WebSocket will handle all logs
		data := ShowPageData{
			Content:     "", // Empty - WebSocket will populate
			ContainerID: containerID,
		}

		tmpl := template.Must(template.ParseFS(templateFS, "cmd/logs/show.gohtml"))
		tmpl.Execute(w, data)
	}
}
