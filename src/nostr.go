package main

import (
	"context"
	"log"
	"sync"
	"time"

	"nostr-static/src/types"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func fetchEvents(relays []string, naddrs []string) ([]types.Event, map[string]string, error) {
	var events []types.Event
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
			eventIDToNaddr[ev.ID] = naddr
			mu.Unlock()
			return true
		})
	}

	log.Printf("Finished fetching events. Total events received: %d", len(events))
	return events, eventIDToNaddr, nil
}

func fetchProfiles(
	relays []string,
	pubkeys []string,
) (map[string]types.Event, error) {
	// pubkey to kind 0 event map
	profileMap := make(map[string]types.Event)

	var mu sync.Mutex
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool := nostr.NewSimplePool(ctx)

	log.Printf("pubkeys: %v", pubkeys)

	// Create a filter for kind 0 events
	filter := nostr.Filter{
		Kinds:   []int{0},
		Authors: pubkeys,
	}

	// Fetch profile events from all relays
	eventMap := pool.FetchMany(ctx, relays, filter)

	// Process events
	for ev := range eventMap {
		mu.Lock()
		profileMap[ev.PubKey] = types.Event{
			ID:        ev.ID,
			PubKey:    ev.PubKey,
			CreatedAt: int64(ev.CreatedAt),
			Kind:      ev.Kind,
			Content:   ev.Content,
			Sig:       ev.Sig,
		}

		mu.Unlock()
	}

	return profileMap, nil
}
