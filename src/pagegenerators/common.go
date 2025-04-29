package pagegenerators

import (
	"bytes"
	"strconv"
	"strings"

	"nostr-static/src/utils"

	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func displayNameOrName(displayName, name string) string {
	if displayName != "" {
		return displayName
	}
	return name
}

// Common CSS styles
const CommonStyles = `
    /* Base styles */
    body {
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
        line-height: 1.6;
        max-width: 800px;
        margin: 0 auto;
        padding: 20px;
        font-size: 16px;
    }

    /* Theme styles */
    body.light {
        background-color: #ffffff;
        color: #333333;
    }

    body.dark {
        background-color: #1a1a1a;
        color: #ffffff;
    }

    body.light a {
        color: #0066cc;
    }

    body.dark a {
        color: #4a9eff;
    }

    /* Feed links styles */
    .feed-links {
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 0.5em;
        padding: 1em;
    }

    .feed-icon {
        width: 16px;
        height: 16px;
    }

    /* Common components */
    .logo {
        text-align: center;
    }

    .logo img {
        max-height: 100px;
        width: auto;
    }

    .article-card {
        border: 1px solid #ddd;
        border-radius: 8px;
        padding: 20px;
        margin-bottom: 20px;
    }

    .article-card h2 {
        margin-top: 0;
        text-align: center;
    }

    .tags {
        display: flex;
        flex-wrap: wrap;
        gap: 0 8px;
        margin-top: 1em;
    }

    .tag {
        display: inline-block;
        padding: 4px 8px;
        border-radius: 4px;
        margin-bottom: 8px;
        font-size: 0.9em;
    }

    .tag a {
        text-decoration: none;
        color: inherit;
    }

    /* Theme-specific tag styles */
    body.light .tag {
        background-color: #f0f0f0;
        border: 1px solid #dee2e6;
        color: #666666;
    }

    body.dark .tag {
        background-color: #1a1a1a;
        border: 1px solid #404040;
        color: #e0e0e0;
    }

    body.index .page-container,
    body.tagspage .page-container {
        display: flex;
        align-items: flex-start;
    }

    .logo-container {
        margin-top: 20px;
        margin-right: 20px;
    }

    .logo-container img {
        max-height: 50px;
        width: auto;
    }

    body.index .main-content,
    body.tagspage .main-content {
        flex: 1;
    }

    /* Image container styles */
    .image-container {
        margin: 20px auto 10px;
        text-align: center;
    }

    .image-container img {
        max-width: 100%;
        height: auto;
        border-radius: 4px;
    }

    /* Index page specific image styles */
    body.index .image-container {
        max-width: 300px;
    }

    /* Theme-specific component styles */
    body.light .article-card {
        background-color: #f8f9fa;
        border: 1px solid #dee2e6;
    }

    body.dark .article-card {
        background-color: #2d2d2d;
        border: 1px solid #404040;
    }

    body.light pre {
        background-color: #f5f5f5 !important;
        border: 1px solid #e9ecef;
    }

    body.dark pre {
        background-color: #2d2d2d !important;
        border: 1px solid #404040;
    }

    body.light .article-card .summary {
        color: #666666;
    }

    body.dark .article-card .summary {
        color: #e0e0e0;
    }

    .time-ago {
        display: inline-block;
        min-width: 70px;
    }

    /* Compact profile styles */
    .compact-profile {
        display: flex;
        align-items: center;
        gap: 8px;
    }

    .compact-profile-link {
        display: flex;
        align-items: center;
        text-decoration: none;
        color: inherit;
        gap: 8px;
    }

    .compact-profile-ago {
        margin-bottom: 1px;
        font-size: 0.8em;
    }

    .compact-profile-picture {
        width: 15px;
        height: 15px;
        border-radius: 50%;
        object-fit: cover;
    }

    .compact-profile-name {
        font-size: 0.9em;
        font-weight: 500;
    }

    /* Theme-specific compact profile styles */
    body.light .compact-profile-link {
        color: #000000;
    }

    body.light .compact-profile-link:hover {
        color: #0066cc;
    }

    body.dark .compact-profile-link {
        color: #e0e0e0;
    }

    body.dark .compact-profile-link:hover {
        color: #4a9eff;
    }

    /* Footer styles */
    .footer {
        margin-top: 2em;
        padding-top: 1em;
        border-top: 1px solid #eee;
        text-align: center;
        font-size: 0.9em;
        color: #666;
    }

    body.dark .footer {
        border-top-color: #333;
        color: #999;
    }

    .footer a {
        text-decoration: none;
    }

    body.light .footer a {
        color: #0066cc;
    }

    body.dark .footer a {
        color: #4a9eff;
    }
`

const CommonResponsiveStyles = `
    @media (max-width: 768px) {
        body {
            font-size: 18px;
            padding: 15px;
        }

        .page-container {
            flex-direction: column;
        }

        .main-content {
            width: 100%;
        }

        .article-card {
            padding: 15px;
        }
    }
`

func convertMarkdownToHTML(markdown string, stripTitle bool) (string, error) {
	if stripTitle {
		// Remove the first heading if it exists
		lines := strings.Split(markdown, "\n")
		if len(lines) > 0 && strings.HasPrefix(strings.TrimSpace(lines[0]), "# ") {
			lines = lines[1:]
		}
		markdown = strings.Join(lines, "\n")
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Common HTML rendering functions
func renderLogo(logo, baseFolder string) HTML {
	if logo == "" {
		return Text("")
	}

	baseFolder = strings.Trim(baseFolder, "/")
	if baseFolder != "" {
		baseFolder = baseFolder + "/"
	}

	return Div(Attr(a.Class("logo")),
		A(Attr(a.Href(baseFolder+"index.html")),
			Img(Attr(
				a.Src(baseFolder+logo),
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

func renderSummaryHTML(summary string) HTML {
	if summary == "" {
		return Text("")
	}
	return P(Attr(a.Class("summary")), Text(summary))
}

func renderCompactProfile(
	authorName,
	authorPicture,
	authorNProfile,
	naddr string,
	createdAt int64,
) HTML {
	if authorName == "" {
		return Text("")
	}

	pictureHTML := utils.Ternary(authorPicture == "",
		Text(""),
		Img(Attr(
			a.Src(authorPicture),
			a.Alt(authorName),
			a.Class("compact-profile-picture"),
		)),
	)

	return Div(Attr(a.Class("compact-profile")),
		A(Attr(
			a.Href("profile/"+authorNProfile+".html"),
			a.Class("compact-profile-link"),
		),
			pictureHTML,
			Span(Attr(a.Class("compact-profile-name")),
				Text(authorName),
			),
		),
		A(Attr(
			a.Href(naddr+".html"),
			a.Class("compact-profile-ago"),
		),
			Span(Attr(
				a.Class("time-ago"),
				a.Dataset("timestamp", strconv.FormatInt(createdAt, 10)),
			)),
		),
	)
}

// Should be used in the head section
func rssFeedLink(fileName string) HTML {
	return Link(Attr(
		a.Rel("alternate"),
		a.Type("application/rss+xml"),
		a.Title("RSS Feed"),
		a.Href(fileName+"-rss.xml"),
	))
}

// Should be used in the head section
func atomFeedLink(fileName string) HTML {
	return Link(Attr(
		a.Rel("alternate"),
		a.Type("application/atom+xml"),
		a.Title("Atom Feed"),
		a.Href(fileName+"-atom.xml"),
	))
}

// Should be after the footer section
func renderFeed(fileName string) HTML {
	return Div(Attr(a.Class("feed-links")),
		Img(Attr(
			a.Src_("data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiA/PjxzdmcgaGVpZ2h0PSIyNCIgdmVyc2lvbj0iMS4xIiB3aWR0aD0iMjQiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgeG1sbnM6Y2M9Imh0dHA6Ly9jcmVhdGl2ZWNvbW1vbnMub3JnL25zIyIgeG1sbnM6ZGM9Imh0dHA6Ly9wdXJsLm9yZy9kYy9lbGVtZW50cy8xLjEvIiB4bWxuczpyZGY9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMiPjxnIHRyYW5zZm9ybT0idHJhbnNsYXRlKDAgLTEwMjguNCkiPjxnPjxwYXRoIGQ9Im00IDEwMzEuNGMtMS4xMDQ2IDAtMiAwLjktMiAydjE2YzAgMS4xIDAuODk1NCAyIDIgMmgxNmMxLjEwNSAwIDItMC45IDItMnYtMTZjMC0xLjEtMC44OTUtMi0yLTJoLTE2eiIgZmlsbD0iI2QzNTQwMCIvPjxwYXRoIGQ9Im00IDJjLTEuMTA0NiAwLTIgMC44OTU0LTIgMnYxNmMwIDEuMTA1IDAuODk1NCAyIDIgMmgxNmMxLjEwNSAwIDItMC44OTUgMi0ydi0xNmMwLTEuMTA0Ni0wLjg5NS0yLTItMmgtMTZ6IiBmaWxsPSIjZTY3ZTIyIiB0cmFuc2Zvcm09InRyYW5zbGF0ZSgwIDEwMjguNCkiLz48cGF0aCBkPSJtNSAxMDM0LjR2Mi4zYzYuNDQzIDAgMTEuNjY3IDUuMiAxMS42NjcgMTEuN2gyLjMzM2MwLTcuOC02LjI2OC0xNC0xNC0xNHptMCA0LjZ2Mi40YzMuODY2IDAgNyAzLjEgNyA3aDIuMzMzYzAtNS4yLTQuMTc4LTkuNC05LjMzMy05LjR6bTIuMDQxNyA1LjNjLTEuMTI3NiAwLTIuMDQxNyAwLjktMi4wNDE3IDJzMC45MTQxIDIuMSAyLjA0MTcgMi4xYzEuMTI3NSAwIDIuMDQxNi0xIDIuMDQxNi0yLjFzLTAuOTE0MS0yLTIuMDQxNi0yeiIgZmlsbD0iI2QzNTQwMCIvPjxwYXRoIGQ9Im01IDEwMzMuNHYyLjNjNi40NDMgMCAxMS42NjcgNS4yIDExLjY2NyAxMS43aDIuMzMzYzAtNy44LTYuMjY4LTE0LTE0LTE0em0wIDQuNnYyLjRjMy44NjYgMCA3IDMuMSA3IDdoMi4zMzNjMC01LjItNC4xNzgtOS40LTkuMzMzLTkuNHptMi4wNDE3IDUuM2MtMS4xMjc2IDAtMi4wNDE3IDAuOS0yLjA0MTcgMnMwLjkxNDEgMi4xIDIuMDQxNyAyLjFjMS4xMjc1IDAgMi4wNDE2LTEgMi4wNDE2LTIuMXMtMC45MTQxLTItMi4wNDE2LTJ6IiBmaWxsPSIjZWNmMGYxIi8+PC9nPjwvZz48L3N2Zz4="),
			a.Alt("RSS Feed"),
			a.Width("16"),
			a.Height("16"),
			a.Class("feed-icon"),
		)),
		A(Attr(
			a.Href(fileName+"-rss.xml"),
			a.Title("RSS"),
			a.Target("_blank"),
		),
			Text("RSS"),
		),
		A(Attr(
			a.Href(fileName+"-atom.xml"),
			a.Title("Atom"),
			a.Target("_blank"),
		),
			Text("Atom"),
		),
	)
}
