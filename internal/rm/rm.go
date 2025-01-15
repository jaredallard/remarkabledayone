// Copyright (C) 2025 remarkabledayone contributors
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

// Package rm contains logic for interacting with the Remarkable API for
// the purposes of reading and writing documents.
package rm

import (
	"archive/zip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/juruen/rmapi/api"
	"github.com/juruen/rmapi/model"
	"github.com/juruen/rmapi/transport"
)

type Client struct {
	log  *slog.Logger
	rm   api.ApiCtx
	user *api.UserInfo
}

type Document struct {
	Path string
	Zip  *Zip
}

// remarkableAuth attempts to authenticate with remarkable's API.
func remarkableAuth(_ *slog.Logger) (api.ApiCtx, *api.UserInfo, error) {
	for range 3 {
		ctx, userInf, err := apiCtxForTransport(api.AuthHttpCtx(true, false))
		if err != nil {
			log.With("error", err).Error("failed to authenticate")
			continue
		}

		return ctx, userInf, nil
	}

	// If we've reached this point, we've failed to authenticate.
	return nil, nil, fmt.Errorf("failed to authenticate")
}

// apiCtxForTransport creates an API context for the given transport.
func apiCtxForTransport(tctx *transport.HttpClientCtx) (api.ApiCtx, *api.UserInfo, error) {
	userInfo, err := api.ParseToken(tctx.Tokens.UserToken)
	if err != nil {
		return nil, nil, err
	}

	apiCtx, err := api.CreateApiCtx(tctx, userInfo.SyncVersion)
	if err != nil {
		return nil, nil, err
	}

	return apiCtx, userInfo, nil
}

// New creates a new Client for interacting with the Remarkable API.
//
//nolint:gocritic // Why: Acceptable shadow.
func New(log *slog.Logger) (*Client, error) {
	ctx, userInfo, err := remarkableAuth(log)
	if err != nil {
		return nil, err
	}

	return &Client{log, ctx, userInfo}, nil
}

// ListDocuments returns a list of all documents on the Remarkable.
func (c *Client) ListDocuments() []*model.Node {
	mDocs := c.rm.Filetree().Root().Children

	// Convert to []*model.Node
	docs := make([]*model.Node, 0, len(mDocs))
	for i := range mDocs {
		docs = append(docs, mDocs[i])
	}

	// Sort by name
	sort.Slice(docs, func(i, j int) bool {
		return docs[i].Name() < docs[j].Name()
	})

	return docs
}

// sanitizeArchivePath to mitigate "G305".
func sanitizeArchivePath(d, t string) (v string, err error) {
	v = filepath.Join(d, t)
	if strings.HasPrefix(v, filepath.Clean(d)) {
		return v, nil
	}

	return "", fmt.Errorf("%s: %s", "content filepath is tainted", t)
}

func (c *Client) zipFromArchive(tmpDir, path string) (*Zip, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fInf, err := f.Stat()
	if err != nil {
		return nil, err
	}

	zr, err := zip.NewReader(f, fInf.Size())
	if err != nil {
		return nil, err
	}

	for _, f := range zr.File {
		if err := func() error {
			zf, err := f.Open()
			if err != nil {
				return err
			}
			defer zf.Close()

			// Skip directories.
			if f.FileInfo().IsDir() {
				return nil
			}

			// Create the file in the temporary directory.
			outPath, err := sanitizeArchivePath(tmpDir, f.Name)
			if err != nil {
				return err
			}

			// Ensure the directory exists.
			if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
				return err
			}

			out, err := os.Create(outPath)
			if err != nil {
				return err
			}
			defer out.Close()

			// Copy the file.
			//#nosec:G110 // Why: This is acceptable for our use case.
			if _, err := io.Copy(out, zf); err != nil {
				return err
			}

			return nil
		}(); err != nil {
			return nil, err
		}
	}

	return newZipFromDir(tmpDir)
}

// DownloadDocument downloads a document from Remarkable and returns it.
func (c *Client) DownloadDocument(doc *model.Document) (*Document, error) {
	tmpDir, err := os.MkdirTemp("", "remarkabledayone")
	if err != nil {
		log.With("error", err).Error("failed to create temporary directory")
		os.Exit(1)
	}
	tmpFile := filepath.Join(tmpDir, strings.TrimSuffix(doc.Name, ".zip")+".zip")

	if err := c.rm.FetchDocument(doc.ID, tmpFile); err != nil {
		return nil, err
	}

	c.log.Info("downloaded document", "name", doc.Name, "path", tmpFile)
	z, err := c.zipFromArchive(tmpDir, tmpFile)
	if err != nil {
		return nil, err
	}

	return &Document{tmpFile, z}, nil
}
