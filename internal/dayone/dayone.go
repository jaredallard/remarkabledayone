// Copyright (C) 2024 remarkabledayone contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// SPDX-License-Identifier: AGPL-3.0

// Package dayone implements a small wrapper around the dayone2 CLI for
// creating entries in DayOne.
package dayone

import (
	"os"
	"os/exec"
)

// EntryFromPNG creates a new DayOne entry from a PNG file.
func EntryFromPNG(src, title string, tags []string) error {
	args := []string{"--attachments", src}

	if len(tags) > 0 {
		args = append(args, "--tags")
		args = append(args, tags...)
		args = append(args, "--")
	}

	// Add the new, title and attachment arguments.
	args = append(args, "new", title, "[{attachment}]")

	//#nosec:G204 // Why: Safe for our usecase.
	cmd := exec.Command("dayone2", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
