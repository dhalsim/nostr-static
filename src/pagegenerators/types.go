package pagegenerators

import "html/template"

// ArticleData represents the data needed for article templates
type ArticleData struct {
	Title     string
	Content   template.HTML
	Color     string
	Summary   string
	Tags      []string
	Logo      string
	Image     string
	ImageLink string
}

// TagData represents the data needed for tag templates
type TagData struct {
	Tag      string
	Color    string
	Logo     string
	Articles []struct {
		ID        string
		Title     string
		Summary   string
		Image     string
		ImageLink string
		Tags      []string
	}
}

// ArticleMetadata represents the metadata extracted from an event
type ArticleMetadata struct {
	Title     string
	Summary   string
	Image     string
	ImageLink string
	Tags      []string
}
