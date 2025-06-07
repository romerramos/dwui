package terminal

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ShowPageData struct {
	ContainerID string
}

func Show(templateFS embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var containerID = chi.URLParam(req, "containerID")

		data := ShowPageData{
			ContainerID: containerID,
		}

		tmpl := template.Must(template.ParseFS(templateFS, "terminal/show.gohtml"))

		tmpl.Execute(w, data)
	}
}
