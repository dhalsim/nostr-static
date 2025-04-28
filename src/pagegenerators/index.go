package pagegenerators

import (
	"os"
	"path/filepath"

	"nostr-static/src/types"

	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
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
	CreatedAt     int64
}

type IndexData struct {
	BaseFolder string
	Color      string
	Logo       string
	Title      string
	Articles   []IndexArticleData
}

type GenerateIndexParams struct {
	BaseFolder       string
	Events           []types.Event
	Profiles         map[string]types.Event
	OutputDir        string
	Layout           types.Layout
	EventIDToNaddr   map[string]string
	PubkeyToNProfile map[string]string
}

// Index-specific HTML rendering functions
func renderIndexHeader(title string) HTML {
	return H1_(Text(title))
}

func renderIndexArticle(article IndexArticleData, baseFolder string) HTML {
	return Div(Attr(a.Class("article-card")),
		renderCompactProfile(
			article.AuthorName,
			article.AuthorPicture,
			article.Nprofile,
			article.Naddr,
			article.CreatedAt,
		),
		renderImageHTML(article.Image, article.Title, article.Naddr, ""),
		H2_(
			A(Attr(a.Href(article.Naddr+".html")),
				Text(article.Title),
			),
		),
		renderSummaryHTML(article.Summary),
		renderTagsHTML(article.Tags, baseFolder),
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
	indexData.Title = params.Layout.Title

	indexData.Articles = make([]IndexArticleData, 0, len(params.Events))

	for _, event := range params.Events {
		parsedProfile, err := parseProfile(params.Profiles[event.PubKey])
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

	// Generate the HTML using htmlgo
	html := Html5_(
		Head_(
			Meta(Attr(a.Charset("UTF-8"))),
			Meta(Attr(
				a.Name("viewport"),
				a.Content("width=device-width, initial-scale=1.0"),
			)),
			Title_(Text(indexData.Title)),
			Style_(Text_(CommonStyles+CommonResponsiveStyles)),
		),
		Body(Attr(a.Class(indexData.Color+" index")),
			Div(Attr(a.Class("page-container")),
				Div(Attr(a.Class("logo-container")),
					renderLogo(indexData.Logo, ""),
				),
				renderIndexArticles(indexData),
			),
			renderFooter(),
			renderTimeAgoScript(),
		),
	)

	// Write the HTML to file
	outputPath := filepath.Join(params.OutputDir, "index.html")
	return os.WriteFile(outputPath, []byte(html), 0644)
}
