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
// GNU Affero General Public License for more.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package session

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/dwui/cmd/database"
)

type SessionData struct {
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}

var (
	SessionCookieName = "dwui_session"
)

func Create(w http.ResponseWriter) error {
	if database.Instance == nil {
		if err := database.Init(); err != nil {
			return err
		}
	}

	sessionID := make([]byte, 32)
	if _, err := rand.Read(sessionID); err != nil {
		return err
	}

	sessionToken := hex.EncodeToString(sessionID)
	expiry := time.Now().Add(24 * time.Hour) // 24 hour session

	if err := add(sessionToken, expiry); err != nil {
		return err
	}

	setCookie(w, sessionToken)
	return nil
}

func Validate(sessionToken string) bool {
	if database.Instance == nil {
		database.Init()
	}

	return validate(sessionToken)
}

func Clear(w http.ResponseWriter, sessionToken string) {
	if database.Instance == nil {
		database.Init()
	}

	remove(sessionToken)
	clearCookie(w)
}

func setCookie(w http.ResponseWriter, sessionToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   86400, // 24 hours
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})
}

func clearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func add(token string, expiry time.Time) error {
	if database.Instance == nil {
		if err := database.Init(); err != nil {
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

	return database.Instance.Update(func(txn *badger.Txn) error {
		// Set TTL for automatic expiration (BadgerInstance will clean up automatically)
		return txn.SetEntry(badger.NewEntry([]byte(token), sessionData).WithTTL(24 * time.Hour))
	})
}

func validate(token string) bool {
	if database.Instance == nil || token == "" {
		return false
	}

	var session SessionData
	found := false

	err := database.Instance.View(func(txn *badger.Txn) error {
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
		remove(token)
		return false
	}

	err = extend(session)
	return err == nil
}

func extend(session SessionData) error {
	if database.Instance == nil {
		return errors.New("database not initialized")
	}

	now := time.Now()

	newExpiry := now.Add(24 * time.Hour)
	session.Expiry = newExpiry

	sessionData, err := json.Marshal(session)
	if err != nil {
		return err
	}

	database.Instance.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(badger.NewEntry([]byte(session.Token), sessionData).WithTTL(24 * time.Hour))
	})

	return nil
}

// remove removes a session from BadgerInstance
func remove(token string) {
	if database.Instance == nil {
		return
	}

	database.Instance.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(token))
	})
}
