package discovery

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

type GetMostPopularEventsParams struct {
	IndexDir    string
	Tags        []string
	Limit       int
	KeepEachTag bool
}

func GetMostPopularEvents(params GetMostPopularEventsParams) ([]PopularItem, error) {
	allItems := []PopularItem{}

	for _, tag := range params.Tags {
		popular, err := os.ReadFile(filepath.Join(params.IndexDir, tag+"-popular.json"))
		if err != nil {
			return nil, err
		}

		var popularStats PopularStats
		err = json.Unmarshal(popular, &popularStats)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, popularStats.Top10...)
	}

	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].Score > allItems[j].Score
	})

	if params.KeepEachTag {
		allItems = groupByTag(allItems, params.Limit)
	} else {
		allItems = allItems[:params.Limit]
	}

	return allItems, nil
}

func groupByTag(items []PopularItem, limit int) []PopularItem {
	// First, group items by tag
	groupedItems := make(map[string][]PopularItem)
	for _, item := range items {
		groupedItems[item.Tag] = append(groupedItems[item.Tag], item)
	}

	// Ensure we have at least one item from each tag
	result := []PopularItem{}
	addedNaddrs := make(map[string]bool) // Track added naddrs to prevent duplicates

	for _, items := range groupedItems {
		if len(items) > 0 {
			item := items[0]
			if !addedNaddrs[item.Naddr] {
				result = append(result, item)
				addedNaddrs[item.Naddr] = true
			}
		}
	}

	// Add remaining items up to the limit
	remainingSlots := limit - len(result)
	if remainingSlots > 0 {
		// Sort all items by score
		sort.Slice(items, func(i, j int) bool {
			return items[i].Score > items[j].Score
		})

		// Add remaining items, skipping those we already included
		for _, item := range items {
			if len(result) >= limit {
				break
			}
			if !addedNaddrs[item.Naddr] {
				result = append(result, item)
				addedNaddrs[item.Naddr] = true
			}
		}
	}

	return result
}
