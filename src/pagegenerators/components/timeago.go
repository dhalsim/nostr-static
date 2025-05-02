package components

import (
	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
)

func RenderTimeAgoScript() HTML {
	return Script(Attr(a.Src("/output/static/js/time-ago.js")), JavaScript(""))
}
