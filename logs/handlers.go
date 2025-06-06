package logs

import (
	"html/template"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/go-chi/chi/v5"
)

type ShowPageData struct {
	Content     string
	ContainerID string
}

func Show(w http.ResponseWriter, req *http.Request) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	defer cli.Close()
	var containerID = chi.URLParam(req, "containerID")

	data := ShowPageData{
		Content:     GetByContainer(containerID),
		ContainerID: containerID,
	}

	tmpl := template.Must(template.ParseFiles(("logs/show.html")))

	tmpl.Execute(w, data)
}
