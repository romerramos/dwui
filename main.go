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

//go:embed assets/javascript/*
var javascriptFiles embed.FS

//go:embed assets/stylesheets/*
var stylesheetsFiles embed.FS

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// Serve static files
	javascriptFS, _ := fs.Sub(javascriptFiles, "assets/javascript")
	stylesheetFS, _ := fs.Sub(stylesheetsFiles, "assets/stylesheets")
	r.Handle("/assets/javascript/*", http.StripPrefix("/assets/javascript/", http.FileServer(http.FS(javascriptFS))))
	r.Handle("/assets/stylesheets/*", http.StripPrefix("/assets/stylesheets/", http.FileServer(http.FS(stylesheetFS))))

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
