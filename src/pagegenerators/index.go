package pagegenerators

import (
	"html/template"
	"os"
	"path/filepath"

	"nostr-static/src/types"
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
	Ago           string
}

type IndexData struct {
	Color    string
	Logo     string
	Title    string
	Articles []IndexArticleData
}

const indexTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        ` + CommonStyles + `
				` + ResponsiveStyles + `
    </style>
</head>
<body class="{{.Color}} index">
    <div class="page-container">
        <div class="logo-container">
            {{renderLogo .Logo ""}}
        </div>
        <div class="main-content">
            <h1>{{.Title}}</h1>
            {{range .Articles}}
            <div class="article-card">
                {{renderCompactProfile .AuthorName .Nprofile .AuthorPicture .Ago ""}}
                {{renderImage .Image .Title .Naddr ""}}
                <h2><a href="{{.Naddr}}.html">{{.Title}}</a></h2>
                {{renderSummary .Summary}}
                {{renderTags .Tags ""}}
            </div>
            {{end}}
            {{renderFooter}}
        </div>
    </div>
</body>
</html>`

type GenerateIndexParams struct {
	Events           []types.Event
	Profiles         map[string]types.Event
	OutputDir        string
	Layout           types.Layout
	EventIDToNaddr   map[string]string
	PubkeyToNProfile map[string]string
}

func GenerateIndexHTML(params GenerateIndexParams) error {
	var indexData IndexData

	indexData.Color = params.Layout.Color
	indexData.Logo = params.Layout.Logo
	indexData.Title = params.Layout.Title

	indexData.Articles = make([]IndexArticleData, 0, len(params.Events))

	for _, event := range params.Events {
		parsedProfile := parseProfile(params.Profiles[event.PubKey])
		authorName := displayNameOrName(
			parsedProfile.DisplayName,
			parsedProfile.Name,
		)

		article := IndexArticleData{
			Naddr:         params.EventIDToNaddr[event.ID],
			AuthorName:    authorName,
			AuthorPicture: parsedProfile.Picture,
			Nprofile:      params.PubkeyToNProfile[event.PubKey],
			Ago:           diffString(ago(event)),
		}

		for _, tag := range event.Tags {
			if len(tag) < 2 {
				continue
			}

			switch tag[0] {
			case "title":
				article.Title = tag[1]
			case "summary":
				article.Summary = tag[1]
			case "image":
				article.Image = tag[1]
			case "t":
				article.Tags = append(article.Tags, tag[1])
			}
		}

		if article.Title == "" {
			article.Title = "Untitled Article"
		}

		indexData.Articles = append(indexData.Articles, article)
	}

	tmpl, err := template.New("index").Funcs(ComponentFuncs).Parse(indexTemplate)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(params.OutputDir, "index.html")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, indexData)
}
