package pagegenerators

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
	"github.com/nbd-wtf/go-nostr"

	"nostr-static/src/helpers"
	"nostr-static/src/pagegenerators/components"
	"nostr-static/src/types"
	"nostr-static/src/utils"
)

// ProfileData represents all data needed for profile templates
type ProfileData struct {
	BaseFolder    string
	Color         string
	Logo          string
	Nprofile      string
	PubKey        string
	Name          string
	About         string
	Picture       string
	Website       string
	DisplayName   string
	Nip05         string
	Nip05Verified bool
	Lud16         string
	Articles      []ProfileArticleData
}

type ProfileArticleData struct {
	Naddr     string
	Title     string
	Summary   string
	Image     string
	Tags      []string
	CreatedAt nostr.Timestamp
}

type GenerateProfilePagesParams struct {
	BaseFolder       string
	NostrLinks       string
	BlogURL          string
	Profiles         map[string]nostr.Event
	Events           []nostr.Event
	OutputDir        string
	Layout           types.Layout
	PubkeyToNProfile map[string]string
	EventIDToNaddr   map[string]string
}

// Profile-specific HTML rendering functions
func renderProfilePicture(data ProfileData) HTML {
	if data.Picture == "" {
		return Text("")
	}

	return Img(Attr(
		a.Src(data.Picture),
		a.Alt(data.Name),
		a.Class("profile-picture"),
	))
}

func renderProfileName(data ProfileData) HTML {
	return H2(Attr(a.Class("profile-name")),
		Text(displayNameOrName(data.DisplayName, data.Name)),
	)
}

func renderProfileAbout(data ProfileData) HTML {
	if data.About == "" {
		return Text("")
	}

	return P(Attr(a.Class("profile-about")), Text(data.About))
}

func renderProfileLinks(data ProfileData) HTML {
	return Div(Attr(a.Class("profile-links")),
		renderWebsite(data),
		renderNip05(data),
		renderLud16(data),
	)
}

func renderWebsite(data ProfileData) HTML {
	if data.Website == "" {
		return Text("")
	}

	return A(Attr(
		a.Href(data.Website),
		a.Target("_blank"),
		a.Rel("noopener noreferrer"),
	),
		Img(Attr(
			a.Src_("data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9ImN1cnJlbnRDb2xvciIgc3Ryb2tlLXdpZHRoPSIyIiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiIGNsYXNzPSJsdWNpZGUgbHVjaWRlLWxpbmsyLWljb24gbHVjaWRlLWxpbmstMiI+PHBhdGggZD0iTTkgMTdIN0E1IDUgMCAwIDEgNyA3aDIiLz48cGF0aCBkPSJNMTUgN2gyYTUgNSAwIDEgMSAwIDEwaC0yIi8+PGxpbmUgeDE9IjgiIHgyPSIxNiIgeTE9IjEyIiB5Mj0iMTIiLz48L3N2Zz4="),
			a.Alt("Website"),
			a.Class("author-website"),
			a.Width("16"),
			a.Height("16"),
		)),
		Text(data.Website),
	)
}

func renderNip05(data ProfileData) HTML {
	if data.Nip05 == "" {
		return Text("")
	}

	badgeClass := utils.Ternary(data.Nip05Verified, "verified", "unverified")

	badgeIcon := utils.Ternary(data.Nip05Verified,
		`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="green" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>`,
		`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="red" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>`,
	)

	return Span(Attr(a.Class("verified-badge "+badgeClass)),
		Text_(badgeIcon),
		Text(data.Nip05),
	)
}

func renderLud16(data ProfileData) HTML {
	if data.Lud16 == "" {
		return Text("")
	}

	return A(Attr(a.Href("lightning:"+data.Lud16)),
		Text("⚡ "+data.Lud16),
	)
}

func renderProfileArticles(data ProfileData) HTML {
	var articleElements []HTML
	for _, article := range data.Articles {
		articleElements = append(articleElements,
			Div(Attr(a.Class("article-card")),
				components.RenderImageHTML(article.Image, article.Title, article.Naddr, "../"),
				H3_(
					A(Attr(a.Href("../"+article.Naddr+".html")),
						Text(article.Title),
					),
				),
				components.RenderSummaryHTML(article.Summary),
				components.RenderTagsHTML(article.Tags, "../"),
			),
		)
	}

	return Div(Attr(a.Class("profile-articles")),
		append([]HTML{
			H2_(Text("Articles of " + displayNameOrName(data.DisplayName, data.Name))),
		}, articleElements...)...,
	)
}

func GenerateProfilePages(params GenerateProfilePagesParams) error {
	// Create profile directory if it doesn't exist
	profileDir := filepath.Join(params.OutputDir, "profile")
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return err
	}

	// Create a map to track articles by author
	authorArticles := make(map[string][]nostr.Event)
	for _, event := range params.Events {
		authorArticles[event.PubKey] = append(authorArticles[event.PubKey], event)
	}

	// Generate a page for each profile
	for pubkey, profileEvent := range params.Profiles {
		// Parse profile metadata
		parsedProfile, err := helpers.ParseProfile(profileEvent)
		if err != nil {
			return err
		}

		// Verify Nip05 if present
		nip05Verified := false
		if parsedProfile.Nip05 != "" {
			parts := strings.Split(parsedProfile.Nip05, "@")
			if len(parts) == 2 {
				username := parts[0]
				domain := parts[1]
				url := fmt.Sprintf("https://%s/.well-known/nostr.json", domain)

				log.Println("requesting nip05 url: ", url)

				resp, err := http.Get(url)
				if err == nil {
					defer resp.Body.Close()
					var data struct {
						Names map[string]string `json:"names"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
						if pubkey, ok := data.Names[username]; ok {
							nip05Verified = pubkey == profileEvent.PubKey
						}
					}
				}
			}
		}

		// Add articles for this profile
		articleEvents, ok := authorArticles[pubkey]
		if !ok {
			return fmt.Errorf("no articles found for profile: %s", pubkey)
		}

		// Generate feeds for this profile
		if err := GenerateFeeds(GenerateFeedParams{
			BlogURL:        params.BlogURL,
			Folder:         profileDir,
			FileName:       params.PubkeyToNProfile[pubkey],
			Events:         articleEvents,
			Profiles:       params.Profiles,
			Layout:         params.Layout,
			EventIDToNaddr: params.EventIDToNaddr,
		}); err != nil {
			return err
		}

		articles := make([]ProfileArticleData, len(articleEvents))

		for i, event := range articleEvents {
			metadata := helpers.ExtractArticleMetadata(event.Tags)

			articles[i] = ProfileArticleData{
				Naddr:     params.EventIDToNaddr[event.ID],
				Title:     metadata.Title,
				Summary:   metadata.Summary,
				Image:     metadata.Image,
				Tags:      metadata.Tags,
				CreatedAt: event.CreatedAt,
			}
		}

		data := ProfileData{
			BaseFolder:    params.BaseFolder,
			Color:         params.Layout.Color,
			Logo:          params.Layout.Logo,
			Nprofile:      params.PubkeyToNProfile[pubkey],
			PubKey:        pubkey,
			Name:          parsedProfile.Name,
			About:         parsedProfile.About,
			Picture:       parsedProfile.Picture,
			Website:       parsedProfile.Website,
			DisplayName:   parsedProfile.DisplayName,
			Nip05:         parsedProfile.Nip05,
			Nip05Verified: nip05Verified,
			Lud16:         parsedProfile.Lud16,
			Articles:      articles,
		}

		// Generate the HTML using htmlgo
		html := Html5_(
			Head_(
				Meta(Attr(a.Charset("UTF-8"))),
				Meta(Attr(
					a.Name("viewport"),
					a.Content("width=device-width, initial-scale=1.0"),
				)),
				Title_(Text("Profile: "+data.Name)),
				Style_(Text_(CommonCSS+
					CommonResponsiveStyles+
					components.DotMenuCSS+
					components.LogoCSS+
					components.FeedLinksCSS+
					components.ArticleCardCSS+
					components.TagsCSS+
					components.ImageCSS)),
				Style_(Text_(ProfileStyles)),
				components.RenderFeedLinks(data.Nprofile),
				components.RenderAtomFeedLink(data.Nprofile),
			),
			Body(Attr(a.Class(data.Color+" profile")),
				Div(Attr(a.Class("page-container")),
					Div(Attr(a.Class("logo-container")),
						components.RenderLogo(data.Logo, data.BaseFolder),
					),
					Div(Attr(a.Class("main-content")),
						Div(Attr(a.Class("profile-header")),
							Div(Attr(a.Class("profile-header-left")),
								renderProfilePicture(data),
								Div(Attr(a.Class("profile-header-name")),
									renderProfileName(data),
									components.RenderNostrLinks(
										"",
										data.Nprofile,
										params.NostrLinks,
									),
								),
							),
							renderProfileAbout(data),
							renderProfileLinks(data),
						),
						renderProfileArticles(data),
					),
				),
				components.RenderFooter(),
				components.RenderFeed(data.Nprofile),
				components.RenderTimeAgoScript(),
				components.RenderDropdownScript(),
			),
		)

		// Write the HTML to file
		outputFile := filepath.Join(profileDir, data.Nprofile+".html")
		if err := os.WriteFile(outputFile, []byte(html), 0644); err != nil {
			return err
		}
	}

	return nil
}

const ProfileStyles = `
body.profile .page-container {
		display: flex;
		align-items: flex-start;
		max-width: 1200px;
		margin: 0 auto;
}

body.profile .main-content {
		flex: 1;
		max-width: 800px;
}

body.profile img {
		max-width: 100%;
		height: auto;
}

/* Theme-specific author styles */
body.light .author-website {
		color: #000000;
}

body.light .author-website:hover {
		color: #0066cc;
}

body.dark .author-website {
		color: #e0e0e0;
		filter: invert(1);
}

body.dark .author-website:hover {
		color: #4a9eff;
}

.profile-name {
	display: inline-block;
}

.profile-picture {
		width: 100px;
		height: 100px;
		border-radius: 50%;
		object-fit: cover;
}

.profile-links {
		display: flex;
		align-items: center;
		gap: 10px;
}

.profile-links a {
		text-decoration: none;
}

.profile-header-name {
		display: flex;
		align-items: center;
}

.profile-header {
		display: flex;
    flex-direction: column;
    align-items: center;
}

@media (max-width: 768px) {
		.profile-header {
				flex-direction: column;
				align-items: center;
				text-align: center;
				gap: 15px;
		}

		.profile-header-left {
				display: flex;
    		flex-direction: column;
    		align-items: center;
		}

		.profile-links {
				flex-direction: column;
		}
}
`
