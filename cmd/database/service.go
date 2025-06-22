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

package database

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v4"
)

var (
	Instance *badger.DB
)

func Init() error {
	if Instance != nil {
		return nil
	}

	tempDir := os.TempDir()
	dbDir := filepath.Join(tempDir, "dwui_sessions")

	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create sessions directory: %v", err)
	}

	opts := badger.DefaultOptions(dbDir)
	opts.Logger = nil

	var err error
	Instance, err = badger.Open(opts)
	if err != nil {
		return fmt.Errorf("failed to open BadgerInstance: %v", err)
	}

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			Instance.RunValueLogGC(0.5)
		}
	}()

	return nil
}

// CloseInstance closes the BadgerInstance instance
func CloseInstance() {
	if Instance != nil {
		Instance.Close()
	}
}
