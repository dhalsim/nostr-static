package helpers

import (
	"context"
	"encoding/json"
	"log"
	"sort"
	"sync"
	"time"

	"nostr-static/src/types"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func FetchEvents(
	relays []string,
	naddrs []string,
) ([]nostr.Event, map[string]string, error) {
	var events []nostr.Event
	var mu sync.Mutex
	eventIDToNaddr := make(map[string]string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool := nostr.NewSimplePool(ctx)

	// Process each naddr
	for _, naddr := range naddrs {
		log.Printf("Processing naddr: %s", naddr)

		// Decode naddr to get pubkey and d tag
		prefix, data, err := nip19.Decode(naddr)
		if err != nil {
			log.Printf("Error decoding naddr %s: %v", naddr, err)
			continue
		}

		if prefix != "naddr" {
			log.Printf("Invalid naddr prefix for %s: %s", naddr, prefix)
			continue
		}

		addr := data.(nostr.EntityPointer)

		if addr.Kind != nostr.KindArticle {
			log.Printf("Invalid kind for %s: %d", naddr, addr.Kind)
			continue
		}

		// Create filter for replaceable event
		filter := nostr.Filter{
			Kinds:   []int{nostr.KindArticle},
			Authors: []string{addr.PublicKey},
			Tags: nostr.TagMap{
				"d": []string{addr.Identifier},
			},
		}

		allRelays := append(relays, addr.Relays...)

		// Fetch replaceable event from all relays
		eventMap := pool.FetchManyReplaceable(ctx, allRelays, filter)

		// Process events
		eventMap.Range(func(key nostr.ReplaceableKey, ev *nostr.Event) bool {
			log.Printf("Received event with ID: %s", ev.ID)

			tags := make(nostr.Tags, len(ev.Tags))
			for i, tag := range ev.Tags {
				tags[i] = nostr.Tag(tag)
			}

			mu.Lock()

			events = append(events, nostr.Event{
				ID:        ev.ID,
				PubKey:    ev.PubKey,
				CreatedAt: ev.CreatedAt,
				Kind:      ev.Kind,
				Tags:      tags,
				Content:   ev.Content,
				Sig:       ev.Sig,
			})
			eventIDToNaddr[ev.ID] = naddr

			mu.Unlock()

			return true
		})
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt > events[j].CreatedAt
	})

	log.Printf("Finished fetching events. Total events received: %d", len(events))

	return events, eventIDToNaddr, nil
}

func FetchProfiles(
	relays []string,
	pubkeys []string,
) (map[string]nostr.Event, error) {
	// pubkey to kind 0 event map
	pubkeyToKind0 := make(map[string]nostr.Event)

	var mu sync.Mutex
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool := nostr.NewSimplePool(ctx)

	// Create a filter for kind 0 events
	filter := nostr.Filter{
		Kinds:   []int{0},
		Authors: pubkeys,
	}

	log.Printf("Filter: %v", filter)

	// Fetch profile events from all relays
	eventMap := pool.FetchMany(ctx, relays, filter)

	// Process events
	for ev := range eventMap {
		mu.Lock()
		pubkeyToKind0[ev.PubKey] = nostr.Event{
			ID:        ev.ID,
			PubKey:    ev.PubKey,
			CreatedAt: nostr.Timestamp(ev.CreatedAt),
			Kind:      ev.Kind,
			Content:   ev.Content,
			Sig:       ev.Sig,
			Tags:      ev.Tags,
		}

		mu.Unlock()
	}

	return pubkeyToKind0, nil
}

// ExtractArticleMetadata extracts metadata from event tags
func ExtractArticleMetadata(tags nostr.Tags) types.ArticleMetadata {
	var metadata types.ArticleMetadata
	for _, tag := range tags {
		if len(tag) < 2 {
			continue
		}

		switch tag[0] {
		case "title":
			metadata.Title = tag[1]
		case "summary":
			metadata.Summary = tag[1]
		case "image":
			metadata.Image = tag[1]
		case "t":
			metadata.Tags = append(metadata.Tags, tag[1])
		}
	}
	return metadata
}

func ParseProfile(event nostr.Event) (*types.ParsedProfile, error) {
	var profile types.ParsedProfile

	if err := json.Unmarshal([]byte(event.Content), &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
