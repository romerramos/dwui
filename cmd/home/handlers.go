package home

import (
	"embed"
	"net/http"
	"text/template"
)

type ShowPageData struct {
	PageTitle string
}

func Show(templateFS embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data := ShowPageData{
			PageTitle: "Docker Web UI",
		}

		tmpl := template.Must(template.ParseFS(templateFS, "cmd/home/show.gohtml"))
		tmpl.Execute(w, data)
	}
}
