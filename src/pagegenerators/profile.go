package pagegenerators

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"nostr-static/src/types"
)

const profileTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Profile: {{.Name}}</title>
    <style>
        ` + CommonStyles + `
        .profile-picture {
            width: 50px;
            height: 50px;
						border-radius: 50%;
        }
        .profile-links {
            margin-top: 10px;
            display: flex;
            gap: 15px;
            align-items: center;
        }
        .profile-links a {
            color: inherit;
            text-decoration: none;
            display: flex;
            align-items: center;
            gap: 5px;
        }
        .profile-links a:hover {
            text-decoration: underline;
        }
        .verified-badge {
            display: inline-flex;
            align-items: center;
            gap: 5px;
        }
        .verified-badge.verified {
            color: #22c55e;
        }
        .verified-badge.unverified {
            color: #ef4444;
        }

				.profile-header {
					display: flex;
					flex-direction: column;
					gap: 10px;
				}

				.profile-header-left {
						display: flex;
						align-content: space-around;
						flex-wrap: wrap;
						align-items: flex-end;
				}

				.profile-header-left h2 {
					margin: 0;
				}` + ResponsiveStyles + `
    </style>
</head>
<body class="{{.Color}} profile">
    <div class="page-container">
        <div class="logo-container">
            {{renderLogo .}}
        </div>
        <div class="main-content">
            <div class="profile-header">
						    <div class="profile-header-left">
									{{renderProfilePicture .}}
									{{renderProfileName .}}
								</div>
								{{renderProfileAbout .}}
								{{renderProfileLinks .}}
            </div>
            <div class="profile-articles">
                <h2>Articles of {{displayNameOrName .DisplayName .Name}}</h2>
                {{range .Articles}}
                <div class="article-card">
                    {{renderImage .Image .Title .Naddr "../"}}
                    <h3><a href="../{{.Naddr}}.html">{{.Title}}</a></h3>
                    {{renderSummary .Summary}}
                    {{renderTags .Tags "../"}}
                </div>
                {{end}}
            </div>
            {{RenderFooter}}
        </div>
    </div>
    {{RenderTimeAgoScript}}
</body>
</html>`

type profileData struct {
	BaseFolder    string
	Color         string
	Logo          string
	Nprofile      string
	PubKey        string
	Name          string
	About         string
	Picture       string
	Website       string
	DisplayName   string
	Nip05         string
	Nip05Verified bool
	Lud16         string
	Articles      []ProfileArticleData
}

func NewProfileData(
	baseFolder string,
	color string,
	logo string,
	nprofile string,
	pubkey string,
	name string,
	about string,
	picture string,
	website string,
	displayName string,
	nip05 string,
	nip05Verified bool,
	lud16 string,
	articles []ProfileArticleData,
) profileData {
	return profileData{
		BaseFolder:    baseFolder,
		Color:         color,
		Logo:          logo,
		Nprofile:      nprofile,
		PubKey:        pubkey,
		Name:          name,
		About:         about,
		Picture:       picture,
		Website:       website,
		DisplayName:   displayName,
		Nip05:         nip05,
		Nip05Verified: nip05Verified,
		Lud16:         lud16,
		Articles:      articles,
	}
}

type ProfileArticleData struct {
	Naddr     string
	Title     string
	Summary   string
	Image     string
	Tags      []string
	CreatedAt int64
}

type generateProfilePagesParams struct {
	BaseFolder       string
	Profiles         map[string]types.Event
	Events           []types.Event
	OutputDir        string
	Layout           types.Layout
	PubkeyToNProfile map[string]string
	EventIDToNaddr   map[string]string
}

func NewGenerateProfilePagesParams(
	baseFolder string,
	profiles map[string]types.Event,
	events []types.Event,
	outputDir string,
	layout types.Layout,
	pubkeyToNProfile map[string]string,
	eventIDToNaddr map[string]string,
) generateProfilePagesParams {
	return generateProfilePagesParams{
		BaseFolder:       baseFolder,
		Profiles:         profiles,
		Events:           events,
		OutputDir:        outputDir,
		Layout:           layout,
		PubkeyToNProfile: pubkeyToNProfile,
		EventIDToNaddr:   eventIDToNaddr,
	}
}

func GenerateProfilePages(params generateProfilePagesParams) error {
	// Create profile directory if it doesn't exist
	profileDir := filepath.Join(params.OutputDir, "profile")
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return err
	}

	// Create a map to track articles by author
	authorArticles := make(map[string][]types.Event)
	for _, event := range params.Events {
		authorArticles[event.PubKey] = append(authorArticles[event.PubKey], event)
	}

	// Generate a page for each profile
	for pubkey, profileEvent := range params.Profiles {
		// Parse profile metadata
		parsedProfile := parseProfile(profileEvent)

		// Verify Nip05 if present
		nip05Verified := false
		if parsedProfile.Nip05 != "" {
			parts := strings.Split(parsedProfile.Nip05, "@")
			if len(parts) == 2 {
				username := parts[0]
				domain := parts[1]
				url := fmt.Sprintf("https://%s/.well-known/nostr.json", domain)

				log.Println("requesting nip05 url: ", url)

				resp, err := http.Get(url)
				if err == nil {
					defer resp.Body.Close()
					var data struct {
						Names map[string]string `json:"names"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
						if pubkey, ok := data.Names[username]; ok {
							nip05Verified = pubkey == profileEvent.PubKey
						}
					}
				}
			}
		}

		// Add articles for this profile
		articleEvents, ok := authorArticles[pubkey]
		if !ok {
			return fmt.Errorf("no articles found for profile: %s", pubkey)
		}

		articles := make([]ProfileArticleData, len(articleEvents))

		for i, event := range articleEvents {
			metadata := ExtractArticleMetadata(event.Tags)

			articles[i] = ProfileArticleData{
				Naddr:     params.EventIDToNaddr[event.ID],
				Title:     metadata.Title,
				Summary:   metadata.Summary,
				Image:     metadata.Image,
				Tags:      metadata.Tags,
				CreatedAt: event.CreatedAt,
			}
		}

		data := NewProfileData(
			params.BaseFolder,
			params.Layout.Color,
			params.Layout.Logo,
			params.PubkeyToNProfile[pubkey],
			pubkey,
			parsedProfile.Name,
			parsedProfile.About,
			parsedProfile.Picture,
			parsedProfile.Website,
			parsedProfile.DisplayName,
			parsedProfile.Nip05,
			nip05Verified,
			parsedProfile.Lud16,
			articles,
		)

		funcs := template.FuncMap{
			"RenderTimeAgoScript": RenderTimeAgoScript,
			"RenderFooter":        RenderFooter,
			"renderLogo": func(data profileData) template.HTML {
				return RenderLogo(data.Logo, data.BaseFolder)
			},
			"renderProfilePicture": func(data profileData) template.HTML {
				if data.Picture == "" {
					return ""
				}

				return template.HTML(
					`<img src="` + data.Picture + `" alt="` + data.Name + `" class="profile-picture">`,
				)
			},
			"renderProfileName": func(data profileData) template.HTML {
				return template.HTML(
					`<h2>` + displayNameOrName(data.DisplayName, data.Name) + `</h2>`,
				)
			},
			"renderProfileAbout": func(data profileData) template.HTML {
				if data.About == "" {
					return ""
				}

				return template.HTML(`<p class="profile-about">` + data.About + `</p>`)
			},
			"renderProfileLinks": func(data profileData) template.HTML {
				return template.HTML(`<div class="profile-links">` +
					renderWebsite(data) +
					renderNip05(data) +
					renderLud16(data) +
					`</div>`)
			},
			"renderProfileArticles": func(data profileData) template.HTML {
				return template.HTML(`
				  <div class="profile-articles">` +
					`<h2>Articles of ` + displayNameOrName(data.DisplayName, data.Name) + `</h2>
					  {{range .Articles}}
               {{renderArticle .}}
						{{end}}
					</div>`,
				)
			},
			"renderArticle": func(data ProfileArticleData) template.HTML {
				return template.HTML(`
					<div class="article-card">
						{{renderImage .Image .Title .Naddr "../"}}
					</div>
				`)
			},
			"renderImage": func(data ProfileArticleData) template.HTML {
				return template.HTML(`
					<img src="` + data.Image + `" alt="` + data.Title + `" class="article-image">
				`)
			},
		}

		tmpl, err := template.New("profile").Funcs(funcs).Parse(profileTemplate)
		if err != nil {
			return err
		}

		outputFile := filepath.Join(profileDir, data.Nprofile+".html")
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

func renderWebsite(data profileData) template.HTML {
	if data.Website == "" {
		return ""
	}

	return template.HTML(
		`<a href="` + data.Website + `" target="_blank" rel="noopener noreferrer">
		  <img src="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9ImN1cnJlbnRDb2xvciIgc3Ryb2tlLXdpZHRoPSIyIiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiIGNsYXNzPSJsdWNpZGUgbHVjaWRlLWxpbmsyLWljb24gbHVjaWRlLWxpbmstMiI+PHBhdGggZD0iTTkgMTdIN0E1IDUgMCAwIDEgNyA3aDIiLz48cGF0aCBkPSJNMTUgN2gyYTUgNSAwIDEgMSAwIDEwaC0yIi8+PGxpbmUgeDE9IjgiIHgyPSIxNiIgeTE9IjEyIiB5Mj0iMTIiLz48L3N2Zz4=" alt="Website" class="author-website" width="16" height="16">
			` + data.Website + `
		</a>`,
	)
}

func renderNip05(data profileData) template.HTML {
	if data.Nip05 == "" {
		return ""
	}

	badgeClass := ternary(data.Nip05Verified, "verified", "unverified")

	badgeIcon := ternary(data.Nip05Verified,
		`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>`,
		`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>`,
	)

	return template.HTML(`<span class="verified-badge ` + badgeClass + `">` +
		badgeIcon +
		data.Nip05 +
		`</span>`)
}

func renderLud16(data profileData) template.HTML {
	if data.Lud16 == "" {
		return ""
	}

	return template.HTML(`<a href="lightning:` + data.Lud16 + `">
	    ⚡ ` + data.Lud16 +
		`</a>`)
}

func parseProfile(event types.Event) ParsedProfile {
	var profile ParsedProfile
	if err := json.Unmarshal([]byte(event.Content), &profile); err != nil {
		return ParsedProfile{}
	}

	return profile
}
