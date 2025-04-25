package pagegenerators

import (
	"fmt"
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
    <style>
        ` + CommonStyles + `
        .page-container {
            display: flex;
            align-items: flex-start;
            max-width: 1200px;
            margin: 0 auto;
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
<body class="{{.Color}} article">
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
                    {{renderImage .Image .Title "" ""}}
                </div>
                {{.Content}}
            </article>
        </div>
    </div>
</body>
</html>`

type GenerateArticleParams struct {
	Event          types.Event
	OutputDir      string
	Layout         types.Layout
	EventIDToNaddr map[string]string
}

func GenerateArticleHTML(params GenerateArticleParams) error {
	event := params.Event
	outputDir := params.OutputDir
	layout := params.Layout
	eventIDToNaddr := params.EventIDToNaddr

	htmlContent, err := convertMarkdownToHTML(event.Content, true)
	if err != nil {
		return err
	}

	metadata := ExtractArticleMetadata(event.Tags)
	if metadata.Title == "" {
		metadata.Title = "Untitled Article"
	}

	naddr := eventIDToNaddr[event.ID]

	if naddr == "" {
		return fmt.Errorf("naddr not found for event %s", event.ID)
	}

	data := ArticleData{
		Title:   metadata.Title,
		Content: template.HTML(htmlContent),
		Color:   layout.Color,
		Summary: metadata.Summary,
		Tags:    metadata.Tags,
		Logo:    layout.Logo,
		Image:   metadata.Image,
		Naddr:   naddr,
	}

	tmpl, err := template.New("article").Funcs(ComponentFuncs).Parse(articleTemplate)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(outputDir, naddr+".html")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}
