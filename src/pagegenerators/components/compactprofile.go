package components

import (
	"strconv"

	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
)

func RenderCompactProfile(
	authorName,
	authorPicture,
	authorNProfile,
	naddr string,
	createdAt int64,
) HTML {
	if authorName == "" {
		return Text("")
	}

	pictureHTML := Text("")
	if authorPicture != "" {
		pictureHTML = Img(Attr(
			a.Src(authorPicture),
			a.Alt(authorName),
			a.Class("compact-profile-picture"),
		))
	}

	return Div(Attr(a.Class("compact-profile")),
		A(Attr(
			a.Href("profile/"+authorNProfile+".html"),
			a.Class("compact-profile-link"),
		),
			pictureHTML,
			Span(Attr(a.Class("compact-profile-name")),
				Text(authorName),
			),
		),
		A(Attr(
			a.Href(naddr+".html"),
			a.Class("compact-profile-ago"),
		),
			Span(Attr(
				a.Class("time-ago"),
				a.Dataset("timestamp", strconv.FormatInt(createdAt, 10)),
			)),
		),
	)
}

const CompactProfileCSS = `
.compact-profile {
    display: flex;
    align-items: center;
    gap: 8px;
}

.compact-profile-link {
    display: flex;
    align-items: center;
    text-decoration: none;
    color: inherit;
    gap: 8px;
}

.compact-profile-ago {
    margin-bottom: 1px;
    font-size: 0.8em;
}

.compact-profile-picture {
    width: 15px;
    height: 15px;
    border-radius: 50%;
    object-fit: cover;
}

.compact-profile-name {
    font-size: 0.9em;
    font-weight: 500;
}

/* Theme-specific compact profile styles */
body.light .compact-profile-link {
    color: #000000;
}

body.light .compact-profile-link:hover {
    color: #0066cc;
}

body.dark .compact-profile-link {
    color: #e0e0e0;
}

body.dark .compact-profile-link:hover {
    color: #4a9eff;
}

.time-ago {
    display: inline-block;
    min-width: 70px;
}
`
