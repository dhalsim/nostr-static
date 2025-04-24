package pagegenerators

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Component functions for templates
var ComponentFuncs = template.FuncMap{
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
	"renderImage": func(image, title string) template.HTML {
		if image == "" {
			return ""
		}
		return template.HTML(`
		<div class="image-container">
			<img src="` + image + `" alt="` + title + `">
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

// Common data structures
type ArticleData struct {
	Title    string
	Content  template.HTML
	ThemeCSS template.HTML
	Summary  string
	Tags     []string
	Logo     string
}

type TagData struct {
	Tag      string
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

// Common CSS styles
const CommonStyles = `
    body {
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
        line-height: 1.6;
        max-width: 800px;
        margin: 0 auto;
        padding: 20px;
    }
    .logo {
        text-align: center;
        margin-bottom: 2em;
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
`

// Theme CSS functions
func GetThemeCSS(color string) template.HTML {
	if color == "dark" {
		return template.HTML(`
		<style>
			body {
				background-color: #1a1a1a;
				color: #ffffff;
			}
			a {
				color: #4a9eff;
			}
			.article-card {
				background-color: #2d2d2d;
				border: 1px solid #404040;
			}
			pre {
				background-color: #2d2d2d !important;
				border: 1px solid #404040;
			}
			.article-card .tag {
				display: inline-block;
				background-color: #1a1a1a;
				border: 1px solid #404040;
				color: #e0e0e0;
				padding: 4px 8px;
				border-radius: 4px;
				margin-bottom: 8px;
				font-size: 0.9em;
			}
			.article-card .tags {
				display: flex;
				flex-wrap: wrap;
				gap: 8px;
				margin-top: 1em;
			}
			.article-card .image-container {
				margin: 0 auto 20px;
				max-width: 600px;
				text-align: center;
			}
			.article-card .image-container img {
				max-width: 100%;
				height: auto;
				border-radius: 4px;
			}
			.article-card .summary {
				color: #e0e0e0;
			}
		</style>
		`)
	}

	return template.HTML(`
	<style>
		body {
			background-color: #ffffff;
			color: #333333;
		}
		a {
			color: #0066cc;
		}
		.article-card {
			background-color: #f8f9fa;
			border: 1px solid #dee2e6;
		}
		pre {
			background-color: #f5f5f5 !important;
			border: 1px solid #e9ecef;
		}
		.article-card .tag {
			display: inline-block;
			background-color: #f0f0f0;
			border: 1px solid #dee2e6;
			color: #666666;
			padding: 4px 8px;
			border-radius: 4px;
			margin-bottom: 8px;
			font-size: 0.9em;
		}
		.article-card .tags {
			display: flex;
			flex-wrap: wrap;
			gap: 8px;
			margin-top: 1em;
		}
		.article-card .image-container {
			margin: 0 auto 20px;
			max-width: 600px;
			text-align: center;
		}
		.article-card .image-container img {
			max-width: 100%;
			height: auto;
			border-radius: 4px;
		}
		.article-card .summary {
			color: #666666;
		}
	</style>
	`)
}

func convertMarkdownToHTML(markdown string) (string, error) {
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

// ArticleMetadata represents the metadata extracted from an event
type ArticleMetadata struct {
	Title   string
	Summary string
	Image   string
	Tags    []string
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
