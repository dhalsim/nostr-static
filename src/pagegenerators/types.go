package pagegenerators

// ArticleMetadata represents the metadata extracted from an event
type ArticleMetadata struct {
	Title   string
	Summary string
	Image   string
	Tags    []string
}

type ParsedProfile struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	About       string `json:"about"`
	Picture     string `json:"picture"`
	Website     string `json:"website"`
	Nip05       string `json:"nip05"`
	Lud16       string `json:"lud16"`
}

// TagData represents all data needed for tag templates
type TagData struct {
	BaseFolder string
	Tag        string
	Color      string
	Logo       string
	Articles   []TagArticleData
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
