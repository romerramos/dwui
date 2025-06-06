package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/dwui/containers"
	"github.com/dwui/logs"
	"github.com/dwui/terminal"
)

//go:embed containers/*.html logs/*.html terminal/*.html
var templateFiles embed.FS

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/containers", containers.Index(templateFiles))
	r.Get("/logs/{containerID}", logs.Show(templateFiles))
	r.Get("/terminal/{containerID}", terminal.Socket)
	r.Get("/terminal/view/{containerID}", terminal.Show(templateFiles))

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
