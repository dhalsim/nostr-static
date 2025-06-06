package commands

import (
	"log"
	"path/filepath"

	"nostr-static/src/helpers"
	"nostr-static/src/pagegenerators"
	"nostr-static/src/types"

	"github.com/nbd-wtf/go-nostr"
)

type GenerateCommandParams struct {
	ConfigPath       string
	OutputDir        string
	IndexDir         string
	Config           *types.Config
	PubkeyToKind0    map[string]nostr.Event
	PubkeyToNprofile map[string]string
	Events           []nostr.Event
	EventIDToNaddr   map[string]string
}

func Generate(params GenerateCommandParams) error {
	configPath := params.ConfigPath
	outputDir := params.OutputDir
	indexDir := params.IndexDir
	config := params.Config
	events := params.Events
	eventIDToNaddr := params.EventIDToNaddr
	pubkeyToKind0 := params.PubkeyToKind0
	pubkeyToNprofile := params.PubkeyToNprofile

	// Copy logo file if specified
	if config.Layout.Logo != "" {
		logoPath := filepath.Join(filepath.Dir(configPath), config.Layout.Logo)
		outputLogoPath := filepath.Join(outputDir, config.Layout.Logo)
		if err := helpers.CopyFile(logoPath, outputLogoPath); err != nil {
			log.Printf("Warning: Failed to copy logo file: %v", err)
		}
	}

	// Copy favicon dir if specified
	if config.Layout.FaviconDir != "" {
		faviconPath := filepath.Join(filepath.Dir(configPath), config.Layout.FaviconDir)
		outputFaviconPath := filepath.Join(outputDir, config.Layout.FaviconDir)
		if err := helpers.CopyDir(faviconPath, outputFaviconPath); err != nil {
			log.Printf("Warning: Failed to copy favicon folder: %v", err)
		}
	}

	log.Println("generating article pages")

	// Article pages
	for _, event := range events {
		if err := pagegenerators.GenerateArticleHTML(pagegenerators.GenerateArticleParams{
			BaseFolder: "",
			BlogURL:    config.BlogURL,
			Event:      event,
			OutputDir:  outputDir,
			IndexDir:   indexDir,
			Settings:   config.Settings,
			Layout:     config.Layout,
			Features:   config.Features,
			Naddr:      eventIDToNaddr[event.ID],
			Profile:    pubkeyToKind0[event.PubKey],
			Nprofile:   pubkeyToNprofile[event.PubKey],
			Relays:     config.Relays,
			NostrLinks: config.Features.NostrLinks,
		}); err != nil {
			log.Printf("Failed to generate HTML for event %s: %v", event.ID, err)
			continue
		}
	}

	log.Println("generating tag pages")

	// Tag pages
	if err := pagegenerators.GenerateTagPages(pagegenerators.GenerateTagPagesParams{
		BaseFolder:       "../",
		BlogURL:          config.BlogURL,
		Events:           events,
		Profiles:         pubkeyToKind0,
		OutputDir:        outputDir,
		Layout:           config.Layout,
		EventIDToNaddr:   eventIDToNaddr,
		PubkeyToNProfile: pubkeyToNprofile,
	}); err != nil {
		return err
	}

	log.Println("generating profile pages")

	// Profile pages
	if err := pagegenerators.GenerateProfilePages(pagegenerators.GenerateProfilePagesParams{
		BaseFolder:       "../",
		NostrLinks:       config.Features.NostrLinks,
		BlogURL:          config.BlogURL,
		Profiles:         pubkeyToKind0,
		Events:           events,
		OutputDir:        outputDir,
		Layout:           config.Layout,
		PubkeyToNProfile: pubkeyToNprofile,
		EventIDToNaddr:   eventIDToNaddr,
	}); err != nil {
		return err
	}

	log.Println("generating index page")

	// Index page
	if err := pagegenerators.GenerateIndexHTML(pagegenerators.GenerateIndexParams{
		BaseFolder:       "",
		BlogURL:          config.BlogURL,
		Events:           events,
		Profiles:         pubkeyToKind0,
		OutputDir:        outputDir,
		Layout:           config.Layout,
		EventIDToNaddr:   eventIDToNaddr,
		PubkeyToNProfile: pubkeyToNprofile,
	}); err != nil {
		return err
	}

	log.Printf("Successfully generated static site in %s", outputDir)
	return nil
}
