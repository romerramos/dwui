package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/dwui/cmd/containers"
	"github.com/dwui/cmd/home"
	"github.com/dwui/cmd/logs"
	"github.com/dwui/cmd/terminal"
)

//go:embed cmd/**/*.gohtml
var templateFiles embed.FS

//go:embed javascript/*
var javascriptFiles embed.FS

//go:embed assets/stylesheets/*
var stylesheetsFiles embed.FS

//go:embed assets/images/*
var imagesFiles embed.FS

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// Serve static files
	javascriptFS, _ := fs.Sub(javascriptFiles, "javascript")
	stylesheetFS, _ := fs.Sub(stylesheetsFiles, "assets/stylesheets")
	imagesFS, _ := fs.Sub(imagesFiles, "assets/images")
	r.Handle("/javascript/*", http.StripPrefix("/javascript/", http.FileServer(http.FS(javascriptFS))))
	r.Handle("/assets/stylesheets/*", http.StripPrefix("/assets/stylesheets/", http.FileServer(http.FS(stylesheetFS))))
	r.Handle("/assets/images/*", http.StripPrefix("/assets/images/", http.FileServer(http.FS(imagesFS))))

	r.Get("/", home.Show(templateFiles))
	r.Get("/containers", containers.Index(templateFiles))
	r.Get("/logs/{containerID}", logs.Show(templateFiles))
	r.Get("/logs/stream/{containerID}", logs.Socket)
	r.Get("/terminal/{containerID}", terminal.Socket)
	r.Get("/terminal/view/{containerID}", terminal.Show(templateFiles))

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
