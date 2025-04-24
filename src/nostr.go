package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"nostr-static/src/types"

	"github.com/nbd-wtf/go-nostr"
)

func fetchEvents(relays []string, articleIDs []string) ([]types.Event, error) {
	var events []types.Event
	var mu sync.Mutex

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Starting to fetch events from %d relays for %d article IDs", len(relays), len(articleIDs))
	log.Printf("Article IDs: %v", articleIDs)

	pool := nostr.NewSimplePool(ctx)

	// Create filter for the article IDs
	filter := nostr.Filter{
		IDs: articleIDs,
	}

	// Fetch events from all relays
	eventChan := pool.FetchMany(ctx, relays, filter)

	// Process events
	for ev := range eventChan {
		log.Printf("Received event with ID: %s", ev.ID)
		tags := make([][]string, len(ev.Tags))
		for i, tag := range ev.Tags {
			tags[i] = []string(tag)
		}

		mu.Lock()
		events = append(events, types.Event{
			ID:        ev.ID,
			PubKey:    ev.PubKey,
			CreatedAt: int64(ev.CreatedAt),
			Kind:      ev.Kind,
			Tags:      tags,
			Content:   ev.Content,
			Sig:       ev.Sig,
		})
		mu.Unlock()
	}

	log.Printf("Finished fetching events. Total events received: %d", len(events))
	return events, nil
}

func saveEvents(events []types.Event, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for _, event := range events {
		if err := encoder.Encode(event); err != nil {
			return err
		}
	}

	return nil
}
