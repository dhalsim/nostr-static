package components

import (
	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
)

func RenderFooter() HTML {
	return Div(Attr(a.Class("footer")),
		Text("Built with "),
		A(Attr(
			a.Href("https://github.com/dhalsim/nostr-static"),
			a.Target("_blank"),
		),
			Text("nostr-static"),
		),
	)
}

const FooterCSS = `
.footer {
    margin-top: 2em;
    padding-top: 1em;
    border-top: 1px solid #eee;
    text-align: center;
    font-size: 0.9em;
    color: #666;
}

body.dark .footer {
    border-top-color: #333;
    color: #999;
}

.footer a {
    text-decoration: none;
}

body.light .footer a {
    color: #0066cc;
}

body.dark .footer a {
    color: #4a9eff;
}
`
