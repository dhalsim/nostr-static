package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"nostr-static/src/commands"
	"nostr-static/src/helpers"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:           "nostr-static",
		Usage:          "Generate static HTML pages from Nostr articles",
		DefaultCommand: "generate",
		Commands: []*cli.Command{
			{
				Name:  "generate",
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
					&cli.StringFlag{
						Name:    "index-dir",
						Value:   "index-tag-discovery",
						Usage:   "directory to save the tag discovery index",
						Aliases: []string{"i"},
					},
					&cli.BoolFlag{
						Name:    "clean",
						Value:   true,
						Usage:   "clean output directory before generating",
						Aliases: []string{"C"},
					},
					&cli.BoolFlag{
						Name:    "save-files",
						Value:   false,
						Usage:   "save events and profiles to JSON files",
						Aliases: []string{"s"},
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					configPath := cmd.String("config")
					outputDir := cmd.String("output")
					indexDir := cmd.String("index-dir")
					clean := cmd.Bool("clean")
					saveFiles := cmd.Bool("save-files")

					// Load configuration
					log.Println("loading config from: ", configPath)

					config, err := LoadConfig(configPath)
					if err != nil {
						return err
					}

					// Fetch events
					log.Println("fetching events from relays: ", config.Relays)

					events, eventIDToNaddr, err := helpers.FetchEvents(
						config.Relays,
						config.Articles,
					)
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

					// Copy static folder to output directory
					log.Println("copying static files to output directory")
					if err := helpers.CopyDir("src/static", filepath.Join(outputDir, "static")); err != nil {
						return fmt.Errorf("failed to copy static files: %w", err)
					}

					// Save events to JSON file
					if saveFiles {
						eventsPath := filepath.Join(outputDir, "events.json")
						if err := helpers.SaveEvents(events, eventsPath); err != nil {
							return err
						}

						// Save naddr mapping
						naddrPath := filepath.Join(outputDir, "naddr.json")
						if err := helpers.SaveNaddrMapping(eventIDToNaddr, naddrPath); err != nil {
							return err
						}
					}

					// Fetch profiles
					pubkeyToNprofile := make(map[string]string)
					pubkeys := make([]string, len(config.Profiles))

					for i, nprofile := range config.Profiles {
						prefix, data, err := nip19.Decode(nprofile)
						if err != nil {
							return err
						}

						if prefix != "nprofile" {
							return fmt.Errorf("invalid nprofile prefix: %s", prefix)
						}

						profilePointer := data.(nostr.ProfilePointer)
						pubkeyToNprofile[profilePointer.PublicKey] = nprofile
						pubkeys[i] = profilePointer.PublicKey
					}

					log.Println("fetching profiles from relays: ", config.Relays)

					pubkeyToKind0, err := helpers.FetchProfiles(config.Relays, pubkeys)
					if err != nil {
						return err
					}

					if saveFiles {
						// Save profiles to JSON file
						profilesPath := filepath.Join(outputDir, "profiles.json")
						if err := helpers.SaveProfiles(pubkeyToKind0, profilesPath); err != nil {
							return err
						}

						// Save nprofile mapping
						nprofilePath := filepath.Join(outputDir, "nprofile.json")
						if err := helpers.SaveNprofileMapping(pubkeyToNprofile, nprofilePath); err != nil {
							return err
						}
					}

					return commands.Generate(commands.GenerateCommandParams{
						ConfigPath:       configPath,
						OutputDir:        outputDir,
						IndexDir:         indexDir,
						Config:           config,
						PubkeyToKind0:    pubkeyToKind0,
						PubkeyToNprofile: pubkeyToNprofile,
						Events:           events,
						EventIDToNaddr:   eventIDToNaddr,
					})
				},
			},
			{
				Name:  "index-tag-discovery",
				Usage: "Index tag discovery",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Value:   "config.yaml",
						Usage:   "path to config file",
						Aliases: []string{"c"},
					},
					&cli.BoolFlag{
						Name:    "reset",
						Value:   false,
						Usage:   "reset the tag discovery index",
						Aliases: []string{"r"},
					},
					&cli.StringFlag{
						Name:    "index-dir",
						Value:   "index-tag-discovery",
						Usage:   "directory to save the tag discovery index",
						Aliases: []string{"i"},
					},
					&cli.StringFlag{
						Name:    "output-dir",
						Value:   "output",
						Usage:   "output directory for generated files",
						Aliases: []string{"o"},
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					reset := cmd.Bool("reset")
					indexDir := cmd.String("index-dir")
					outputDir := cmd.String("output-dir")
					configPath := cmd.String("config")

					config, err := LoadConfig(configPath)
					if err != nil {
						return err
					}

					return commands.IndexTagDiscovery(commands.IndexTagDiscoveryCommandParams{
						Reset:     reset,
						IndexDir:  indexDir,
						OutputDir: outputDir,
						Config:    config,
					})
				},
			},
			{
				Name:  "calculate-tag-discovery",
				Usage: "Calculate tag discovery",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Value:   "config.yaml",
						Usage:   "path to config file",
						Aliases: []string{"c"},
					},
					&cli.StringFlag{
						Name:    "index-dir",
						Value:   "index-tag-discovery",
						Usage:   "directory to save the tag discovery index",
						Aliases: []string{"i"},
					},
					&cli.StringFlag{
						Name:    "output-dir",
						Value:   "output",
						Usage:   "output directory for generated files",
						Aliases: []string{"o"},
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					indexDir := cmd.String("index-dir")
					outputDir := cmd.String("output-dir")
					configPath := cmd.String("config")

					config, err := LoadConfig(configPath)
					if err != nil {
						return err
					}

					return commands.CalculateTagDiscovery(commands.CalculateTagDiscoveryCommandParams{
						IndexDir:  indexDir,
						OutputDir: outputDir,
						Config:    config,
					})
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
