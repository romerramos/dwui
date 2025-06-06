package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/dwui/containers"
	"github.com/dwui/logs"
	"github.com/dwui/terminal"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/containers", containers.Index)
	r.Get("/logs/{containerID}", logs.Show)
	r.Get("/terminal/{id}", terminal.Socket)
	r.Get("/terminal/view/{id}", terminal.Show)

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
