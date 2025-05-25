package pagegenerators

import "github.com/nbd-wtf/go-nostr"

// TagData represents all data needed for tag templates
type TagData struct {
	BaseFolder string
	Tag        string
	Color      string
	Logo       string
	FaviconDir string
	Articles   []TagArticleData
}

type TagArticleData struct {
	Naddr         string
	Title         string
	Summary       string
	Image         string
	Tags          []string
	AuthorName    string
	Nprofile      string
	AuthorPicture string
	CreatedAt     nostr.Timestamp
}
