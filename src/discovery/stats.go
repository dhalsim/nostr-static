package discovery

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"strings"

	"nostr-static/src/types"
)

type ZapStats struct {
	Count             int `json:"count"`
	ZapperCount       int `json:"zapper_count"`
	TargetEventCount  int `json:"target_event_count"`
	TargetPubkeyCount int `json:"target_pubkey_count"`
	ProviderCount     int `json:"provider_count"`
	Msats             int `json:"msats"`
	MinMsats          int `json:"min_msats"`
	MaxMsats          int `json:"max_msats"`
	AvgMsats          int `json:"avg_msats"`
	MedianMsats       int `json:"median_msats"`
}

type EventStats struct {
	EventID             string   `json:"event_id"`
	ReactionCount       int      `json:"reaction_count"`
	ReactionPubkeyCount int      `json:"reaction_pubkey_count"`
	RepostCount         int      `json:"repost_count"`
	RepostPubkeyCount   int      `json:"repost_pubkey_count"`
	ReplyCount          int      `json:"reply_count"`
	ReplyPubkeyCount    int      `json:"reply_pubkey_count"`
	ReportCount         int      `json:"report_count"`
	ReportPubkeyCount   int      `json:"report_pubkey_count"`
	Zaps                ZapStats `json:"zaps"`
}

type StatsResponse struct {
	EventIdToStatsMap     map[string]EventStats `json:"stats"`
	EventIdToNaddrMap     map[string]string     `json:"naddr"`
	PubkeyToAuthorDataMap map[string]AuthorData `json:"author_data"`
}

type PopularItem struct {
	EventID       string  `json:"EventID"`
	Naddr         string  `json:"Naddr"`
	Nprofile      string  `json:"Nprofile"`
	AuthorName    string  `json:"AuthorName"`
	AuthorPicture string  `json:"AuthorPicture"`
	Tag           string  `json:"Tag"`
	Title         string  `json:"Title"`
	Summary       string  `json:"Summary"`
	Image         string  `json:"Image"`
	Score         float64 `json:"Score"`
	Stats         EventStats
}

type PopularStats struct {
	Tag   string        `json:"tag"`
	Top10 []PopularItem `json:"top_10"`
}

// chunkSlice splits a slice into chunks of specified size
func chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// GetEventStats fetches stats for a batch of event IDs from nostr.band API
func GetEventStats(eventIDs []string) (map[string]EventStats, error) {
	if len(eventIDs) == 0 {
		return map[string]EventStats{}, nil
	}

	const maxBatchSize = 10
	chunks := chunkSlice(eventIDs, maxBatchSize)

	eventIdToStatsMap := make(map[string]EventStats)

	// Process each chunk
	for _, chunk := range chunks {
		// Join event IDs with commas
		idsString := strings.Join(chunk, "%2C")

		// Construct the API URL
		url := fmt.Sprintf("https://api.nostr.band/v0/stats/event/batch?objects=%s", idsString)

		// Make the HTTP request
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}

		// Check if the response status is 200
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status code: %d, url: %s", resp.StatusCode, url)
		}

		// Parse the response
		var chunkResponse StatsResponse
		if err := json.NewDecoder(resp.Body).Decode(&chunkResponse); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		resp.Body.Close()

		// Merge the chunk response into the combined response
		maps.Copy(eventIdToStatsMap, chunkResponse.EventIdToStatsMap)
	}

	return eventIdToStatsMap, nil
}

// CalculateEventScore calculates a popularity score for an event based on its stats and weights
func CalculateEventScore(stats EventStats, weights *types.ScoringWeights) float64 {
	// Base engagement score
	engagementScore :=
		float64(stats.ReactionCount)*weights.GetOrDefault(weights.ReactionWeight, types.DefaultReactionWeight) +
			float64(stats.RepostCount)*weights.GetOrDefault(weights.RepostWeight, types.DefaultRepostWeight) +
			float64(stats.ReplyCount)*weights.GetOrDefault(weights.ReplyWeight, types.DefaultReplyWeight)

	// Negative impact from reports
	reportPenalty := float64(stats.ReportCount) * weights.GetOrDefault(weights.ReportPenalty, types.DefaultReportPenalty)

	// Zap score calculation
	zapScore := 0.0
	if stats.Zaps.Msats > 0 {
		// Convert msats to sats (1000 msats = 1 sat)
		sats := float64(stats.Zaps.Msats) / 1000.0

		// Calculate zap score based on configured weights
		avgSats := float64(stats.Zaps.Msats) / float64(stats.Zaps.Count) / 1000.0
		avgSatsScore := avgSats / weights.GetOrDefault(weights.MaxSats, types.DefaultMaxSats)

		zapScore = (sats * weights.GetOrDefault(weights.ZapAmountWeight, types.DefaultZapAmountWeight)) + // Amount score
			(float64(stats.Zaps.ZapperCount) * weights.GetOrDefault(weights.ZapperCountWeight, types.DefaultZapperCountWeight)) + // Zapper count score
			(avgSatsScore * weights.GetOrDefault(weights.ZapAvgWeight, types.DefaultZapAvgWeight)) // Average zap amount score
	}

	return engagementScore - reportPenalty + (zapScore * weights.GetOrDefault(weights.ZapTotalWeight, types.DefaultZapTotalWeight))
}
