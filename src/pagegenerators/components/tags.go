package components

import (
	"strings"

	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
)

func RenderTagsHTML(tags []string, baseFolder string) HTML {
	if len(tags) == 0 {
		return Text("")
	}

	var tagElements []HTML
	for _, tag := range tags {
		baseFolder = strings.Trim(baseFolder, "/")
		if baseFolder != "" {
			baseFolder = baseFolder + "/"
		}
		tagElements = append(tagElements,
			Span(Attr(a.Class("tag")),
				A(Attr(
					a.Href(baseFolder+"tag/"+strings.ToLower(tag)+".html"),
				),
					Text(tag),
				),
			),
		)
	}

	return Div(Attr(a.Class("tags")), tagElements...)
}

const TagsCSS = `
.tags {
	display: flex;
	flex-wrap: wrap;
	gap: 0 8px;
	margin-top: 1em;
	margin-bottom: 1em;
}

.tag {
	display: inline-block;
	padding: 4px 8px;
	border-radius: 4px;
	margin-bottom: 8px;
	font-size: 0.9em;
}

.tag a {
	text-decoration: none;
	color: inherit;
}

body.light .tag {
	background-color: #f0f0f0;
	border: 1px solid #dee2e6;
}	

body.dark .tag {
	background-color: #1a1a1a;
	border: 1px solid #404040;
}
`
