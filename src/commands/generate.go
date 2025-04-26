package commands

import (
	"log"
	"path/filepath"

	"nostr-static/src/helpers"
	"nostr-static/src/pagegenerators"
	"nostr-static/src/types"
)

type generateCommandParams struct {
	configPath       string
	outputDir        string
	config           *types.Config
	pubkeyToKind0    map[string]types.Event
	pubkeyToNprofile map[string]string
	events           []types.Event
	eventIDToNaddr   map[string]string
}

func NewGenerateCommandParams(
	configPath string,
	outputDir string,
	config *types.Config,
	pubkeyToKind0 map[string]types.Event,
	pubkeyToNprofile map[string]string,
	events []types.Event,
	eventIDToNaddr map[string]string,
) generateCommandParams {
	return generateCommandParams{
		configPath:       configPath,
		outputDir:        outputDir,
		config:           config,
		pubkeyToKind0:    pubkeyToKind0,
		pubkeyToNprofile: pubkeyToNprofile,
		events:           events,
		eventIDToNaddr:   eventIDToNaddr,
	}
}

func Generate(params generateCommandParams) error {
	configPath := params.configPath
	outputDir := params.outputDir
	config := params.config
	events := params.events
	eventIDToNaddr := params.eventIDToNaddr
	pubkeyToKind0 := params.pubkeyToKind0
	pubkeyToNprofile := params.pubkeyToNprofile

	// Copy logo file if specified
	if config.Layout.Logo != "" {
		logoPath := filepath.Join(filepath.Dir(configPath), config.Layout.Logo)
		outputLogoPath := filepath.Join(outputDir, config.Layout.Logo)
		if err := helpers.CopyFile(logoPath, outputLogoPath); err != nil {
			log.Printf("Warning: Failed to copy logo file: %v", err)
		}
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
			pubkeyToKind0[event.PubKey],
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
		pubkeyToKind0,
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
		pubkeyToKind0,
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
		pubkeyToKind0,
		outputDir,
		config.Layout,
		eventIDToNaddr,
		pubkeyToNprofile,
	)); err != nil {
		return err
	}

	log.Printf("Successfully generated static site in %s", outputDir)
	return nil
}
