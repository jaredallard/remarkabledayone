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

// Package implements a syncer between remarkable and dayone. It does
// not currently write back to Remarkable.
package syncer

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jaredallard/remarkabledayone/internal/config"
	"github.com/jaredallard/remarkabledayone/internal/dayone"
	"github.com/jaredallard/remarkabledayone/internal/rm"
	"github.com/jaredallard/remarkabledayone/internal/state"
	"github.com/juruen/rmapi/model"
)

// Syncer implements a syncer between remarkable and dayone. Create with
// the [New] function.
type Syncer struct {
	cfg   *config.Config
	log   *slog.Logger
	state *state.State
	rm    *rm.Client
}

// New creates a new syncer.
func New(log *slog.Logger, cfg *config.Config) (*Syncer, error) {
	st := state.Load(log.With("component", "state"))
	if st.SyncedPages == nil {
		st.SyncedPages = make(map[string]struct{})
	}

	//nolint:gocritic // Why: Acceptable shadow.
	rm, err := rm.New(log.With("component", "remarkable"))
	if err != nil {
		log.With("error", err).Error("failed to create remarkable client")
		os.Exit(1)
	}

	return &Syncer{
		cfg:   cfg,
		log:   log.With("component", "syncer"),
		state: st,
		rm:    rm,
	}, nil
}

// Sync syncs the configured document with DayOne.
func (s *Syncer) Sync() error {
	s.log.Info("syncing document", "name", s.cfg.DocumentName)
	var docMeta *model.Document
	for _, n := range s.rm.ListDocuments() {
		if n.Name() == s.cfg.DocumentName {
			docMeta = n.Document
			break
		}
	}
	if docMeta == nil {
		return fmt.Errorf("document not found")
	}

	doc, err := s.rm.DownloadDocument(docMeta)
	if err != nil {
		return fmt.Errorf("failed to download document: %w", err)
	}
	defer os.RemoveAll(doc.Path)

	s.log.Info("fetched document", "path", doc.Path, "pages", len(doc.Zip.Pages))

	// Compare the pages we have synced with the pages in the document.
	pagesHM := make(map[string]struct{})
	needToSync := make([]int, 0)
	for i, p := range doc.Zip.Pages {
		// Used for cleanup later.
		pagesHM[p.ID] = struct{}{}

		if _, ok := s.state.SyncedPages[p.ID]; ok {
			s.log.Debug("page already synced", "page", p.ID)
			continue
		}

		needToSync = append(needToSync, i)
	}

	// When we're done, cleanup the state.
	defer func() {
		// Remove pages that no longer exist from the state.
		s.log.Info("cleaning up pages in state that no longer exist")
		for id := range s.state.SyncedPages {
			if _, ok := pagesHM[id]; !ok {
				s.log.With("page", id).Info("removing page from state")
				delete(s.state.SyncedPages, id)
			}
		}
		if err := s.state.Save(); err != nil {
			s.log.Warn("failed to save state", "error", err)
		}
	}()

	if len(needToSync) == 0 {
		s.log.Info("no pages to sync")
		return nil
	}

	s.log.With("pages", len(needToSync)).Info("syncing pages")
	for _, p := range needToSync {
		page := doc.Zip.Pages[p]
		s.log.With("page", page.ID).Info("syncing page")

		// Render the page to a PNG.
		if err := page.Render(); err != nil {
			s.log.With("error", err).Error("failed to render page")
			continue
		}

		if err := dayone.EntryFromPNG(page.PNGPath, "Remarkable Entry", []string{"Remarkable"}); err != nil {
			s.log.With("error", err).Error("failed to create dayone entry")
			continue
		}

		s.state.SyncedPages[page.ID] = struct{}{}
	}

	if err := s.state.Save(); err != nil {
		s.log.Warn("failed to save state", "error", err)
	}

	s.log.With("pages", len(needToSync)).Info("synced pages")

	return nil
}
