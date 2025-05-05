package components

import (
	"nostr-static/src/utils"

	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
)

func RenderDropdownScript() HTML {
	return Script(Attr(a.Src("/static/js/dropdown.js")), JavaScript(""))
}

func RenderNostrLinks(naddr, nprofile, nostrLinks string) HTML {
	showLinks := (naddr != "" || nprofile != "") && nostrLinks != ""
	if !showLinks {
		return Text("")
	}

	return Div(Attr(a.Class("dropdown")),
		Button(Attr(
			a.Class("dropdown-button"),
			a.Type("button"),
		),
			Text("⋯"),
		),
		Div(Attr(a.Class("dropdown-content")),
			utils.Ternary(
				naddr != "",
				A(Attr(
					a.Href(`https://`+nostrLinks+`/`+naddr),
					a.Class("nostr-link"),
					a.Target("_blank"),
				),
					Text("Open Article in Nostr"),
				),
				Text(""),
			),
			utils.Ternary(
				nprofile != "",
				A(Attr(
					a.Href(`https://`+nostrLinks+`/`+nprofile),
					a.Class("nostr-link"),
					a.Target("_blank"),
				),
					Text("Open Author in Nostr"),
				),
				Text(""),
			),
		),
	)
}

var DotMenuCSS = `
.dropdown {
	position: relative;
	display: inline-block;
}

.dropdown-button {
	background-color: #f0f0f0;
	border: none;
	border-radius: 50%;
	width: 20px;
	height: 20px;
	margin-left: 10px;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 0.8em;
	cursor: pointer;
	color: #666;
	transition: background-color 0.2s;
}

.dropdown-button:hover {
	background-color: #e0e0e0;
}

.dropdown-content {
	display: none;
	position: absolute;
	right: -10px;
	top: 100%;
	background-color: white;
	min-width: 180px;
	box-shadow: 0px 8px 16px 0px rgba(0,0,0,0.1);
	z-index: 1;
	border-radius: 8px;
	margin-top: 8px;
	border: 1px solid #eee;
}

.dropdown.show .dropdown-content {
	display: block;
}

.dropdown-content a {
	font-size: 0.8em;
	color: #333;
	padding: 12px 16px;
	text-decoration: none;
	display: block;
	transition: background-color 0.2s;
	border-bottom: 1px solid #eee;
	border-radius: 8px;
}

.dropdown-content a:last-child {
	border-bottom: none;
}

.dropdown-content a:hover {
	background-color: #f5f5f5;
}
`
