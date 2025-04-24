package main

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"nostr-static/src/pagegenerators"

	"github.com/urfave/cli/v3"
)

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func main() {
	app := &cli.Command{
		Name:  "nostr-static",
		Usage: "Generate static HTML pages from Nostr articles",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Value:   "config.yaml",
				Usage:   "path to config file",
				Aliases: []string{"c"},
			},
			&cli.StringFlag{
				Name:    "output",
				Value:   "output",
				Usage:   "output directory for generated files",
				Aliases: []string{"o"},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configPath := cmd.String("config")
			outputDir := cmd.String("output")

			// Load configuration
			config, err := LoadConfig(configPath)
			if err != nil {
				return err
			}

			// Create output directory if it doesn't exist
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return err
			}

			// Copy logo file if specified
			if config.Layout.Logo != "" {
				logoPath := filepath.Join(filepath.Dir(configPath), config.Layout.Logo)
				outputLogoPath := filepath.Join(outputDir, config.Layout.Logo)
				if err := copyFile(logoPath, outputLogoPath); err != nil {
					log.Printf("Warning: Failed to copy logo file: %v", err)
				}
			}

			// Fetch events from relays
			events, err := fetchEvents(config.Relays, config.ArticleIDs)
			if err != nil {
				return err
			}

			// Save events to JSON file
			eventsPath := filepath.Join(outputDir, "events.json")
			if err := saveEvents(events, eventsPath); err != nil {
				return err
			}

			// Generate HTML files for each event
			for _, event := range events {
				if err := pagegenerators.GenerateArticleHTML(event, outputDir, config.Layout); err != nil {
					log.Printf("Failed to generate HTML for event %s: %v", event.ID, err)
					continue
				}
			}

			// Generate tag pages
			if err := pagegenerators.GenerateTagPages(events, outputDir, config.Layout); err != nil {
				return err
			}

			// Generate index.html
			if err := pagegenerators.GenerateIndexHTML(events, outputDir, config.Layout); err != nil {
				return err
			}

			log.Printf("Successfully generated static site in %s", outputDir)
			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
