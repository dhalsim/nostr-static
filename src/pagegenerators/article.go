package pagegenerators

import (
	"html/template"
	"nostr-static/src/types"
	"os"
	"path/filepath"
	"strings"

	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
)

// ArticleData represents all data needed for article templates
type ArticleData struct {
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
	BaseFolder     string
}

type GenerateArticleParams struct {
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

// Regular Go functions for rendering HTML components
func renderArticleHeader(data ArticleData) HTML {
	return Div(Attr(a.Class("article-header")),
		renderCompactProfile(
			data.AuthorName,
			data.AuthorPicture,
			data.AuthorNProfile,
			data.Naddr,
			data.CreatedAt,
		),
		H1_(Text(data.Title)),
		renderSummaryHTML(data.Summary),
		renderTagsHTML(data.Tags, data.BaseFolder),
		renderImageHTML(data.Image, data.Title, data.Naddr, data.BaseFolder),
	)
}

func renderCommentsHTML(data ArticleData) HTML {
	if !data.Comments {
		return Text("")
	}

	return Text_(`<zap-threads 
		anchor="` + data.Naddr + `" 
		relays="` + data.Relays + `"
		urls="naddr:njump.me/,npub:njump.me/,nprofile:njump.me/,nevent:njump.me/,note:njump.me/,tag:snort.social/t/"
		disable="replyAnonymously" />`)
}

func renderCommentsScript(data ArticleData) HTML {
	if !data.Comments {
		return Text("")
	}

	return Text_(`
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
	`)
}

const ArticleStyles = `
body.article .page-container {
		display: flex;
		align-items: flex-start;
		max-width: 1200px;
		margin: 0 auto;
}

body.article .logo-container {
		flex: 0 0 200px;
		position: sticky;
		top: 20px;
}

body.article .main-content {
		flex: 1;
		max-width: 800px;
}

body.article img {
		max-width: 100%;
		height: auto;
}

body.article pre {
		background-color: #f5f5f5;
		padding: 15px;
		border-radius: 5px;
		overflow-x: auto;
}

body.article code {
		font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
}

body.article .article-header {
		margin-bottom: 2em;
}

body.article .article-header h1 {
		margin-bottom: 0.5em;
}

body.article article {
		word-wrap: break-word;
		overflow-wrap: break-word;
		hyphens: auto;
}

body.article article p {
		max-width: 100%;
		overflow-wrap: break-word;
		word-wrap: break-word;
		word-break: break-word;
}
		
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
		
@media (max-width: 768px) {
	.article-header {
		display: flex;
		flex-direction: column;
		align-items: baseline;
		margin-top: 15px;
	}
}
`

func GenerateArticleHTML(params GenerateArticleParams) error {
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

	profileData, err := parseProfile(profile)
	if err != nil {
		return err
	}

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
		CreatedAt:      event.CreatedAt,
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
		Comments:       features.Comments,
		Relays:         strings.Join(relays, ","),
		BaseFolder:     params.BaseFolder,
	}

	// Generate the HTML using htmlgo
	html := Html5_(
		Head_(
			Meta(Attr(a.Charset("UTF-8"))),
			Meta(Attr(
				a.Name("viewport"),
				a.Content("width=device-width, initial-scale=1.0"),
			)),
			Title_(Text(data.Title)),
			Style_(Text_(CommonStyles+
				ArticleStyles+
				CommonResponsiveStyles)),
		),
		Body(Attr(a.Class(data.Color+" article")),
			Div(Attr(a.Class("page-container")),
				Div(Attr(a.Class("logo-container")),
					renderLogo(data.Logo, data.BaseFolder),
				),
				Div(Attr(a.Class("main-content")),
					Article_(
						renderArticleHeader(data),
						Text_(string(data.Content)),
					),
					renderCommentsHTML(data),
				),
			),
			renderFooter(),
			renderCommentsScript(data),
			renderTimeAgoScript(),
		),
	)

	// Write the HTML to file
	outputPath := filepath.Join(outputDir, data.Naddr+".html")
	return os.WriteFile(outputPath, []byte(html), 0644)
}
