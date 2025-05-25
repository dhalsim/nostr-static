package pagegenerators

import (
	"os"
	"path/filepath"

	"nostr-static/src/helpers"
	"nostr-static/src/pagegenerators/components"
	"nostr-static/src/types"

	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
	"github.com/nbd-wtf/go-nostr"
)

type IndexArticleData struct {
	Title         string
	Summary       string
	Image         string
	Tags          []string
	Naddr         string
	AuthorName    string
	AuthorPicture string
	Nprofile      string
	CreatedAt     nostr.Timestamp
}

type IndexData struct {
	BaseFolder string
	Color      string
	Logo       string
	FaviconDir string
	Title      string
	Articles   []IndexArticleData
}

type GenerateIndexParams struct {
	BaseFolder       string
	BlogURL          string
	Events           []nostr.Event
	Profiles         map[string]nostr.Event
	OutputDir        string
	Layout           types.Layout
	EventIDToNaddr   map[string]string
	PubkeyToNProfile map[string]string
}

// Index-specific HTML rendering functions
func renderIndexHeader(title string) HTML {
	return H1_(
		// Attr(a.Class("index-title")),
		Text(title),
	)
}

func renderIndexArticle(article IndexArticleData, baseFolder string) HTML {
	return Div(Attr(a.Class("article-card")),
		components.RenderCompactProfile(
			article.AuthorName,
			article.AuthorPicture,
			article.Nprofile,
			article.Naddr,
			article.CreatedAt,
		),
		components.RenderImageHTML(article.Image, article.Title, article.Naddr, ""),
		H2_(
			A(Attr(a.Href(article.Naddr+".html")),
				Text(article.Title),
			),
		),
		components.RenderSummaryHTML(article.Summary),
		components.RenderTagsHTML(article.Tags, baseFolder),
	)
}

func renderIndexArticles(data IndexData) HTML {
	var articleElements []HTML
	for _, article := range data.Articles {
		articleElements = append(
			articleElements,
			renderIndexArticle(article, data.BaseFolder),
		)
	}

	return Div(Attr(a.Class("main-content")),
		append([]HTML{
			renderIndexHeader(data.Title),
		}, articleElements...)...,
	)
}

func GenerateIndexHTML(params GenerateIndexParams) error {
	var indexData IndexData

	indexData.Color = params.Layout.Color
	indexData.Logo = params.Layout.Logo
	indexData.FaviconDir = params.Layout.FaviconDir
	indexData.Title = params.Layout.Title

	indexData.Articles = make([]IndexArticleData, 0, len(params.Events))

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
		tags := []string{}

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
			case "t":
				tags = append(tags, tag[1])
			}
		}

		article := IndexArticleData{
			Title:         title,
			Summary:       summary,
			Image:         image,
			Tags:          tags,
			Naddr:         params.EventIDToNaddr[event.ID],
			AuthorName:    authorName,
			AuthorPicture: parsedProfile.Picture,
			Nprofile:      params.PubkeyToNProfile[event.PubKey],
			CreatedAt:     event.CreatedAt,
		}

		indexData.Articles = append(indexData.Articles, article)
	}

	// Generate feeds
	if err := GenerateFeeds(GenerateFeedParams{
		Folder:         params.OutputDir,
		FileName:       "index",
		BlogURL:        params.BlogURL,
		Events:         params.Events,
		Profiles:       params.Profiles,
		Layout:         params.Layout,
		EventIDToNaddr: params.EventIDToNaddr,
	}); err != nil {
		return err
	}

	// Generate the HTML using htmlgo
	html := Html5_(
		Head_(
			append(
				components.RenderFaviconLinks(indexData.FaviconDir),
				Meta(Attr(a.Charset("UTF-8"))),
				Meta(Attr(
					a.Name("viewport"),
					a.Content("width=device-width, initial-scale=1.0"),
				)),
				Title_(Text(indexData.Title)),
				Style_(Text_(CommonCSS+
					CommonResponsiveStyles+
					components.LogoCSS+
					components.CompactProfileCSS+
					components.FeedLinksCSS+
					components.ArticleCardCSS+
					components.TagsCSS+
					components.ImageCSS+
					components.FooterCSS+
					indexCSS)),
				components.RenderFeedLinks("index"),
				components.RenderAtomFeedLink("index"),
			)...,
		),
		Body(Attr(a.Class(indexData.Color+" index")),
			Div(Attr(a.Class("page-container")),
				Div(Attr(a.Class("logo-container")),
					components.RenderLogo(indexData.Logo, ""),
				),
				renderIndexArticles(indexData),
			),
			components.RenderFooter(),
			components.RenderFeed("index"),
			components.RenderTimeAgoScript(),
		),
	)

	// Write the HTML to file
	outputPath := filepath.Join(params.OutputDir, "index.html")
	return os.WriteFile(outputPath, []byte(html), 0644)
}

var indexCSS = `
body.index .main-content {
  flex: 1;
}

body.index .page-container {
	display: flex;
	align-items: flex-start;
}

body.index .image-container {
	max-width: 300px;
	margin: 20px auto 10px;
}
`
