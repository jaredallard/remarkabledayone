// Copyright (C) 2024 Jared Allard
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package rm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Zip is a representation of the inside of a "rm" file version 6.
type Zip struct {
	// ID is the ID of the document contained in the zip file.
	ID string

	// Metadata is the contents of the "<id>.metadata" file.
	Metadata Metadata

	// Pages is a list of pages in the zip file.
	Pages []Page
}

type Metadata struct {
	CreatedTime    string `json:"createdTime"`
	LastModified   string `json:"lastModified"`
	LastOpened     string `json:"lastOpened"`
	LastOpenedPage int    `json:"lastOpenedPage"`
	Parent         string `json:"parent"`
	Pinned         bool   `json:"pinned"`
	Type           string `json:"type"`
	VisibleName    string `json:"visibleName"`
}

type Page struct {
	// ID is the ID of the page.
	ID string

	// Path is the path to the page. This is the "<id>.rm" file.
	Path string

	// PNGPath is the path to the rendered PNG file. To set, call "Render"
	// on the page.
	PNGPath string
}

// Render populates the PNGPath field of the page by rendering the page
// to a PNG file.
func (p *Page) Render() error {
	p.PNGPath = fmt.Sprintf("%s.png", strings.TrimSuffix(p.Path, ".rm"))
	return RenderRmToPng(p.Path, p.PNGPath)
}

// newZipFromDir creates a new Zip from a directory containing a
// Remarkable document.
func newZipFromDir(path string) (*Zip, error) {
	z := &Zip{
		Metadata: Metadata{},
		Pages:    make([]Page, 0),
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Find the document ID by looking for the first .metadata file.
	var metadataPath string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".metadata") {
			metadataPath = f.Name()
			break
		}
	}
	if metadataPath == "" {
		return nil, fmt.Errorf("no metadata file found")
	}

	id := strings.TrimSuffix(metadataPath, ".metadata")

	// Load the metadata file.
	f, err := os.Open(filepath.Join(path, metadataPath))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&z.Metadata); err != nil {
		return nil, err
	}

	// Find all the pages.
	pageFiles, err := os.ReadDir(filepath.Join(path, id))
	if err != nil {
		return nil, err
	}

	for _, f := range pageFiles {
		if !strings.HasSuffix(f.Name(), ".rm") {
			continue
		}

		z.Pages = append(z.Pages, Page{
			ID:   strings.TrimSuffix(f.Name(), ".rm"),
			Path: filepath.Join(path, id, f.Name()),
		})
	}

	return z, nil
}
