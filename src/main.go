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
			events, eventIDToNaddr, err := fetchEvents(config.Relays, config.Articles)
			if err != nil {
				return err
			}

			// Save events to JSON file
			eventsPath := filepath.Join(outputDir, "events.json")
			if err := saveEvents(events, eventsPath); err != nil {
				return err
			}

			// Fetch profiles for all authors
			profiles, pubkeyToNprofile, err := fetchProfiles(config.Relays, events)
			if err != nil {
				return err
			}

			// Save profiles to JSON file
			profilesPath := filepath.Join(outputDir, "profiles.json")
			if err := saveProfiles(profiles, profilesPath); err != nil {
				return err
			}

			// Save nprofile mapping
			nprofilePath := filepath.Join(outputDir, "nprofile.json")
			if err := saveNprofileMapping(pubkeyToNprofile, nprofilePath); err != nil {
				return err
			}

			// Save naddr mapping
			naddrPath := filepath.Join(outputDir, "naddr.json")
			if err := saveNaddrMapping(eventIDToNaddr, naddrPath); err != nil {
				return err
			}

			// Article pages
			for _, event := range events {
				params := pagegenerators.GenerateArticleParams{
					Event:     event,
					OutputDir: outputDir,
					Layout:    config.Layout,
					Naddr:     eventIDToNaddr[event.ID],
					Nprofile:  pubkeyToNprofile[event.PubKey],
					Profile:   profiles[event.PubKey],
				}

				if err := pagegenerators.GenerateArticleHTML(params); err != nil {
					log.Printf("Failed to generate HTML for event %s: %v", event.ID, err)
					continue
				}
			}

			// Tag pages
			if err := pagegenerators.GenerateTagPages(pagegenerators.GenerateTagPagesParams{
				Events:           events,
				Profiles:         profiles,
				OutputDir:        outputDir,
				Layout:           config.Layout,
				EventIDToNaddr:   eventIDToNaddr,
				PubkeyToNProfile: pubkeyToNprofile,
			}); err != nil {
				return err
			}

			// Profile pages
			if err := pagegenerators.GenerateProfilePages(pagegenerators.GenerateProfilePagesParams{
				Profiles:         profiles,
				Events:           events,
				OutputDir:        outputDir,
				Layout:           config.Layout,
				EventIDToNaddr:   eventIDToNaddr,
				PubkeyToNProfile: pubkeyToNprofile,
			}); err != nil {
				return err
			}

			// Index page
			if err := pagegenerators.GenerateIndexHTML(pagegenerators.GenerateIndexParams{
				Events:           events,
				Profiles:         profiles,
				OutputDir:        outputDir,
				Layout:           config.Layout,
				EventIDToNaddr:   eventIDToNaddr,
				PubkeyToNProfile: pubkeyToNprofile,
			}); err != nil {
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
