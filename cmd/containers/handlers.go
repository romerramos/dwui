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
	Containers []containertypes.Summary
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
			Containers: containers,
		}

		funcMap := template.FuncMap{
			"shortenID":   ShortenID,
			"shortenName": ShortenName,
			"urlQuery":    template.URLQueryEscaper,
		}

		tmpl := template.Must(template.New("index.gohtml").Funcs(funcMap).ParseFS(templateFS, "cmd/containers/index.gohtml"))

		tmpl.Execute(w, data)
	}
}
