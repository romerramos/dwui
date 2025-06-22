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
	"net/http"

	"github.com/dwui/cmd/session"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie(session.SessionCookieName)
		if err != nil || !session.Validate(cookie.Value) {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
