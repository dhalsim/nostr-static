package components

import (
	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
)

// Should be used in the head section
func RenderFeedLinks(fileName string) HTML {
	return Link(Attr(
		a.Rel("alternate"),
		a.Type("application/rss+xml"),
		a.Title("RSS Feed"),
		a.Href(fileName+"-rss.xml"),
	))
}

// Should be used in the head section
func RenderAtomFeedLink(fileName string) HTML {
	return Link(Attr(
		a.Rel("alternate"),
		a.Type("application/atom+xml"),
		a.Title("Atom Feed"),
		a.Href(fileName+"-atom.xml"),
	))
}

// Should be after the footer section
func RenderFeed(fileName string) HTML {
	return Div(Attr(a.Class("feed-links")),
		Img(Attr(
			a.Src_("data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiA/PjxzdmcgaGVpZ2h0PSIyNCIgdmVyc2lvbj0iMS4xIiB3aWR0aD0iMjQiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgeG1sbnM6Y2M9Imh0dHA6Ly9jcmVhdGl2ZWNvbW1vbnMub3JnL25zIyIgeG1sbnM6ZGM9Imh0dHA6Ly9wdXJsLm9yZy9kYy9lbGVtZW50cy8xLjEvIiB4bWxuczpyZGY9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMiPjxnIHRyYW5zZm9ybT0idHJhbnNsYXRlKDAgLTEwMjguNCkiPjxnPjxwYXRoIGQ9Im00IDEwMzEuNGMtMS4xMDQ2IDAtMiAwLjktMiAydjE2YzAgMS4xIDAuODk1NCAyIDIgMmgxNmMxLjEwNSAwIDItMC45IDItMnYtMTZjMC0xLjEtMC44OTUtMi0yLTJoLTE2eiIgZmlsbD0iI2QzNTQwMCIvPjxwYXRoIGQ9Im00IDJjLTEuMTA0NiAwLTIgMC44OTU0LTIgMnYxNmMwIDEuMTA1IDAuODk1NCAyIDIgMmgxNmMxLjEwNSAwIDItMC44OTUgMi0ydi0xNmMwLTEuMTA0Ni0wLjg5NS0yLTItMmgtMTZ6IiBmaWxsPSIjZTY3ZTIyIiB0cmFuc2Zvcm09InRyYW5zbGF0ZSgwIDEwMjguNCkiLz48cGF0aCBkPSJtNSAxMDM0LjR2Mi4zYzYuNDQzIDAgMTEuNjY3IDUuMiAxMS42NjcgMTEuN2gyLjMzM2MwLTcuOC02LjI2OC0xNC0xNC0xNHptMCA0LjZ2Mi40YzMuODY2IDAgNyAzLjEgNyA3aDIuMzMzYzAtNS4yLTQuMTc4LTkuNC05LjMzMy05LjR6bTIuMDQxNyA1LjNjLTEuMTI3NiAwLTIuMDQxNyAwLjktMi4wNDE3IDJzMC45MTQxIDIuMSAyLjA0MTcgMi4xYzEuMTI3NSAwIDIuMDQxNi0xIDIuMDQxNi0yLjFzLTAuOTE0MS0yLTIuMDQxNi0yeiIgZmlsbD0iI2QzNTQwMCIvPjxwYXRoIGQ9Im01IDEwMzMuNHYyLjNjNi40NDMgMCAxMS42NjcgNS4yIDExLjY2NyAxMS43aDIuMzMzYzAtNy44LTYuMjY4LTE0LTE0LTE0em0wIDQuNnYyLjRjMy44NjYgMCA3IDMuMSA3IDdoMi4zMzNjMC01LjItNC4xNzgtOS40LTkuMzMzLTkuNHptMi4wNDE3IDUuM2MtMS4xMjc2IDAtMi4wNDE3IDAuOS0yLjA0MTcgMnMwLjkxNDEgMi4xIDIuMDQxNyAyLjFjMS4xMjc1IDAgMi4wNDE2LTEgMi4wNDE2LTIuMXMtMC45MTQxLTItMi4wNDE2LTJ6IiBmaWxsPSIjZWNmMGYxIi8+PC9nPjwvZz48L3N2Zz4="),
			a.Alt("RSS Feed"),
			a.Width("16"),
			a.Height("16"),
			a.Class("feed-icon"),
		)),
		A(Attr(
			a.Href(fileName+"-rss.xml"),
			a.Title("RSS"),
			a.Target("_blank"),
		),
			Text("RSS"),
		),
		A(Attr(
			a.Href(fileName+"-atom.xml"),
			a.Title("Atom"),
			a.Target("_blank"),
		),
			Text("Atom"),
		),
	)
}

const FeedLinksCSS = `
.feed-links {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 0.5em;
    padding: 1em;
}

.feed-icon {
    width: 16px;
    height: 16px;
}
`
