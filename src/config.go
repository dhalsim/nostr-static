package main

import (
	"os"

	"gopkg.in/yaml.v3"

	"nostr-static/src/types"
)

func LoadConfig(path string) (*types.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config types.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Set default theme if not specified
	if config.Layout.Color == "" {
		config.Layout.Color = "light"
	}

	// Set default features if not specified
	if config.Features.Comments {
		config.Features.Comments = true
	}

	return &config, nil
}
