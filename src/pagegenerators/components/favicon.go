package components

import (
	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
)

// RenderFaviconLinks returns the HTML elements for favicon links
func RenderFaviconLinks(faviconDir string) []HTML {
	return []HTML{
		Link(Attr(a.Rel_("icon"), a.Href(faviconDir+"/favicon-96x96.png"), a.Sizes("96x96"))),
		Link(Attr(a.Rel_("icon"), a.Href(faviconDir+"/favicon.svg"), a.Type("image/svg+xml"))),
		Link(Attr(a.Rel_("shortcut icon"), a.Href(faviconDir+"/favicon.ico"))),
		Link(Attr(a.Rel_("apple-touch-icon"), a.Href(faviconDir+"/apple-touch-icon.png"), a.Sizes("180x180"))),
		Link(Attr(a.Rel_("manifest"), a.Href(faviconDir+"/site.webmanifest"))),
	}
}
