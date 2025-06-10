package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

var (
	// Global password hash for comparison
	passwordHash string
	// Simple session store (in production, use Redis or database)
	sessions = make(map[string]time.Time)
	// Session cookie name
	sessionCookieName = "dwui_session"
)

// SetPassword sets the application password
func SetPassword(password string) {
	hash := sha256.Sum256([]byte(password))
	passwordHash = hex.EncodeToString(hash[:])
}

// GenerateRandomPassword generates a random password
func GenerateRandomPassword() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert to a more readable format
	password := fmt.Sprintf("%x", bytes)[:12] // 12 character password
	return password, nil
}

// ValidatePassword checks if the provided password is correct
func ValidatePassword(password string) bool {
	hash := sha256.Sum256([]byte(password))
	providedHash := hex.EncodeToString(hash[:])
	return providedHash == passwordHash
}

// CreateSession creates a new session for the user
func CreateSession() (string, error) {
	sessionID := make([]byte, 32)
	if _, err := rand.Read(sessionID); err != nil {
		return "", err
	}

	sessionToken := hex.EncodeToString(sessionID)
	sessions[sessionToken] = time.Now().Add(24 * time.Hour) // 24 hour session

	return sessionToken, nil
}

// ValidateSession checks if a session is valid
func ValidateSession(sessionToken string) bool {
	if sessionToken == "" {
		return false
	}

	expiry, exists := sessions[sessionToken]
	if !exists {
		return false
	}

	if time.Now().After(expiry) {
		delete(sessions, sessionToken)
		return false
	}

	// Extend session
	sessions[sessionToken] = time.Now().Add(24 * time.Hour)
	return true
}

// ClearSession removes a session
func ClearSession(sessionToken string) {
	delete(sessions, sessionToken)
}

// RequireAuth middleware that checks for valid authentication
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for sign-in page and static files
		if r.URL.Path == "/signin" || r.URL.Path == "/auth/signin" ||
			r.URL.Path == "/auth/signout" ||
			isStaticFile(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Check for session cookie
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil || !ValidateSession(cookie.Value) {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// isStaticFile checks if the request is for a static file
func isStaticFile(path string) bool {
	staticPaths := []string{"/javascript/", "/assets/"}
	for _, staticPath := range staticPaths {
		if len(path) >= len(staticPath) && path[:len(staticPath)] == staticPath {
			return true
		}
	}
	return false
}

// SetSessionCookie sets the session cookie
func SetSessionCookie(w http.ResponseWriter, sessionToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   86400, // 24 hours
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearSessionCookie removes the session cookie
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}
