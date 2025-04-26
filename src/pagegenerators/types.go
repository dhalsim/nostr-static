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
