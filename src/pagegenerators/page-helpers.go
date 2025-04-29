package pagegenerators

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
