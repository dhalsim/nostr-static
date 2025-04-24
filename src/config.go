package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Layout struct {
	Color string `yaml:"color"`
	Logo  string `yaml:"logo"`
}

type Config struct {
	Relays     []string `yaml:"relays"`
	ArticleIDs []string `yaml:"article_ids"`
	Layout     Layout   `yaml:"layout"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Set default theme if not specified
	if config.Layout.Color == "" {
		config.Layout.Color = "light"
	}

	return &config, nil
}
