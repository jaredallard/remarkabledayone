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

// Package main implements the entrypoint for the remarkabledayone
// utility.
package main

import (
	"log/slog"
	"os"

	charmlog "github.com/charmbracelet/log"
	"github.com/jaredallard/remarkabledayone/internal/config"
	"github.com/jaredallard/remarkabledayone/internal/syncer"
)

// main is the entrypoint for the remarkabledayone utility.
func main() {
	handler := charmlog.New(os.Stderr)
	log := slog.New(handler)

	cfg, err := config.Load(log.With("component", "config"))
	if err != nil {
		log.With("error", err).Error("failed to load config")
		os.Exit(1)
	}

	if os.Getenv("ENV") == "development" {
		handler.SetLevel(charmlog.DebugLevel)
		log.Debug("debug logging enabled")
	}

	syncer, err := syncer.New(log, cfg)
	if err != nil {
		log.With("error", err).Error("failed to create syncer")
		os.Exit(1)
	}

	if err := syncer.Sync(); err != nil {
		log.With("error", err).Error("failed to sync")
		os.Exit(1)
	}
}
