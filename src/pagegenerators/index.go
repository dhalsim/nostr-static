package pagegenerators

import (
	"html/template"
	"os"
	"path/filepath"

	"nostr-static/src/types"
)

type IndexData struct {
	Color    string
	Logo     string
	Articles []struct {
		ID        string
		Title     string
		Summary   string
		Image     string
		ImageLink string
		Tags      []string
	}
}

const indexTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Nostr Articles</title>
    <style>
        ` + CommonStyles + `
    </style>
</head>
<body class="{{.Color}} index">
    <div class="page-container">
        <div class="logo-container">
            {{renderLogo .Logo ""}}
        </div>
        <div class="main-content">
            <h1>Nostr Articles</h1>
            {{range .Articles}}
            <div class="article-card">
                {{renderImage .Image .Title .ID}}
                <h2><a href="{{.ID}}.html">{{.Title}}</a></h2>
                {{renderSummary .Summary}}
                {{renderTags .Tags ""}}
            </div>
            {{end}}
        </div>
    </div>
</body>
</html>`

func GenerateIndexHTML(
	events []types.Event,
	outputDir string,
	layout types.Layout,
) error {
	var indexData IndexData

	indexData.Color = layout.Color
	indexData.Logo = layout.Logo
	indexData.Articles = make([]struct {
		ID        string
		Title     string
		Summary   string
		Image     string
		ImageLink string
		Tags      []string
	}, 0, len(events))

	for _, event := range events {
		article := struct {
			ID        string
			Title     string
			Summary   string
			Image     string
			ImageLink string
			Tags      []string
		}{
			ID: event.ID,
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

	outputPath := filepath.Join(outputDir, "index.html")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Generate tag pages
	if err := GenerateTagPages(events, outputDir, layout); err != nil {
		return err
	}

	return tmpl.Execute(file, indexData)
}
