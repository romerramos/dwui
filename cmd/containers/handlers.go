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

func shortenID(id string) string {
	return shortenWithAmount(id, 12)
}

func shortenName(name string) string {
	// Remove leading '/' from container name if present
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}
	return shortenWithAmount(name, 25)
}

func shortenWithAmount(text string, amount int) string {
	if len(text) > amount {
		return text[:amount]
	}
	return text
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
			"shortenID":   shortenID,
			"shortenName": shortenName,
			"urlQuery":    template.URLQueryEscaper,
		}

		tmpl := template.Must(template.New("index.gohtml").Funcs(funcMap).ParseFS(templateFS, "cmd/containers/index.gohtml"))

		tmpl.Execute(w, data)
	}
}
