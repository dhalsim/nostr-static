package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"nostr-static/src/pagegenerators"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
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
			&cli.BoolFlag{
				Name:    "clean",
				Value:   true,
				Usage:   "clean output directory before generating",
				Aliases: []string{"C"},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configPath := cmd.String("config")
			outputDir := cmd.String("output")
			clean := cmd.Bool("clean")

			// Load configuration
			log.Println("loading config from: ", configPath)

			config, err := LoadConfig(configPath)
			if err != nil {
				return err
			}

			if clean {
				if err := os.RemoveAll(outputDir); err != nil {
					return err
				}
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

			log.Println("fetching events from relays: ", config.Relays)

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

			log.Println("fetching profiles for all authors")

			pubkeys := make([]string, 0)
			pubkeyToNprofile := make(map[string]string)

			for _, nprofile := range config.Profiles {
				prefix, pubkey, err := nip19.Decode(nprofile)
				if err != nil {
					return err
				}

				if prefix != "nprofile" {
					return fmt.Errorf("invalid nprofile prefix: %s", prefix)
				}

				profilePointer := pubkey.(nostr.ProfilePointer)
				pubkeys = append(pubkeys, profilePointer.PublicKey)
				pubkeyToNprofile[profilePointer.PublicKey] = nprofile
			}

			// Fetch profiles for all authors
			profiles, err := fetchProfiles(config.Relays, pubkeys)
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

			log.Println("generating article pages")

			// Article pages
			for _, event := range events {
				params := pagegenerators.NewGenerateArticleParams(
					event,
					outputDir,
					config.Layout,
					config.Features,
					eventIDToNaddr[event.ID],
					profiles[event.PubKey],
					pubkeyToNprofile[event.PubKey],
					config.Relays,
				)

				if err := pagegenerators.GenerateArticleHTML(params); err != nil {
					log.Printf("Failed to generate HTML for event %s: %v", event.ID, err)
					continue
				}
			}

			log.Println("generating tag pages")

			// Tag pages
			if err := pagegenerators.GenerateTagPages(pagegenerators.NewGenerateTagPagesParams(
				events,
				profiles,
				outputDir,
				config.Layout,
				eventIDToNaddr,
				pubkeyToNprofile,
			)); err != nil {
				return err
			}

			log.Println("generating profile pages")

			// Profile pages
			if err := pagegenerators.GenerateProfilePages(pagegenerators.NewGenerateProfilePagesParams(
				profiles,
				events,
				outputDir,
				config.Layout,
				pubkeyToNprofile,
				eventIDToNaddr,
			)); err != nil {
				return err
			}

			log.Println("generating index page")

			// Index page
			if err := pagegenerators.GenerateIndexHTML(pagegenerators.NewGenerateIndexParams(
				events,
				profiles,
				outputDir,
				config.Layout,
				eventIDToNaddr,
				pubkeyToNprofile,
			)); err != nil {
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
