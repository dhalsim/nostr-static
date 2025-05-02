package components

import (
	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
)

func RenderSummaryHTML(summary string) HTML {
	if summary == "" {
		return Text("")
	}
	return P(Attr(a.Class("summary")), Text(summary))
}
