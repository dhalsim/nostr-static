package pagegenerators

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"nostr-static/src/types"

	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
)

type GenerateTagPagesParams struct {
	BaseFolder       string
	Events           []types.Event
	Profiles         map[string]types.Event
	OutputDir        string
	Layout           types.Layout
	EventIDToNaddr   map[string]string
	PubkeyToNProfile map[string]string
}

// Tag-specific HTML rendering functions
func renderTagHeader(tag string) HTML {
	return H1_(Text("Articles tagged with \"" + tag + "\""))
}

func renderTagArticle(article TagArticleData, baseFolder string) HTML {
	return Div(Attr(a.Class("article-card")),
		renderTagCompactProfile(article),
		renderImageHTML(article.Image, article.Title, article.Naddr, baseFolder),
		H2_(
			A(Attr(a.Href(article.Naddr+".html")),
				Text(article.Title),
			),
		),
		renderSummaryHTML(article.Summary),
		renderTagsHTML(article.Tags, baseFolder),
	)
}

func renderTagArticles(data TagData) HTML {
	var articleElements []HTML
	for _, article := range data.Articles {
		articleElements = append(
			articleElements,
			renderTagArticle(article, data.BaseFolder),
		)
	}

	return Div(Attr(a.Class("main-content")),
		append([]HTML{
			renderTagHeader(data.Tag),
		}, articleElements...)...,
	)
}

func renderTagCompactProfile(data TagArticleData) HTML {
	if data.AuthorName == "" {
		return Text("")
	}

	var pictureHTML HTML
	if data.AuthorPicture != "" {
		pictureHTML = Img(Attr(
			a.Src(data.AuthorPicture),
			a.Alt(data.AuthorName),
			a.Class("compact-profile-picture"),
		))
	}

	return Div(Attr(a.Class("compact-profile")),
		A(Attr(
			a.Href("profile/"+data.Nprofile+".html"),
			a.Class("compact-profile-link"),
		),
			pictureHTML,
			Span(Attr(a.Class("compact-profile-name")),
				Text(data.AuthorName),
			),
		),
		A(Attr(
			a.Href(data.Naddr+".html"),
			a.Class("compact-profile-ago"),
		),
			Span(Attr(
				a.Class("time-ago"),
				a.Dataset("timestamp", strconv.FormatInt(data.CreatedAt, 10)),
			)),
		),
	)
}

func GenerateTagPages(params GenerateTagPagesParams) error {
	// Create a map to track tags and their associated events
	tagMap := make(map[string][]types.Event)

	// Populate the tag map
	for _, event := range params.Events {
		// Extract tags from the event
		for _, tagArray := range event.Tags {
			if len(tagArray) > 1 && tagArray[0] == "t" {
				tag := strings.ToLower(tagArray[1])
				tagMap[tag] = append(tagMap[tag], event)
			}
		}
	}

	// Create tag directory if it doesn't exist
	tagDir := filepath.Join(params.OutputDir, "tag")
	if err := os.MkdirAll(tagDir, 0755); err != nil {
		return err
	}

	// Generate a page for each tag
	for tag, tagEvents := range tagMap {
		data := TagData{
			BaseFolder: params.BaseFolder,
			Tag:        tag,
			Color:      params.Layout.Color,
			Logo:       params.Layout.Logo,
			Articles:   make([]TagArticleData, len(tagEvents)),
		}

		for i, event := range tagEvents {
			// Extract metadata from tags
			metadata := ExtractArticleMetadata(event.Tags)

			parsedProfile, err := parseProfile(params.Profiles[event.PubKey])
			if err != nil {
				return err
			}

			data.Articles[i] = TagArticleData{
				Naddr:         params.EventIDToNaddr[event.ID],
				Title:         metadata.Title,
				Summary:       metadata.Summary,
				Image:         metadata.Image,
				Tags:          metadata.Tags,
				AuthorName:    displayNameOrName(parsedProfile.DisplayName, parsedProfile.Name),
				AuthorPicture: parsedProfile.Picture,
				CreatedAt:     event.CreatedAt,
				Nprofile:      params.PubkeyToNProfile[event.PubKey],
			}
		}

		// Generate the HTML using htmlgo
		html := Html5_(
			Head_(
				Meta(Attr(a.Charset("UTF-8"))),
				Meta(Attr(
					a.Name("viewport"),
					a.Content("width=device-width, initial-scale=1.0"),
				)),
				Title_(Text("Tag: "+data.Tag)),
				Style_(Text_(CommonStyles+CommonResponsiveStyles)),
			),
			Body(Attr(a.Class(data.Color+" tags")),
				Div(Attr(a.Class("page-container")),
					Div(Attr(a.Class("logo-container")),
						renderLogo(data.Logo, "../"),
					),
					renderTagArticles(data),
				),
				renderFooter(),
				renderTimeAgoScript(),
			),
		)

		// Write the HTML to file
		outputFile := filepath.Join(tagDir, strings.ToLower(tag)+".html")
		if err := os.WriteFile(outputFile, []byte(html), 0644); err != nil {
			return err
		}
	}

	return nil
}
