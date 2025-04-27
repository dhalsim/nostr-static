package pagegenerators

import (
	"html/template"
	"nostr-static/src/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
)

// TemplateData represents all data needed for article templates
type TemplateData struct {
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

// Regular Go functions for rendering HTML components
func renderArticleHeader(data TemplateData) HTML {
	return Div(Attr(a.Class("article-header")),
		renderCompactProfileHTML(data),
		H1_(Text(data.Title)),
		renderSummaryHTML(data.Summary),
		renderTagsHTML(data.Tags, data.BaseFolder),
		renderImageHTML(data.Image, data.Title, data.Naddr, data.BaseFolder),
	)
}

func renderCompactProfileHTML(data TemplateData) HTML {
	if data.AuthorName == "" {
		return Text("")
	}

	var pictureHTML HTML
	if data.AuthorPicture != "" {
		pictureHTML = Img(Attr(
			a.Src(data.AuthorPicture),
			a.Alt(data.AuthorName),
			a.Class("compact-profile-picture"),
		))
	}

	return Div(Attr(a.Class("compact-profile")),
		A(Attr(
			a.Href("profile/"+data.AuthorNProfile+".html"),
			a.Class("compact-profile-link"),
		),
			pictureHTML,
			Span(Attr(a.Class("compact-profile-name")),
				Text(data.AuthorName),
			),
		),
		A(Attr(
			a.Href(data.Naddr+".html"),
			a.Class("compact-profile-ago"),
		),
			Span(Attr(
				a.Class("time-ago"),
				a.Dataset("timestamp", strconv.FormatInt(data.CreatedAt, 10)),
			)),
		),
	)
}

func renderSummaryHTML(summary string) HTML {
	if summary == "" {
		return Text("")
	}
	return P(Attr(a.Class("summary")), Text(summary))
}

func renderTagsHTML(tags []string, baseFolder string) HTML {
	if len(tags) == 0 {
		return Text("")
	}

	var tagElements []HTML
	for _, tag := range tags {
		baseFolder = strings.Trim(baseFolder, "/")
		if baseFolder != "" {
			baseFolder = baseFolder + "/"
		}
		tagElements = append(tagElements,
			Span(Attr(a.Class("tag")),
				A(Attr(
					a.Href(baseFolder+"tag/"+strings.ToLower(tag)+".html"),
				),
					Text(tag),
				),
			),
		)
	}

	return Div(Attr(a.Class("tags")), tagElements...)
}

func renderImageHTML(image, alt, imageLink, baseFolder string) HTML {
	if image == "" {
		return Text("")
	}

	if imageLink == "" {
		return Div(Attr(a.Class("image-container")),
			Img(Attr(
				a.Src(image),
				a.Alt(alt),
			)),
		)
	}

	return Div(Attr(a.Class("image-container")),
		A(Attr(a.Href(baseFolder+imageLink+".html")),
			Img(Attr(
				a.Src(image),
				a.Alt(alt),
			)),
		),
	)
}

func renderCommentsHTML(data TemplateData) HTML {
	if !data.Comments {
		return Text("")
	}

	return Text(`<zap-threads 
		anchor="` + data.Naddr + `" 
		relays="` + data.Relays + `"
		urls="naddr:njump.me/,npub:njump.me/,nprofile:njump.me/,nevent:njump.me/,note:njump.me/,tag:snort.social/t/"
		disable="replyAnonymously" />`)
}

func renderCommentsScript(data TemplateData) HTML {
	if !data.Comments {
		return Text("")
	}

	return Text(`
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

func renderLogo(data TemplateData) HTML {
	if data.Logo == "" {
		return Text("")
	}

	baseFolder := strings.Trim(data.BaseFolder, "/")
	if baseFolder != "" {
		baseFolder = baseFolder + "/"
	}

	return Div(Attr(a.Class("logo")),
		A(Attr(a.Href(baseFolder+"index.html")),
			Img(Attr(
				a.Src(baseFolder+data.Logo),
				a.Alt("Logo"),
			)),
		),
	)
}

func renderFooter() HTML {
	return Div(Attr(a.Class("footer")),
		Text("Built with "),
		A(Attr(
			a.Href("https://github.com/dhalsim/nostr-static"),
			a.Target("_blank"),
		),
			Text("nostr-static"),
		),
	)
}

func renderTimeAgoScript() HTML {
	return Script(Attr(a.Src("/output/static/js/time-ago.js")), JavaScript(""))
}

func GenerateArticleHTML(params generateArticleParams) error {
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

	data := TemplateData{
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
			Style_(Text(CommonStyles+ResponsiveStyles)),
		),
		Body(Attr(a.Class(data.Color+" article")),
			Div(Attr(a.Class("page-container")),
				Div(Attr(a.Class("logo-container")),
					renderLogo(data),
				),
				Div(Attr(a.Class("main-content")),
					Article_(
						renderArticleHeader(data),
						Text(string(data.Content)),
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
