package pagegenerators

import (
	"html/template"
	"nostr-static/src/types"
	"os"
	"path/filepath"
	"strings"
)

// articleData represents all data needed for article templates
type articleData struct {
	BaseFolder     string
	Title          string
	CreatedAt      int64
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
	Comments       bool
	Relays         string
}

func NewArticleData(
	baseFolder string,
	title string,
	createdAt int64,
	content template.HTML,
	color string,
	summary string,
	tags []string,
	logo string,
	image string,
	naddr string,
	authorName string,
	authorNProfile string,
	authorPicture string,
	comments bool,
	relays string,
) articleData {
	return articleData{
		BaseFolder:     baseFolder,
		Title:          title,
		CreatedAt:      createdAt,
		Content:        content,
		Color:          color,
		Summary:        summary,
		Tags:           tags,
		Logo:           logo,
		Image:          image,
		Naddr:          naddr,
		AuthorName:     authorName,
		AuthorNProfile: authorNProfile,
		AuthorPicture:  authorPicture,
		Comments:       comments,
		Relays:         relays,
	}
}

// Regular Go functions for rendering HTML components
func renderArticleHeader(data articleData) template.HTML {
	return template.HTML(`
		<div class="article-header">
			` + RenderCompactProfileHTML(data) + `
			<h1>` + data.Title + `</h1>
			` + RenderSummaryHTML(data.Summary) + `
			` + RenderTagsHTML(data.Tags, "") + `
			` + RenderImageHTML(data.Image, data.Title, "", "") + `
		</div>
	`)
}

func renderCommentsHTML(data articleData) string {
	if !data.Comments {
		return ""
	}

	return `
		<zap-threads 
			anchor="` + data.Naddr + `" 
			relays="` + data.Relays + `"
			urls="naddr:njump.me/,npub:njump.me/,nprofile:njump.me/,nevent:njump.me/,note:njump.me/,tag:snort.social/t/"
			disable="replyAnonymously" />
	`
}

func renderCommentsScript(data articleData) string {
	if !data.Comments {
		return ""
	}

	return `
		<script>
			window.wnjParams = {
				position: 'bottom',
				startHidden: true,
				compactMode: true,
				disableOverflowFix: true,
			}
		</script>
		<script src="https://cdn.jsdelivr.net/npm/window.nostr.js/dist/window.nostr.min.js"></script>
		<script type="text/javascript" src="https://unpkg.com/zapthreads/dist/zapthreads.iife.js"></script>
	`
}

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

        .tags {
            margin-bottom: 1em;
        }
				` + ResponsiveStyles + `
    </style>
</head>
<body class="{{.Color}} article">
    <div class="page-container">
        <div class="logo-container">
            {{RenderLogo .}}
        </div>
        <div class="main-content">
            <article>
                {{renderArticleHeader .}}
                {{.Content}}
            </article>
            {{renderComments .}}
        </div>
    </div>
    {{RenderFooter}}
    {{renderCommentsScript .}}
    <script src="/output/static/js/time-ago.js"></script>
</body>
</html>`

type generateArticleParams struct {
	BaseFolder string
	Event      types.Event
	OutputDir  string
	Layout     types.Layout
	Features   types.Features
	Naddr      string
	Profile    types.Event
	Nprofile   string
	Relays     []string
}

func NewGenerateArticleParams(
	baseFolder string,
	event types.Event,
	outputDir string,
	layout types.Layout,
	features types.Features,
	naddr string,
	profile types.Event,
	nprofile string,
	relays []string,
) generateArticleParams {
	return generateArticleParams{
		BaseFolder: baseFolder,
		Event:      event,
		OutputDir:  outputDir,
		Layout:     layout,
		Features:   features,
		Naddr:      naddr,
		Profile:    profile,
		Nprofile:   nprofile,
		Relays:     relays,
	}
}

func GenerateArticleHTML(params generateArticleParams) error {
	baseFolder := params.BaseFolder
	event := params.Event
	profile := params.Profile
	nprofile := params.Nprofile
	naddr := params.Naddr
	outputDir := params.OutputDir
	layout := params.Layout
	features := params.Features
	relays := params.Relays

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

	data := NewArticleData(
		baseFolder,
		metadata.Title,
		event.CreatedAt,
		template.HTML(htmlContent),
		layout.Color,
		metadata.Summary,
		metadata.Tags,
		layout.Logo,
		metadata.Image,
		naddr,
		authorName,
		nprofile,
		authorPicture,
		features.Comments,
		strings.Join(relays, ","),
	)

	// Create template functions map
	funcs := template.FuncMap{
		"renderArticleHeader": func(data articleData) template.HTML {
			return renderArticleHeader(data)
		},
		"renderComments": func(data articleData) template.HTML {
			return template.HTML(renderCommentsHTML(data))
		},
		"renderCommentsScript": func(data articleData) template.HTML {
			return template.HTML(renderCommentsScript(data))
		},
	}

	tmpl, err := template.New("article").Funcs(funcs).Parse(articleTemplate)
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
