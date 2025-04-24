package pagegenerators

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"nostr-static/src/types"
)

const tagTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Tag: {{.Tag}}</title>
    {{.ThemeCSS}}
    <style>
        ` + CommonStyles + `
        .page-container {
            display: flex;
            gap: 2em;
            align-items: flex-start;
        }
        .logo-container {
            flex: 0 0 200px;
            position: sticky;
            top: 20px;
        }
        .main-content {
            flex: 1;
        }
    </style>
</head>
<body>
    <div class="page-container">
        <div class="logo-container">
            {{renderLogo .Logo "../"}}
        </div>
        <div class="main-content">
            <h1>Articles tagged with "{{.Tag}}"</h1>
            {{range .Articles}}
            <div class="article-card">
                {{renderImage .Image .Title}}
                <h2><a href="../{{.ID}}.html">{{.Title}}</a></h2>
                {{renderSummary .Summary}}
                {{renderTags .Tags "../"}}
            </div>
            {{end}}
        </div>
    </div>
</body>
</html>`

func GenerateTagPages(events []types.Event, outputDir string, themeColor string, logo string) error {
	// Create a map to track tags and their associated events
	tagMap := make(map[string][]types.Event)

	// Populate the tag map
	for _, event := range events {
		// Extract tags from the event
		for _, tagArray := range event.Tags {
			if len(tagArray) > 1 && tagArray[0] == "t" {
				tag := strings.ToLower(tagArray[1])
				tagMap[tag] = append(tagMap[tag], event)
			}
		}
	}

	// Create tag directory if it doesn't exist
	tagDir := filepath.Join(outputDir, "tag")
	if err := os.MkdirAll(tagDir, 0755); err != nil {
		return err
	}

	// Generate a page for each tag
	for tag, tagEvents := range tagMap {
		data := TagData{
			Tag:      tag,
			ThemeCSS: GetThemeCSS(themeColor),
			Logo:     logo,
			Articles: make([]struct {
				ID      string
				Title   string
				Summary string
				Image   string
				Tags    []string
			}, len(tagEvents)),
		}

		for i, event := range tagEvents {
			// Extract metadata from tags
			metadata := ExtractArticleMetadata(event.Tags)

			data.Articles[i] = struct {
				ID      string
				Title   string
				Summary string
				Image   string
				Tags    []string
			}{
				ID:      event.ID,
				Title:   metadata.Title,
				Summary: metadata.Summary,
				Image:   metadata.Image,
				Tags:    metadata.Tags,
			}
		}

		tmpl, err := template.New("tag").Funcs(ComponentFuncs).Parse(tagTemplate)
		if err != nil {
			return err
		}

		outputFile := filepath.Join(tagDir, strings.ToLower(tag)+".html")
		f, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := tmpl.Execute(f, data); err != nil {
			return err
		}
	}

	return nil
}
