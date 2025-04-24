package pagegenerators

import (
	"html/template"
	"os"
	"path/filepath"

	"nostr-static/src/types"
)

type IndexData struct {
	ThemeCSS template.HTML
	Logo     string
	Articles []struct {
		ID      string
		Title   string
		Summary string
		Image   string
		Tags    []string
	}
}

const indexTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Nostr Articles</title>
    {{.ThemeCSS}}
    <style>
        ` + CommonStyles + `
        .page-container {
            display: flex;
            gap: 2em;
            align-items: flex-start;
        }
        .logo-container {
            flex: 0 0 200px;
            position: sticky;
            top: 20px;
        }
        .main-content {
            flex: 1;
        }
    </style>
</head>
<body>
    <div class="page-container">
        <div class="logo-container">
            {{renderLogo .Logo ""}}
        </div>
        <div class="main-content">
            <h1>Nostr Articles</h1>
            {{range .Articles}}
            <div class="article-card">
                {{renderImage .Image .Title}}
                <h2><a href="{{.ID}}.html">{{.Title}}</a></h2>
                {{renderSummary .Summary}}
                {{renderTags .Tags ""}}
            </div>
            {{end}}
        </div>
    </div>
</body>
</html>`

func GenerateIndexHTML(events []types.Event, outputDir string, themeColor string, logo string) error {
	var indexData IndexData
	indexData.ThemeCSS = GetThemeCSS(themeColor)
	indexData.Logo = logo
	indexData.Articles = make([]struct {
		ID      string
		Title   string
		Summary string
		Image   string
		Tags    []string
	}, 0, len(events))

	for _, event := range events {
		article := struct {
			ID      string
			Title   string
			Summary string
			Image   string
			Tags    []string
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
	if err := GenerateTagPages(events, outputDir, themeColor, logo); err != nil {
		return err
	}

	return tmpl.Execute(file, indexData)
}
