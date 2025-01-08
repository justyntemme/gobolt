package dom

import (
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

type Paragraph struct {
	ast.Leaf
	data string
}

// an actual rendering of Paragraph is more complicated
func renderParagraph(w io.Writer, p *ast.Paragraph, entering bool) {
	if entering {
		io.WriteString(w, `<div class="paragraph"<p>`)
		io.WriteString(w, string(p.Content))
	} else {
		io.WriteString(w, "</p></div>")
	}
}

func myRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	// Paragraph Logic
	if para, ok := node.(*ast.Paragraph); ok {
		renderParagraph(w, para, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func newCustomizedRender() *html.Renderer {
	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: myRenderHook,
	}
	return html.NewRenderer(opts)
}
