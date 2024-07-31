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

// Config contains configuration for the syncer.
package config

import (
	"log/slog"
	"os"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	// DocumentName is the name of the document to sync.
	DocumentName string `env:"DOCUMENT_NAME,required"`
}

func Load(log *slog.Logger) (*Config, error) {
	environment := strings.ToLower(os.Getenv("ENV"))

	var envFile string
	switch environment {
	case "dev", "development":
		envFile = ".env.development"
	case "prod", "production":
		envFile = ".env.production"
	default:
		envFile = ".env"
	}

	// If there's an environment file, load it.
	if envFile != "" {
		if _, err := os.Stat(envFile); err == nil {
			if err := godotenv.Load(envFile); err != nil {
				return nil, err
			}
		} else if !os.IsNotExist(err) {
			return nil, err
		}
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
