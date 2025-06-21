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

package containers

func ShortenID(id string) string {
	return shorten(id, 12)
}

func ShortenName(name string) string {
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}
	return shorten(name, 25)
}

func shorten(text string, amount int) string {
	if len(text) > amount {
		return text[:amount]
	}
	return text
}
