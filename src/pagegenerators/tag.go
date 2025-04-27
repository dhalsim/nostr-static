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
    <style>
        ` + CommonStyles + `
				` + ResponsiveStyles + `
    </style>
</head>
<body class="{{.Color}} tags">
    <div class="page-container">
        <div class="logo-container">
            {{renderLogo .Logo "../"}}
        </div>
        <div class="main-content">
            <h1>Articles tagged with "{{.Tag}}"</h1>
            {{range .Articles}}
            <div class="article-card">
								{{renderCompactProfile 
								  .AuthorName 
									.Nprofile 
									.Naddr
									.AuthorPicture 
									.CreatedAt 
									"../"
								}}
                {{renderImage .Image .Title .Naddr "../"}}
                <h2><a href="../{{.Naddr}}.html">{{.Title}}</a></h2>
                {{renderSummary .Summary}}
                {{renderTags .Tags "../"}}
            </div>
            {{end}}
            {{renderFooter}}
        </div>
    </div>
    <script src="/output/static/js/time-ago.js"></script>
</body>
</html>`

type TagData struct {
	Tag      string
	Color    string
	Logo     string
	Articles []TagArticleData
}

type TagArticleData struct {
	Naddr         string
	Title         string
	Summary       string
	Image         string
	Tags          []string
	AuthorName    string
	Nprofile      string
	AuthorPicture string
	CreatedAt     int64
}

type generateTagPagesParams struct {
	Events           []types.Event
	Profiles         map[string]types.Event
	OutputDir        string
	Layout           types.Layout
	EventIDToNaddr   map[string]string
	PubkeyToNProfile map[string]string
}

func NewGenerateTagPagesParams(
	events []types.Event,
	profiles map[string]types.Event,
	outputDir string,
	layout types.Layout,
	eventIDToNaddr map[string]string,
	pubkeyToNProfile map[string]string,
) generateTagPagesParams {
	return generateTagPagesParams{
		Events:           events,
		Profiles:         profiles,
		OutputDir:        outputDir,
		Layout:           layout,
		EventIDToNaddr:   eventIDToNaddr,
		PubkeyToNProfile: pubkeyToNProfile,
	}
}

func GenerateTagPages(params generateTagPagesParams) error {
	// Create a map to track tags and their associated events
	tagMap := make(map[string][]types.Event)

	// Populate the tag map
	for _, event := range params.Events {
		// Extract tags from the event
		for _, tagArray := range event.Tags {
			if len(tagArray) > 1 && tagArray[0] == "t" {
				tag := strings.ToLower(tagArray[1])
				tagMap[tag] = append(tagMap[tag], event)
			}
		}
	}

	// Create tag directory if it doesn't exist
	tagDir := filepath.Join(params.OutputDir, "tag")
	if err := os.MkdirAll(tagDir, 0755); err != nil {
		return err
	}

	// Generate a page for each tag
	for tag, tagEvents := range tagMap {
		data := TagData{
			Tag:      tag,
			Color:    params.Layout.Color,
			Logo:     params.Layout.Logo,
			Articles: make([]TagArticleData, len(tagEvents)),
		}

		for i, event := range tagEvents {
			// Extract metadata from tags
			metadata := ExtractArticleMetadata(event.Tags)

			parsedProfile := parseProfile(params.Profiles[event.PubKey])

			data.Articles[i] = TagArticleData{
				Naddr:   params.EventIDToNaddr[event.ID],
				Title:   metadata.Title,
				Summary: metadata.Summary,
				Image:   metadata.Image,
				Tags:    metadata.Tags,
				AuthorName: displayNameOrName(
					parsedProfile.DisplayName,
					parsedProfile.Name,
				),
				AuthorPicture: parsedProfile.Picture,
				CreatedAt:     event.CreatedAt,
				Nprofile:      params.PubkeyToNProfile[event.PubKey],
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
