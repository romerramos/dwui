package terminal

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ShowPageData struct {
	ContainerID string
}

func Show(w http.ResponseWriter, req *http.Request) {
	var containerID = chi.URLParam(req, "containerID")

	fmt.Println("Container id", containerID)
	data := ShowPageData{
		ContainerID: containerID,
	}

	tmpl := template.Must(template.ParseFiles(("terminal/show.html")))

	tmpl.Execute(w, data)
}
