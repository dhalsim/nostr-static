package components

import (
	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
)

func RenderImageHTML(image, alt, imageLink, baseFolder string) HTML {
	if image == "" {
		return Text("")
	}

	if imageLink == "" {
		return Div(Attr(a.Class("image-container")),
			Img(Attr(
				a.Src(image),
				a.Alt(alt),
			)),
		)
	}

	return Div(Attr(a.Class("image-container")),
		A(Attr(a.Href(baseFolder+imageLink+".html")),
			Img(Attr(
				a.Src(image),
				a.Alt(alt),
			)),
		),
	)
}

const ImageCSS = `
.image-container {
    margin: 20px 0;
    text-align: center;
}

.image-container img {
    max-width: 100%;
    height: auto;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

body.dark .image-container img {
    box-shadow: 0 2px 4px rgba(0,0,0,0.3);
}
`
