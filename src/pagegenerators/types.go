package pagegenerators

import "html/template"

// ArticleData represents the data needed for article templates
type ArticleData struct {
	Title   string
	Content template.HTML
	Color   string
	Summary string
	Tags    []string
	Logo    string
	Image   string
	Naddr   string
}

// ArticleMetadata represents the metadata extracted from an event
type ArticleMetadata struct {
	Title   string
	Summary string
	Image   string
	Tags    []string
}
