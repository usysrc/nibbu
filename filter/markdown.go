package filter

import (
	"html/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
)

func MarkdownFilter(input template.HTML) template.HTML {
	p := parser.New()
	doc := p.Parse([]byte(input))
	htmlRenderer := html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags})
	content := markdown.Render(doc, htmlRenderer)
	policy := bluemonday.UGCPolicy() // User-Generated Content policy.
	sanitized := policy.SanitizeBytes(content)
	return template.HTML(sanitized)
}
