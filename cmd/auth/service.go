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
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

var (
	passwordHash string
)

func ValidatePassword(password string) bool {
	hash := sha256.Sum256([]byte(password))
	providedHash := hex.EncodeToString(hash[:])
	return providedHash == passwordHash
}

func SetPassword(password string) {
	hash := sha256.Sum256([]byte(password))
	passwordHash = hex.EncodeToString(hash[:])
}

func GenerateRandomPassword() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	password := fmt.Sprintf("%x", bytes)[:12]
	return password, nil
}
