// DWUI (Docker Web UI)
// Copyright (C) 2025 Romer Ramos
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

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
	"github.com/dwui/cmd/database"
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
	var port string
	var passwordFile string
	flag.StringVar(&password, "password", "", "Password for authentication (if not provided, a random one will be generated)")
	flag.StringVar(&port, "port", "8300", "Port to run the server on")
	flag.StringVar(&passwordFile, "password-file", "", "File to store the generated password")
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
		fmt.Printf("💡 Use this password to sign in to the web interface\n\n")

		if passwordFile != "" {
			err := os.WriteFile(passwordFile, []byte(password), 0600)
			if err != nil {
				fmt.Printf("⚠️  Warning: Could not write password to %s: %v\n", passwordFile, err)
			} else {
				fmt.Printf("🔑 Password also stored in: %s\n", passwordFile)
			}
		}
	} else {
		fmt.Printf("🔐 Using provided password for authentication\n")
	}

	database.Init()
	auth.SetPassword(password)

	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\n🛑 Gracefully shutting down...")
		database.Instance.Close() // Close BadgerDB before exit
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
		r.Get("/terminal/{containerID}", terminal.Show(templateFiles))
		r.Get("/terminal/stream/{containerID}", terminal.Socket)
		r.Get("/inspect/{containerID}", inspect.Show(templateFiles))
	})

	fmt.Printf("Starting server on :%s\n", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
