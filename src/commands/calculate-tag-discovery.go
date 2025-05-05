package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"nostr-static/src/discovery"
	"nostr-static/src/helpers"
	"nostr-static/src/types"

	"github.com/nbd-wtf/go-nostr"
)

type CalculateTagDiscoveryCommandParams struct {
	IndexDir  string
	OutputDir string
	Config    *types.Config
}

func CalculateTagDiscovery(params CalculateTagDiscoveryCommandParams) error {
	indexDir := params.IndexDir
	outputDir := params.OutputDir
	config := params.Config
	weights := config.Settings.TagDiscovery.Weights

	log.Println("Calculating most popular articles for each tag...")

	// read tags from file
	tags, err := os.ReadFile(filepath.Join(outputDir, "tags.txt"))
	if err != nil {
		return err
	}

	tagsArray := strings.Split(string(tags), "\n")

	for _, tag := range tagsArray {
		if tag == "" {
			continue
		}

		log.Printf("Processing tag: %s", tag)

		// Read stats file
		statsData, err := os.ReadFile(filepath.Join(indexDir, tag+"-stats.json"))
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("No stats file found for tag %s, skipping...", tag)
				continue
			}
			return err
		}

		var statsResponse discovery.StatsResponse
		if err := json.Unmarshal(statsData, &statsResponse); err != nil {
			return fmt.Errorf("failed to unmarshal stats for tag %s: %w", tag, err)
		}

		// Read tag event file
		tagEventsData, err := os.ReadFile(filepath.Join(indexDir, tag+".json"))
		if err != nil {
			return err
		}

		var naddrToEvent map[string]nostr.Event
		err = json.Unmarshal(tagEventsData, &naddrToEvent)
		if err != nil {
			return err
		}

		// Calculate scores for each event
		var scoredEvents []discovery.PopularItem
		for eventID, stats := range statsResponse.EventIdToStatsMap {
			naddr := statsResponse.EventIdToNaddrMap[eventID]
			event := naddrToEvent[naddr]

			metadata := helpers.ExtractArticleMetadata(event.Tags)

			score := discovery.CalculateEventScore(stats, &weights)
			authorData := statsResponse.PubkeyToAuthorDataMap[event.PubKey]

			scoredEvents = append(scoredEvents, discovery.PopularItem{
				EventID:       eventID,
				Naddr:         statsResponse.EventIdToNaddrMap[eventID],
				Nprofile:      authorData.Nprofile,
				AuthorName:    authorData.Name,
				AuthorPicture: authorData.Picture,
				Tag:           tag,
				Title:         metadata.Title,
				Summary:       metadata.Summary,
				Image:         metadata.Image,
				Score:         score,
				Stats:         stats,
			})
		}

		// Sort events by score in descending order
		sort.Slice(scoredEvents, func(i, j int) bool {
			return scoredEvents[i].Score > scoredEvents[j].Score
		})

		// Take top 10 events
		topEvents := scoredEvents
		if len(scoredEvents) > 10 {
			topEvents = scoredEvents[:10]
		}

		// Create output structure
		output := struct {
			Tag   string                  `json:"tag"`
			Top10 []discovery.PopularItem `json:"top_10"`
		}{
			Tag:   tag,
			Top10: topEvents,
		}

		// Write results to file
		outputData, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal output for tag %s: %w", tag, err)
		}

		outputPath := filepath.Join(params.IndexDir, tag+"-popular.json")
		if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
			return fmt.Errorf("failed to write output for tag %s: %w", tag, err)
		}

		log.Printf("Wrote popular articles for tag %s to %s", tag, outputPath)
	}

	return nil
}
