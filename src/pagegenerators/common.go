package pagegenerators

import (
	"bytes"
	"html/template"
	"strconv"
	"strings"

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

// Component functions for templates
var ComponentFuncs = template.FuncMap{
	"displayNameOrName": displayNameOrName,
	"renderFooter": func() template.HTML {
		return template.HTML(`
		<div class="footer">
			Built with <a target="_blank" href="https://github.com/dhalsim/nostr-static">nostr-static</a>
		</div>
		`)
	},
	"renderCompactProfile": func(
		name string,
		nprofile string,
		naddr string,
		picture string,
		createdAt int64,
		baseFolder string,
	) template.HTML {
		if name == "" {
			return ""
		}

		// Normalize the path by removing any leading/trailing slashes
		baseFolder = strings.Trim(baseFolder, "/")
		if baseFolder != "" {
			baseFolder = baseFolder + "/"
		}

		var pictureHTML string
		if picture != "" {
			pictureHTML = `<img src="` + picture + `" alt="` + name + `" class="compact-profile-picture">`
		}

		return template.HTML(`
		<div class="compact-profile">
			<a href="` + baseFolder + `profile/` + nprofile + `.html" class="compact-profile-link">
				` + pictureHTML + `
				<span class="compact-profile-name">` + name + `</span>
            </a>
            <a href="` + baseFolder + naddr + `.html" class="compact-profile-ago">
              <span class="time-ago" data-timestamp="` + strconv.FormatInt(createdAt, 10) + `"></span>
            </a>
		</div>
		`)
	},
	"renderLogo": func(logo string, baseFolder string) template.HTML {
		if logo == "" {
			return ""
		}
		// Normalize the path by removing any leading/trailing slashes
		baseFolder = strings.Trim(baseFolder, "/")
		if baseFolder != "" {
			baseFolder = baseFolder + "/"
		}
		return template.HTML(`
		<div class="logo">
			<a href="` + baseFolder + `index.html">
				<img src="` + baseFolder + logo + `" alt="Logo">
			</a>
		</div>
		`)
	},
	"renderAuthor": func(name, nprofile, picture, baseFolder string) template.HTML {
		if name == "" {
			return ""
		}

		// Normalize the path by removing any leading/trailing slashes
		baseFolder = strings.Trim(baseFolder, "/")
		if baseFolder != "" {
			baseFolder = baseFolder + "/"
		}

		var pictureHTML string
		if picture != "" {
			pictureHTML = `<img src="` + picture + `" alt="` + name + `" class="author-picture">`
		}

		return template.HTML(`
		<div class="author">
			<a href="` + baseFolder + `profile/` + nprofile + `.html" class="author-link">
				` + pictureHTML + `
				<span class="author-name">` + name + `</span>
			</a>
		</div>
		`)
	},
	"renderAgo": func(ago string) template.HTML {
		return template.HTML(`<span class="ago">` + ago + `</span>`)
	},
	"renderImage": func(image, alt, imageLink, baseFolder string) template.HTML {
		if image == "" {
			return ""
		}

		if imageLink == "" {
			return template.HTML(`
			<div class="image-container">
				<img src="` + image + `" alt="` + alt + `">
			</div>
			`)
		}

		return template.HTML(`
			<div class="image-container">
				<a href="` + baseFolder + imageLink + `.html">
					<img src="` + image + `" alt="` + alt + `">
				</a>
			</div>
		`)
	},
	"renderSummary": func(summary string) template.HTML {
		if summary == "" {
			return ""
		}
		return template.HTML(`
		<p class="summary">` + summary + `</p>
		`)
	},
	"renderTags": func(tags []string, baseFolder string) template.HTML {
		if len(tags) == 0 {
			return ""
		}
		var html string
		for _, tag := range tags {
			// Normalize the path by removing any leading/trailing slashes
			baseFolder = strings.Trim(baseFolder, "/")
			if baseFolder != "" {
				baseFolder = baseFolder + "/"
			}
			html += `<span class="tag"><a href="` + baseFolder + `tag/` + strings.ToLower(tag) + `.html">` + tag + `</a></span>`
		}
		return template.HTML(`
		<div class="tags">
			` + html + `
		</div>
		`)
	},
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
        gap: 8px;
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

    /* Page-specific styles */
    body.article .page-container, body.profile .page-container {
        display: flex;
        align-items: flex-start;
        max-width: 1200px;
        margin: 0 auto;
    }

    body.article .logo-container, body.profile .logo-container {
        flex: 0 0 200px;
        position: sticky;
        top: 20px;
    }

    body.article .main-content, body.profile .main-content {
        flex: 1;
        max-width: 800px;
    }

    body.article img, body.profile img {
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

    body.index .page-container,
    body.tags .page-container {
        display: flex;
        align-items: flex-start;
    }

    body.index .logo-container,
    body.tags .logo-container {
        flex: 0 0 200px;
        position: sticky;
        top: 20px;
    }

    body.index .main-content,
    body.tags .main-content {
        flex: 1;
    }

    /* Image container styles */
    .image-container {
        margin: 0 auto 20px;
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

    /* Author styles */
    .author {
        display: flex;
        align-items: center;
        margin-bottom: 1em;
    }

    .author-link {
        display: flex;
        align-items: center;
        text-decoration: none;
        color: inherit;
        gap: 10px;
    }

    .author-picture {
        width: 40px;
        height: 40px;
        border-radius: 50%;
        object-fit: cover;
    }

    .author-name {
        font-weight: 500;
    }

    /* Theme-specific author styles */
    body.light .author-website {
        color: #000000;
    }

    body.light .author-website:hover {
        color: #0066cc;
    }

    body.dark .author-website {
        color: #e0e0e0;
        filter: invert(1);
    }

    body.dark .author-website:hover {
        color: #4a9eff;
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
        margin-bottom: 2px;
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

const ResponsiveStyles = `
    @media (max-width: 768px) {
        body {
            font-size: 18px;
            padding: 15px;
        }

        body.article .page-container,
        body.index .page-container,
        body.tags .page-container,
        body.profile .page-container {
            flex-direction: column;
        }

        body.article .logo-container,
        body.index .logo-container,
        body.tags .logo-container,
        body.profile .logo-container {
            flex: none;
            position: static;
        }

        body.article .main-content,
        body.index .main-content,
        body.tags .main-content,
        body.profile .main-content {
            width: 100%;
        }

        .article-card {
            padding: 15px;
        }

        .logo img {
            max-height: 50px;
            width: auto;
        }

        .logo-container {
            text-align: center;
            padding: 10px 0;
        }

        .article-header {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 15px;
        }

        .author {
            flex-direction: column;
            align-items: center;
            text-align: center;
            margin-bottom: 15px;
        }

        .author-link {
            flex-direction: column;
            gap: 8px;
        }

        .profile-header {
            flex-direction: column;
            align-items: center;
            text-align: center;
            gap: 15px;
        }

        .profile-header-left {
            align-items: center;
        }

        .profile-links {
            flex-direction: column;
            align-items: center;
            gap: 10px;
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

// ExtractArticleMetadata extracts metadata from event tags
func ExtractArticleMetadata(tags [][]string) ArticleMetadata {
	var metadata ArticleMetadata
	for _, tag := range tags {
		if len(tag) < 2 {
			continue
		}

		switch tag[0] {
		case "title":
			metadata.Title = tag[1]
		case "summary":
			metadata.Summary = tag[1]
		case "image":
			metadata.Image = tag[1]
		case "t":
			metadata.Tags = append(metadata.Tags, tag[1])
		}
	}
	return metadata
}
