package components

import (
	"strings"

	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
)

func RenderLogo(logo, baseFolder string) HTML {
	if logo == "" {
		return Text("")
	}

	baseFolder = strings.Trim(baseFolder, "/")
	if baseFolder != "" {
		baseFolder = baseFolder + "/"
	}

	return Div(Attr(a.Class("logo")),
		A(Attr(a.Href(baseFolder+"index.html")),
			Img(Attr(
				a.Src(baseFolder+logo),
				a.Alt("Logo"),
			)),
		),
	)
}

const LogoCSS = `
.logo {
    text-align: center;
}

.logo img {
    max-height: 100px;
    width: auto;
}

.logo-container {
    margin-top: 20px;
    margin-right: 20px;
}

.logo-container img {
    max-height: 50px;
    width: auto;
}
`
