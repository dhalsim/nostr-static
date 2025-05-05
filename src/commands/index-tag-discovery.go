package commands

import (
	"context"
	"encoding/json"
	"log"
	"nostr-static/src/helpers"
	"nostr-static/src/types"
	"os"
	"path/filepath"
	"strings"
	"time"

	"nostr-static/src/discovery"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

type IndexTagDiscoveryCommandParams struct {
	Reset     bool
	IndexDir  string
	OutputDir string
	Config    *types.Config
}

func IndexTagDiscovery(params IndexTagDiscoveryCommandParams) error {
	log.Println("Indexing tags...")

	if params.Reset {
		log.Println("Reset flag is enabled")
		if err := os.RemoveAll(params.IndexDir); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(params.IndexDir, 0755); err != nil {
		return err
	}

	tags, err := os.ReadFile(filepath.Join(params.OutputDir, "tags.txt"))
	if err != nil {
		return err
	}

	tagsArray := strings.Split(string(tags), "\n")
	relays := params.Config.Relays

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool := nostr.NewSimplePool(ctx)

	for _, tag := range tagsArray {
		log.Println("Indexing tag: ", tag)

		// Fetch new events
		naddrToTagEventMap, pubkeyToAuthorDataMap, err := fetchTagEvents(
			ctx,
			pool,
			relays,
			tag,
			params.Config.Settings.TagDiscovery.FetchCountPerTag,
		)
		if err != nil {
			return err
		}

		// Save events and stats
		if err := saveTagEvents(
			params.IndexDir,
			tag,
			naddrToTagEventMap,
			pubkeyToAuthorDataMap,
		); err != nil {
			return err
		}
	}

	return nil
}

func fetchTagEvents(
	ctx context.Context,
	pool *nostr.SimplePool,
	relays []string,
	tag string,
	fetchCountPerTag int,
) (
	map[string]nostr.Event,
	map[string]discovery.AuthorData,
	error,
) {
	// Create a map to store the latest replaceable events
	latestEvents := make(map[nostr.ReplaceableKey]*nostr.Event)

	// Create the check duplicate function
	checkDuplicate := nostr.WithCheckDuplicateReplaceable(func(rk nostr.ReplaceableKey, ts nostr.Timestamp) bool {
		// If we already have this event and it's newer, keep the old one
		if existing, exists := latestEvents[rk]; exists && existing.CreatedAt >= ts {
			return true // skip this event
		}
		return false // process this event
	})

	tagFilter := nostr.Filter{
		Kinds: []int{nostr.KindArticle},
		Tags:  nostr.TagMap{"t": []string{tag}},
		Limit: fetchCountPerTag,
	}

	// Fetch tagEvents
	tagEvents := pool.FetchMany(ctx, relays, tagFilter, checkDuplicate)

	// Create a map to store which relays returned each event
	pubkeys := []string{}
	relayMap := make(map[nostr.ReplaceableKey][]string)
	for ev := range tagEvents {
		rk := nostr.ReplaceableKey{
			PubKey: ev.PubKey,
			D:      ev.Tags.GetD(),
		}

		// Store the event and append its relay
		latestEvents[rk] = ev.Event
		relayMap[rk] = append(relayMap[rk], ev.Relay.URL)

		pubkeys = append(pubkeys, ev.PubKey)
	}

	pubkeyToAuthorDataMap, err := GetAuthorData(pubkeys, relays, pool, ctx)
	if err != nil {
		return nil, nil, err
	}

	// Convert the latest events to the format we need
	naddrToTagEventMap := make(map[string]nostr.Event)
	for rk, event := range latestEvents {
		naddr := nip19.EncodePointer(nostr.EntityPointer{
			PublicKey:  rk.PubKey,
			Relays:     relayMap[rk],
			Kind:       event.Kind,
			Identifier: rk.D,
		})

		naddrToTagEventMap[naddr] = *event
	}

	return naddrToTagEventMap, pubkeyToAuthorDataMap, nil
}

// saveTagEvents saves tag events and their stats to files
func saveTagEvents(
	indexDir,
	tag string,
	naddrToEventMap map[string]nostr.Event,
	pubkeyToAuthorDataMap map[string]discovery.AuthorData,
) error {
	// Save events
	jsonData, err := json.Marshal(naddrToEventMap)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(indexDir, tag+".json"), jsonData, 0644); err != nil {
		return err
	}

	// Get event IDs for stats
	eventIds := make([]string, 0, len(naddrToEventMap))
	for _, event := range naddrToEventMap {
		eventIds = append(eventIds, event.ID)
	}

	eventIdToNaddrMap := make(map[string]string)
	for naddr, event := range naddrToEventMap {
		eventIdToNaddrMap[event.ID] = naddr
	}

	// Get and save stats
	eventIdToStatsMap, err := discovery.GetEventStats(eventIds)
	if err != nil {
		return err
	}

	statsResponse := discovery.StatsResponse{
		EventIdToStatsMap:     eventIdToStatsMap,
		EventIdToNaddrMap:     eventIdToNaddrMap,
		PubkeyToAuthorDataMap: pubkeyToAuthorDataMap,
	}

	statsData, err := json.Marshal(statsResponse)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(indexDir, tag+"-stats.json"), statsData, 0644)
}

func GetAuthorData(
	pubkeys []string,
	relays []string,
	pool *nostr.SimplePool,
	ctx context.Context,
) (map[string]discovery.AuthorData, error) {
	// Fetch profile events
	profileFilter := nostr.Filter{
		Kinds:   []int{nostr.KindProfileMetadata},
		Authors: pubkeys,
	}

	profileEvents := pool.FetchMany(ctx, relays, profileFilter)

	pubkeyToAuthorDataMap := make(map[string]discovery.AuthorData)
	profileRelayMap := make(map[string][]string)
	for ev := range profileEvents {
		profileRelayMap[ev.PubKey] = append(profileRelayMap[ev.PubKey], ev.Relay.URL)

		parsedProfile, err := helpers.ParseProfile(*ev.Event)
		if err != nil {
			return nil, err
		}

		// Create nprofile after we have relay information
		nprofile, err := nip19.EncodeProfile(ev.PubKey, profileRelayMap[ev.PubKey])
		if err != nil {
			return nil, err
		}

		pubkeyToAuthorDataMap[ev.PubKey] = discovery.AuthorData{
			Nprofile: nprofile,
			Name:     helpers.NameOrDisplayName(parsedProfile),
			Picture:  helpers.PictureOrFallback(parsedProfile, ev.PubKey),
		}
	}

	return pubkeyToAuthorDataMap, nil
}
