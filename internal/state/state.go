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

// Package state implements state tracking for the remarkabledayone
// utility.
package state

import (
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const FileName = "state.yml"

type State struct {
	// log is the logger for the state.
	log *slog.Logger `yaml:"-"`

	// path is the location where this state file was loaded from. If not
	// set, it shouldn't be saved.
	path string `yaml:"-"`

	// SyncedPages is a map of synced page IDs.
	SyncedPages map[string]struct{} `yaml:"synced_pages"`
}

// readStateFile reads the state file at the given path and returns the
// state. If an error occurs, it will return the error.
func readStateFile(path string) (*State, error) {
	// Found the state file, load it.
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	st := &State{}
	if err := yaml.NewDecoder(f).Decode(st); err != nil {
		return nil, err
	}

	return st, nil
}

// Load loads the state from disk. If an error occurs, it will return
// a new state.
func Load(log *slog.Logger) *State {
	var defaultState = &State{log: log}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.With("error", err).Error("failed to get user home directory")
		return &State{}
	}

	// nolint:errcheck // Why: Best effort to get the current working directory.
	cwd, _ := os.Getwd()

	// Load the state file from disk.
	var searchDirs = []string{
		// Use the XDG_STATE_HOME environment variable if it's set.
		func() string {
			stHome, ok := os.LookupEnv("XDG_STATE_HOME")
			if !ok {
				return ""
			}
			return filepath.Join(stHome, "remarkabledayone")
		}(),
		filepath.Join(homeDir, ".local", "state", "remarkabledayone"),
		cwd,
	}

	// Attempt each search directory.
	for _, dir := range searchDirs {
		if dir == "" {
			continue
		}

		// Save the first non-empty directory for saving if we created a new
		// state later.
		path := filepath.Join(dir, FileName)
		if defaultState.path == "" {
			defaultState.path = path
		}

		// Failed to read or doesn't exist, continue.
		st, err := readStateFile(path)
		if err != nil {
			continue
		}

		st.log = defaultState.log
		st.path = defaultState.path

		return st
	}

	// Didn't find the state file, return a new state.
	return defaultState
}

// Save saves the state to disk if there was a path provided when it was
// created. Top-level directories are created if they don't exist.
func (s *State) Save() error {
	if s.path == "" {
		s.log.Warn("not saving state, no path provided")
		return nil
	}

	s.log.Info("saving state", "path", s.path)

	// Ensure the directory exists.
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()

	return yaml.NewEncoder(f).Encode(s)
}
