package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/dwui/cmd/auth"
	"github.com/dwui/cmd/containers"
	"github.com/dwui/cmd/home"
	"github.com/dwui/cmd/inspect"
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
	// Parse command line flags
	var password string
	flag.StringVar(&password, "password", "", "Password for authentication (if not provided, a random one will be generated)")
	flag.Parse()

	// Set up authentication
	if password == "" {
		generatedPassword, err := auth.GenerateRandomPassword()
		if err != nil {
			fmt.Printf("Error generating password: %v\n", err)
			os.Exit(1)
		}
		password = generatedPassword
		fmt.Printf("\n🔐 DWUI Authentication Password: %s\n", password)
		fmt.Printf("🌐 Server will be available at: http://localhost:8080\n")
		fmt.Printf("💡 Use this password to sign in to the web interface\n\n")
	} else {
		fmt.Printf("🔐 Using provided password for authentication\n")
		fmt.Printf("🌐 Server will be available at: http://localhost:8080\n\n")
	}

	auth.SetPassword(password)

	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\n🛑 Gracefully shutting down...")
		auth.CloseDB() // Close BadgerDB before exit
		os.Exit(0)
	}()

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// Serve static files (these should be accessible without auth)
	javascriptFS, _ := fs.Sub(javascriptFiles, "javascript")
	stylesheetFS, _ := fs.Sub(stylesheetsFiles, "assets/stylesheets")
	imagesFS, _ := fs.Sub(imagesFiles, "assets/images")
	r.Handle("/javascript/*", http.StripPrefix("/javascript/", http.FileServer(http.FS(javascriptFS))))
	r.Handle("/assets/stylesheets/*", http.StripPrefix("/assets/stylesheets/", http.FileServer(http.FS(stylesheetFS))))
	r.Handle("/assets/images/*", http.StripPrefix("/assets/images/", http.FileServer(http.FS(imagesFS))))

	// Authentication routes (no auth required)
	r.Get("/signin", auth.ShowSignIn(templateFiles))
	r.Post("/auth/signin", auth.HandleSignIn())
	r.Get("/auth/signout", auth.HandleSignOut())

	// Protected routes (require authentication)
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuth)

		r.Get("/", home.Show(templateFiles))
		r.Get("/containers", containers.Index(templateFiles))
		r.Get("/logs/{containerID}", logs.Show(templateFiles))
		r.Get("/logs/stream/{containerID}", logs.Socket)
		r.Get("/terminal/{containerID}", terminal.Socket)
		r.Get("/terminal/view/{containerID}", terminal.Show(templateFiles))
		r.Get("/inspect/{containerID}", inspect.Show(templateFiles))
	})

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
