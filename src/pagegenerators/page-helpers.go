package pagegenerators

import (
	"fmt"
	"time"

	"nostr-static/src/types"
)

func ago(event types.Event) time.Duration {
	now := time.Now()
	diff := now.Sub(time.Unix(event.CreatedAt, 0))

	return diff
}

// convert a duration to a string like "1 day ago"
func diffString(diff time.Duration) string {
	years := int(diff.Hours() / 8760)
	days := int(diff.Hours() / 24)
	hours := int(diff.Hours())
	minutes := int(diff.Minutes())
	seconds := int(diff.Seconds())

	if years > 0 {
		return fmt.Sprintf("%d years ago", years)
	}

	if days > 0 {
		return fmt.Sprintf("%d days ago", days)
	}

	if hours > 0 {
		return fmt.Sprintf("%d hours ago", hours)
	}

	if minutes > 0 {
		return fmt.Sprintf("%d minutes ago", minutes)
	}

	return fmt.Sprintf("%d seconds ago", seconds)
}

func ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

func apply[T any, R any](input []T, mapper func(T) R) []R {
	output := make([]R, len(input))

	for i, v := range input {
		output[i] = mapper(v)
	}

	return output
}
