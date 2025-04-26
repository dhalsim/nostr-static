package pagegenerators

import (
	"html/template"
	"nostr-static/src/types"
	"os"
	"path/filepath"
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
				` + ResponsiveStyles + `
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
                    {{renderAuthor .AuthorName .AuthorNProfile .AuthorPicture ""}}
                    {{renderAgo .Ago}}
                    <h1>{{.Title}}</h1>
                    {{renderSummary .Summary}}
                    {{renderTags .Tags ""}}
                    {{renderImage .Image .Title "" ""}}
                </div>
                {{.Content}}
            </article>
            {{renderFooter}}
        </div>
    </div>
</body>
</html>`

type ArticleData struct {
	Title          string
	Ago            string
	Content        template.HTML
	Color          string
	Summary        string
	Tags           []string
	Logo           string
	Image          string
	Naddr          string
	AuthorName     string
	AuthorNProfile string
	AuthorPicture  string
}

type generateArticleParams struct {
	Event     types.Event
	OutputDir string
	Layout    types.Layout
	Features  types.Features
	Naddr     string
	Profile   types.Event
	Nprofile  string
}

func NewGenerateArticleParams(
	event types.Event,
	outputDir string,
	layout types.Layout,
	features types.Features,
	naddr string,
	profile types.Event,
	nprofile string,
) generateArticleParams {
	return generateArticleParams{
		Event:     event,
		OutputDir: outputDir,
		Layout:    layout,
		Features:  features,
		Naddr:     naddr,
		Profile:   profile,
		Nprofile:  nprofile,
	}
}

func GenerateArticleHTML(params generateArticleParams) error {
	event := params.Event
	profile := params.Profile
	nprofile := params.Nprofile
	naddr := params.Naddr
	outputDir := params.OutputDir
	layout := params.Layout

	htmlContent, err := convertMarkdownToHTML(event.Content, true)
	if err != nil {
		return err
	}

	metadata := ExtractArticleMetadata(event.Tags)
	if metadata.Title == "" {
		metadata.Title = "Untitled Article"
	}

	profileData := parseProfile(profile)

	authorName := "Unknown Author"
	authorPicture := ""

	if profileData.DisplayName != "" {
		authorName = profileData.DisplayName
	} else if profileData.Name != "" {
		authorName = profileData.Name
	}
	authorPicture = profileData.Picture

	data := ArticleData{
		Title:          metadata.Title,
		Ago:            diffString(ago(event)),
		Content:        template.HTML(htmlContent),
		Color:          layout.Color,
		Summary:        metadata.Summary,
		Tags:           metadata.Tags,
		Logo:           layout.Logo,
		Image:          metadata.Image,
		Naddr:          naddr,
		AuthorName:     authorName,
		AuthorNProfile: nprofile,
		AuthorPicture:  authorPicture,
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
