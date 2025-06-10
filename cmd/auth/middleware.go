package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v4"
)

// SessionData represents a session with its expiry time
type SessionData struct {
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}

var (
	// Global password hash for comparison
	passwordHash string
	// BadgerDB instance for session storage
	db *badger.DB
	// Session cookie name
	sessionCookieName = "dwui_session"
)

// initDB initializes the BadgerDB instance for session storage
func initDB() error {
	if db != nil {
		return nil // Already initialized
	}

	// Create sessions directory in temp folder
	tempDir := os.TempDir()
	dbDir := filepath.Join(tempDir, "dwui_sessions")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create sessions directory: %v", err)
	}

	// Open BadgerDB
	opts := badger.DefaultOptions(dbDir)
	opts.Logger = nil // Disable BadgerDB logging for cleaner output

	var err error
	db, err = badger.Open(opts)
	if err != nil {
		return fmt.Errorf("failed to open BadgerDB: %v", err)
	}

	// Start background goroutine for garbage collection
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			// Run garbage collection
			db.RunValueLogGC(0.5)
			// Clean up expired sessions
			cleanupExpiredSessions()
		}
	}()

	return nil
}

// cleanupExpiredSessions removes expired sessions from BadgerDB
func cleanupExpiredSessions() {
	if db == nil {
		return
	}

	now := time.Now()
	var expiredKeys [][]byte

	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()

			err := item.Value(func(val []byte) error {
				var session SessionData
				if err := json.Unmarshal(val, &session); err != nil {
					// Invalid session data, mark for deletion
					expiredKeys = append(expiredKeys, append([]byte(nil), key...))
					return nil
				}

				if now.After(session.Expiry) {
					expiredKeys = append(expiredKeys, append([]byte(nil), key...))
				}
				return nil
			})
			if err != nil {
				continue
			}
		}
		return nil
	})

	if err != nil {
		return
	}

	// Delete expired sessions
	if len(expiredKeys) > 0 {
		db.Update(func(txn *badger.Txn) error {
			for _, key := range expiredKeys {
				txn.Delete(key)
			}
			return nil
		})
	}
}

// addSession adds a new session to BadgerDB
func addSession(token string, expiry time.Time) error {
	if db == nil {
		if err := initDB(); err != nil {
			return err
		}
	}

	session := SessionData{
		Token:  token,
		Expiry: expiry,
	}

	sessionData, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return db.Update(func(txn *badger.Txn) error {
		// Set TTL for automatic expiration (BadgerDB will clean up automatically)
		return txn.SetEntry(badger.NewEntry([]byte(token), sessionData).WithTTL(24 * time.Hour))
	})
}

// validateSession checks if a session is valid and extends it
func validateSession(token string) bool {
	if db == nil || token == "" {
		return false
	}

	var session SessionData
	found := false

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(token))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			found = true
			return json.Unmarshal(val, &session)
		})
	})

	if err != nil || !found {
		return false
	}

	now := time.Now()
	if now.After(session.Expiry) {
		// Session expired, remove it
		removeSession(token)
		return false
	}

	// Extend session
	newExpiry := now.Add(24 * time.Hour)
	session.Expiry = newExpiry

	sessionData, err := json.Marshal(session)
	if err != nil {
		return false
	}

	// Update session with new expiry
	db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(token), sessionData)
	})

	return true
}

// removeSession removes a session from BadgerDB
func removeSession(token string) {
	if db == nil {
		return
	}

	db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(token))
	})
}

// CloseDB closes the BadgerDB instance
func CloseDB() {
	if db != nil {
		db.Close()
	}
}

// SetPassword sets the application password and initializes BadgerDB
func SetPassword(password string) {
	hash := sha256.Sum256([]byte(password))
	passwordHash = hex.EncodeToString(hash[:])

	// Initialize BadgerDB if not already done
	if db == nil {
		initDB()
	}
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
	// Initialize BadgerDB if not already done
	if db == nil {
		if err := initDB(); err != nil {
			return "", err
		}
	}

	sessionID := make([]byte, 32)
	if _, err := rand.Read(sessionID); err != nil {
		return "", err
	}

	sessionToken := hex.EncodeToString(sessionID)
	expiry := time.Now().Add(24 * time.Hour) // 24 hour session

	if err := addSession(sessionToken, expiry); err != nil {
		return "", err
	}

	return sessionToken, nil
}

// ValidateSession checks if a session is valid
func ValidateSession(sessionToken string) bool {
	// Initialize BadgerDB if not already done
	if db == nil {
		initDB()
	}

	return validateSession(sessionToken)
}

// ClearSession removes a session
func ClearSession(sessionToken string) {
	// Initialize BadgerDB if not already done
	if db == nil {
		initDB()
	}

	removeSession(sessionToken)
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
