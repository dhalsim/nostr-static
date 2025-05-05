package pagegenerators

import (
	"bytes"
	"strings"

	. "github.com/julvo/htmlgo"
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
const CommonCSS = `
/* Base styles */
body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    line-height: 1.6;
    max-width: 800px;
	margin: 0 auto;
	padding: 20px;
	font-size: 16px;
    color: #333;
    background-color: #fff;
}

/* Theme styles */
body.light {
    background-color: #ffffff;
    color: #333333;
}

body.dark {
    background-color: #1a1a1a;
    color: #e0e0e0;
}

/* Container */
.container {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
}

/* Links */
a {
    color: #0066cc;
    text-decoration: none;
}

body.light a {
    color: #0066cc;
}

body.dark a {
    color: #4a9eff;
}

a:hover {
    text-decoration: underline;
}

/* Lists */
ul, ol {
    padding-left: 2em;
}

/* Code blocks */
pre {
    background-color: #f5f5f5;
    padding: 1em;
    border-radius: 4px;
    overflow-x: auto;
}

body.dark pre {
    background-color: #2d2d2d;
}

code {
    font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
    font-size: 0.9em;
}

/* Blockquotes */
blockquote {
    margin: 1em 0;
    padding: 0.5em 1em;
    border-left: 4px solid #ddd;
    background-color: #f9f9f9;
}

body.dark blockquote {
    border-left-color: #404040;
    background-color: #2d2d2d;
}

/* Tables */
table {
    border-collapse: collapse;
    width: 100%;
    margin: 1em 0;
}

th, td {
    padding: 0.5em;
    border: 1px solid #ddd;
}

body.dark th, body.dark td {
    border-color: #404040;
}

th {
    background-color: #f5f5f5;
}

body.dark th {
    background-color: #2d2d2d;
}

/* Horizontal rule */
hr {
    border: none;
    border-top: 1px solid #ddd;
    margin: 2em 0;
}

body.dark hr {
    border-top-color: #404040;
}

/* Theme toggle button */
.theme-toggle {
    position: fixed;
    top: 20px;
    right: 20px;
    padding: 8px 16px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    transition: background-color 0.3s;
}

body.light .theme-toggle {
    background-color: #f0f0f0;
    color: #333;
}

body.dark .theme-toggle {
    background-color: #404040;
    color: #fff;
}

.theme-toggle:hover {
    opacity: 0.9;
}

/* Responsive design */
@media (max-width: 600px) {
    .container {
        padding: 10px;
    }
    
    .theme-toggle {
        top: 10px;
        right: 10px;
    }
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

func renderTitleHTML(title string) HTML {
	return H1_(Text(title))
}
