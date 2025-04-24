package pagegenerators

import (
	"html/template"
	"os"
	"path/filepath"

	"nostr-static/src/types"
)

const articleTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    {{.ThemeCSS}}
    <style>
        ` + CommonStyles + `
        .page-container {
            display: flex;
            gap: 2em;
            align-items: flex-start;
            max-width: 1200px;
            margin: 0 auto;
        }
        .logo-container {
            flex: 0 0 200px;
            position: sticky;
            top: 20px;
        }
        .main-content {
            flex: 1;
            max-width: 800px;
        }
        img {
            max-width: 100%;
            height: auto;
        }
        pre {
            background-color: #f5f5f5;
            padding: 15px;
            border-radius: 5px;
            overflow-x: auto;
        }
        code {
            font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
        }
        .article-header {
            margin-bottom: 2em;
        }
        .article-header h1 {
            margin-bottom: 0.5em;
        }
        article {
            word-wrap: break-word;
            overflow-wrap: break-word;
            hyphens: auto;
        }
        article p {
            max-width: 100%;
            overflow-wrap: break-word;
            word-wrap: break-word;
            word-break: break-word;
        }
    </style>
</head>
<body>
    <div class="page-container">
        <div class="logo-container">
            {{renderLogo .Logo ""}}
        </div>
        <div class="main-content">
            <article>
                <div class="article-header">
                    <h1>{{.Title}}</h1>
                    {{renderSummary .Summary}}
                    {{renderTags .Tags ""}}
                </div>
                {{.Content}}
            </article>
        </div>
    </div>
</body>
</html>`

func GenerateArticleHTML(event types.Event, outputDir string, themeColor string, logo string) error {
	htmlContent, err := convertMarkdownToHTML(event.Content)
	if err != nil {
		return err
	}

	// Extract title, summary, and tags from tags
	var title, summary string
	var tags []string
	for _, tag := range event.Tags {
		if len(tag) < 2 {
			continue
		}

		switch tag[0] {
		case "title":
			title = tag[1]
		case "summary":
			summary = tag[1]
		case "t":
			tags = append(tags, tag[1])
		}
	}

	if title == "" {
		title = "Untitled Article"
	}

	data := ArticleData{
		Title:    title,
		Content:  template.HTML(htmlContent),
		ThemeCSS: GetThemeCSS(themeColor),
		Summary:  summary,
		Tags:     tags,
		Logo:     logo,
	}

	tmpl, err := template.New("article").Funcs(ComponentFuncs).Parse(articleTemplate)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(outputDir, event.ID+".html")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}
