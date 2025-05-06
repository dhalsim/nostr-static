package pagegenerators

import (
	"html/template"
	"log"
	"nostr-static/src/discovery"
	"nostr-static/src/types"
	"os"
	"path/filepath"
	"strings"

	"nostr-static/src/helpers"
	"nostr-static/src/pagegenerators/components"

	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
	"github.com/nbd-wtf/go-nostr"
)

// ArticleData represents all data needed for article templates
type ArticleData struct {
	BlogURL        string
	Title          string
	CreatedAt      nostr.Timestamp
	Content        template.HTML
	Color          string
	Summary        string
	Tags           []string
	Logo           string
	Image          string
	EventID        string
	Naddr          string
	AuthorName     string
	AuthorNProfile string
	AuthorPicture  string
	Comments       bool
	TagDiscovery   bool
	DiscoveryData  []discovery.PopularItem
	Relays         string
	BaseFolder     string
}

type GenerateArticleParams struct {
	BaseFolder string
	BlogURL    string
	NostrLinks string
	Settings   types.Settings
	Event      nostr.Event
	OutputDir  string
	IndexDir   string
	Layout     types.Layout
	Features   types.Features
	Naddr      string
	Profile    nostr.Event
	Nprofile   string
	Relays     []string
}

// Regular Go functions for rendering HTML components
func renderArticleHeader(
	data ArticleData,
	nostrLinks string,
) HTML {
	return Div(Attr(a.Class("article-header")),
		Div(Attr(a.Class("article-header-top")),
			components.RenderCompactProfile(
				data.AuthorName,
				data.AuthorPicture,
				data.AuthorNProfile,
				data.Naddr,
				data.CreatedAt,
			),
			components.RenderNostrLinks(data.Naddr, data.AuthorNProfile, nostrLinks),
		),
		renderTitleHTML(data.Title),
		components.RenderSummaryHTML(data.Summary),
		components.RenderTagsHTML(data.Tags, data.BaseFolder),
		components.RenderImageHTML(data.Image, data.Title, data.Naddr, data.BaseFolder),
	)
}

func renderTagDiscoveryHTML(data ArticleData, nostrLinks string) HTML {
	if !data.TagDiscovery {
		return Text("")
	}

	if len(data.DiscoveryData) == 0 {
		return Text("")
	}

	// Filter out the current article
	filteredEvents := []discovery.PopularItem{}
	for _, event := range data.DiscoveryData {
		if event.EventID != data.EventID {
			filteredEvents = append(filteredEvents, event)
		}
	}

	tagDiscoveryItems := []HTML{}
	for _, popularEvent := range filteredEvents {
		tagDiscoveryItems = append(tagDiscoveryItems, renderTagDiscoveryItemHTML(popularEvent, nostrLinks))
	}

	return Div(Attr(a.Class("tag-discovery-section")),
		Div(Attr(a.Class("recommended-articles-title")), Text_("Recommended Articles")),
		Div(Attr(a.Class("tag-discovery-items")), tagDiscoveryItems...),
	)
}

func renderTagDiscoveryItemHTML(data discovery.PopularItem, nostrLinks string) HTML {
	imageContent := HTML(Img(Attr(a.Class("tag-image"), a.Src(data.Image))))
	if data.Image == "" {
		imageContent = Div(Attr(a.Class("tag-image-placeholder")), Text_("No image"))
	}

	return Div(Attr(a.Class("tag-discovery-item")),
		A(Attr(a.Class("tag-image-link"), a.Href(`https://`+nostrLinks+`/`+data.Naddr), a.Target("_blank")),
			imageContent,
		),
		Div(Attr(a.Class("tag-author-section")),
			Img(Attr(a.Class("tag-author-picture"), a.Src(data.AuthorPicture))),
			A(Attr(a.Class("tag-author"), a.Href(`https://`+nostrLinks+`/`+data.Nprofile), a.Target("_blank")), Text_(data.AuthorName)),
		),
		A(Attr(a.Class("tag-title"), a.Href(`https://`+nostrLinks+`/`+data.Naddr), a.Target("_blank")), Text_(data.Title)),
		Div(Attr(a.Class("tag-summary")), Text_(data.Summary)),
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

func attrProperty(data interface{}, templs ...string) a.Attribute {
	attr := a.Attribute{Data: data, Name: "Property"}
	if len(templs) == 0 {
		attr.Templ = `{{define "Property"}}property="{{.}}"{{end}}`
	} else {
		attr.Templ = `{{define "Property"}}property="` + strings.Join(templs, " ") + `"{{end}}`
	}
	return attr
}

func GenerateArticleHTML(params GenerateArticleParams) error {
	event := params.Event
	profile := params.Profile
	nprofile := params.Nprofile
	naddr := params.Naddr
	outputDir := params.OutputDir
	indexDir := params.IndexDir
	settings := params.Settings
	layout := params.Layout
	features := params.Features
	relays := params.Relays
	nostrLinks := params.NostrLinks
	blogURL := params.BlogURL

	htmlContent, err := convertMarkdownToHTML(event.Content, true)
	if err != nil {
		return err
	}

	metadata := helpers.ExtractArticleMetadata(event.Tags)
	if metadata.Title == "" {
		metadata.Title = "Untitled Article"
	}

	profileData, err := helpers.ParseProfile(profile)
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

	discoveryData := []discovery.PopularItem{}

	if features.TagDiscovery {
		discoveryData, err = discovery.GetMostPopularEvents(discovery.GetMostPopularEventsParams{
			IndexDir:    indexDir,
			Tags:        metadata.Tags,
			Limit:       settings.TagDiscovery.PopularArticlesCount,
			KeepEachTag: true,
		})
		if err != nil {
			if os.IsNotExist(err) {
				log.Println("Tag discovery index not found, skipping tag discovery section")
			} else {
				return err
			}
		}
	}

	data := ArticleData{
		BlogURL:        blogURL,
		Title:          metadata.Title,
		CreatedAt:      event.CreatedAt,
		Content:        template.HTML(htmlContent),
		Color:          layout.Color,
		Summary:        metadata.Summary,
		Tags:           metadata.Tags,
		Logo:           layout.Logo,
		Image:          metadata.Image,
		EventID:        event.ID,
		Naddr:          naddr,
		AuthorName:     authorName,
		AuthorNProfile: nprofile,
		AuthorPicture:  authorPicture,
		Comments:       features.Comments,
		TagDiscovery:   features.TagDiscovery,
		DiscoveryData:  discoveryData,
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
			Meta(Attr(
				attrProperty("og:title"),
				a.Content(data.Title),
			)),
			Meta(Attr(
				attrProperty("og:description"),
				a.Content(data.Summary),
			)),
			Meta(Attr(
				attrProperty("og:url"),
				a.Content(blogURL+`/`+data.Naddr),
			)),
			Meta(Attr(
				attrProperty("og:image"),
				a.Content(data.Image),
			)),
			Title_(Text(data.Title)),
			Style_(Text_(CommonCSS+
				ArticleStyles+
				CommonResponsiveStyles+
				components.DotMenuCSS+
				components.LogoCSS+
				components.CompactProfileCSS+
				components.FeedLinksCSS+
				components.TagsCSS+
				components.FooterCSS+
				components.ImageCSS)),
		),
		Body(Attr(a.Class(data.Color+" article")),
			Div(Attr(a.Class("page-container")),
				Div(Attr(a.Class("logo-container")),
					components.RenderLogo(data.Logo, data.BaseFolder),
				),
				Div(Attr(a.Class("main-content")),
					Article_(
						renderArticleHeader(data, nostrLinks),
						Text_(string(data.Content)),
					),
					renderCommentsHTML(data),
				),
			),
			renderTagDiscoveryHTML(data, nostrLinks),
			components.RenderFooter(),
			components.RenderTimeAgoScript(),
			renderCommentsScript(data),
			components.RenderDropdownScript(),
		),
	)

	// Write the HTML to file
	outputPath := filepath.Join(outputDir, data.Naddr+".html")
	return os.WriteFile(outputPath, []byte(html), 0644)
}

const ArticleStyles = `
body.article .logo-container {
	margin-top: 0;
}

body.article .page-container {
	display: flex;
	align-items: flex-start;
	max-width: 1200px;
	margin: 0 auto;
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

.article-header-top {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 1em;
}

.recommended-articles-title {
    font-size: 1.5em;
    font-weight: 600;
		margin-top: 1em;
    margin-bottom: 1em;
    color: #c99;
}

.tag-discovery-items {
    display: flex;
    overflow-x: auto;
    gap: 2em;
    scroll-behavior: smooth;
    -webkit-overflow-scrolling: touch;
    scrollbar-width: thin;
    scrollbar-color: #666 #f0f0f0;
}

.tag-discovery-items::-webkit-scrollbar {
    height: 8px;
}

.tag-discovery-items::-webkit-scrollbar-track {
    background: #f0f0f0;
    border-radius: 4px;
}

.tag-discovery-items::-webkit-scrollbar-thumb {
    background: #666;
    border-radius: 4px;
}

.tag-discovery-items::-webkit-scrollbar-thumb:hover {
    background: #888;
}

.tag-discovery-item {
    flex: 0 0 300px;
    display: flex;
    flex-direction: column;
    gap: 0.5em;
}

.tag-image-link {
    display: block;
    text-decoration: none;
    width: 100%;
    height: 200px;
    overflow: hidden;
    border-radius: 4px;
}

.tag-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.tag-image-placeholder {
    width: 100%;
    height: 100%;
    background-color: #f0f0f0;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #999;
    font-size: 0.9em;
}

.tag-author-section {
    display: flex;
    align-items: center;
    gap: 0.5em;
}

.tag-author-picture {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    object-fit: cover;
}

.tag-author {
    font-size: 0.85em;
    color: #666;
    text-decoration: none;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.tag-author:hover {
    color: #333;
}

.tag-title {
    font-size: 1.2em;
    font-weight: 600;
    color: #333;
    text-decoration: none;
    display: block;
    line-height: 1.3;
}

.tag-title:hover {
    color: #666;
}

.tag-summary {
    font-size: 0.9em;
    color: #666;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
    line-height: 1.4;
}

@media (max-width: 768px) {
    .tag-discovery-item {
        flex: 0 0 260px;
    }
    
    .tag-image-link {
        height: 180px;
    }
    
    .tag-author-picture {
        width: 20px;
        height: 20px;
    }
    
    .tag-author {
        font-size: 0.8em;
    }
}
`
