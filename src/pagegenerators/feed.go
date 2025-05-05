package pagegenerators

import (
	"os"
	"path/filepath"
	"time"

	"nostr-static/src/helpers"
	"nostr-static/src/types"

	"github.com/gorilla/feeds"
	"github.com/nbd-wtf/go-nostr"
)

type GenerateFeedParams struct {
	Folder         string
	FileName       string
	BlogURL        string
	Events         []nostr.Event
	Profiles       map[string]nostr.Event
	Layout         types.Layout
	EventIDToNaddr map[string]string
}

func GenerateFeeds(params GenerateFeedParams) error {
	// Create feed
	now := time.Now()
	feed := &feeds.Feed{
		Title:       params.Layout.Title,
		Link:        &feeds.Link{Href: params.BlogURL},
		Description: "Nostr Articles Feed",
		Author:      &feeds.Author{Name: "Nostr Static"},
		Created:     now,
	}

	// Add items to feed
	for _, event := range params.Events {
		parsedProfile, err := helpers.ParseProfile(params.Profiles[event.PubKey])
		if err != nil {
			return err
		}

		authorName := displayNameOrName(
			parsedProfile.DisplayName,
			parsedProfile.Name,
		)

		title := "Untitled Article"
		summary := ""
		image := ""

		for _, tag := range event.Tags {
			if len(tag) < 2 {
				continue
			}

			switch tag[0] {
			case "title":
				title = tag[1]
			case "summary":
				summary = tag[1]
			case "image":
				image = tag[1]
			}
		}

		// Convert markdown to HTML for feed content
		htmlContent, err := convertMarkdownToHTML(event.Content, true)
		if err != nil {
			return err
		}

		// Create feed item
		articleURL := params.BlogURL + "/" + params.EventIDToNaddr[event.ID] + ".html"
		item := &feeds.Item{
			Title:       title,
			Link:        &feeds.Link{Href: articleURL},
			Description: summary,
			Author:      &feeds.Author{Name: authorName},
			Created:     time.Unix(int64(event.CreatedAt), 0),
			Content:     htmlContent,
		}

		if image != "" {
			item.Enclosure = &feeds.Enclosure{
				Url:    image,
				Type:   "image/jpeg", // Assuming JPEG, adjust if needed
				Length: "0",          // Length unknown
			}
		}

		feed.Items = append(feed.Items, item)
	}

	// Generate RSS feed
	rss, err := feed.ToRss()
	if err != nil {
		return err
	}

	// Generate Atom feed
	atom, err := feed.ToAtom()
	if err != nil {
		return err
	}

	// Write RSS feed to file
	rssPath := filepath.Join(params.Folder, params.FileName+"-rss.xml")
	if err := os.WriteFile(rssPath, []byte(rss), 0644); err != nil {
		return err
	}

	// Write Atom feed to file
	atomPath := filepath.Join(params.Folder, params.FileName+"-atom.xml")
	if err := os.WriteFile(atomPath, []byte(atom), 0644); err != nil {
		return err
	}

	return nil
}
