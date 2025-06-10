package auth

import (
	"embed"
	"html/template"
	"net/http"
)

// SignInData holds data for the sign-in template
type SignInData struct {
	PageTitle string
	Error     string
}

// ShowSignIn displays the sign-in page
func ShowSignIn(templateFiles embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(templateFiles, "cmd/auth/signin.gohtml")
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}

		data := SignInData{
			PageTitle: "Docker Web UI - Sign In",
		}

		// Check for error in query params
		if errorMsg := r.URL.Query().Get("error"); errorMsg != "" {
			data.Error = "Invalid password. Please try again."
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
		}
	}
}

// HandleSignIn processes the sign-in form submission
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

		// Create session
		sessionToken, err := CreateSession()
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		// Set session cookie
		SetSessionCookie(w, sessionToken)

		// Redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// HandleSignOut processes sign-out requests
func HandleSignOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session cookie
		if cookie, err := r.Cookie(sessionCookieName); err == nil {
			ClearSession(cookie.Value)
		}

		// Clear session cookie
		ClearSessionCookie(w)

		// Redirect to sign-in page
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	}
}
