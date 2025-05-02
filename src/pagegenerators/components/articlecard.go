package components

import (
	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
)

func RenderArticleCard(title, summary, image, naddr string, tags []string, baseFolder string) HTML {
	return Div(Attr(a.Class("article-card")),
		RenderImageHTML(image, title, naddr, baseFolder),
		H2_(
			A(Attr(a.Href(baseFolder+naddr+".html")),
				Text(title),
			),
		),
		RenderSummaryHTML(summary),
		RenderTagsHTML(tags, baseFolder),
	)
}

const ArticleCardCSS = `
.article-card {
    border: 1px solid #ddd;
    border-radius: 8px;
    padding: 20px;
    margin-bottom: 20px;
}

.article-card h2 {
    margin-top: 0;
    text-align: center;
}

/* Theme-specific component styles */
body.light .article-card {
    background-color: #f8f9fa;
    border: 1px solid #dee2e6;
}

body.dark .article-card {
    background-color: #2d2d2d;
    border: 1px solid #404040;
}

body.light .article-card .summary {
    color: #666666;
}

body.dark .article-card .summary {
    color: #e0e0e0;
}
`
