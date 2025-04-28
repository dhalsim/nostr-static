package helpers

import "github.com/nbd-wtf/go-nostr"

// ConvertTags converts nostr.Tags to [][]string
func ConvertTags(tags nostr.Tags) [][]string {
	result := make([][]string, len(tags))
	for i, tag := range tags {
		result[i] = []string(tag)
	}
	return result
}
