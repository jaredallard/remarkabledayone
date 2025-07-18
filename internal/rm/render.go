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

package rm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jaredallard/cmdexec"
)

// RenderRmToPng renders a remarkable document to a PNG file.
func RenderRmToPng(src, dest string) error {
	tmpDir, err := os.MkdirTemp("", "rm-render-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir) //nolint:errcheck // Why: Best effort.

	pdfFile := filepath.Join(tmpDir, "output.pdf")
	outputFile := pdfFile + ".png"
	cmds := [][]string{
		{"rmc", "--version"},
		{"rmc", "-t", "pdf", src, "-o", pdfFile},
		{"convert", "-verbose", "-density", "150", "-trim", pdfFile, "-quality", "100", "-flatten", "-sharpen", "0x1.0", outputFile},
	}
	for _, cmd := range cmds {
		cmd := cmdexec.Command(cmd[0], cmd[1:]...)
		cmd.SetStdout(os.Stdout)
		cmd.SetStderr(os.Stderr)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run command %v: %w", cmd, err)
		}
	}

	return os.Rename(outputFile, dest)
}
