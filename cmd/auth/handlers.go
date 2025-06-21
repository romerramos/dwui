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

package auth

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/dwui/cmd/session"
)

type SignInData struct {
	Error string
}

func ShowSignIn(templateFiles embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(templateFiles, "cmd/auth/signin.gohtml")
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}

		data := SignInData{}

		if errorMsg := r.URL.Query().Get("error"); errorMsg != "" {
			data.Error = "Invalid password. Please try again."
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
		}
	}
}

func HandleSignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		password := r.FormValue("password")
		if !ValidatePassword(password) {
			http.Redirect(w, r, "/signin?error=invalid", http.StatusSeeOther)
			return
		}

		err := session.Create(w)
		if err != nil {
			log.Printf("failed to create session: %v", err)
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func HandleSignOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie(session.SessionCookieName); err == nil {
			session.Clear(w, cookie.Value)
		}
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	}
}
