package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/dwui/containers"
	"github.com/dwui/logs"
	"github.com/dwui/terminal"
)

//go:embed containers/*.gohtml logs/*.gohtml terminal/*.gohtml
var templateFiles embed.FS

//go:embed javascript/*
var staticFiles embed.FS

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// Serve static files
	staticFS, _ := fs.Sub(staticFiles, "javascript")
	r.Handle("/javascript/*", http.StripPrefix("/javascript/", http.FileServer(http.FS(staticFS))))

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
