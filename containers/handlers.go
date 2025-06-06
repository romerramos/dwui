package containers

import (
	"context"
	"embed"
	"html/template"
	"net/http"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type IndexPageData struct {
	PageTitle  string
	Containers []containertypes.Summary
}

func formatID(id string) string {
	return id[0:12]
}

func Index(templateFS embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}

		defer cli.Close()
		containers, err := cli.ContainerList(ctx, containertypes.ListOptions{})
		if err != nil {
			panic(err)
		}

		data := IndexPageData{
			PageTitle:  "Docker Web UI",
			Containers: containers,
		}

		funcMap := template.FuncMap{
			"formatID": formatID,
		}

		tmpl := template.Must(template.New("index.html").Funcs(funcMap).ParseFS(templateFS, "containers/index.html"))

		tmpl.Execute(w, data)
	}
}
