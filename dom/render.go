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
		io.WriteString(w, `<div class="paragraph"><p>`)
		io.Writer.Write(w, p.Content)
	} else {
		io.WriteString(w, "</p></div>")
	}
}

func myRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	// renderHookLogger := logrus.New()
	CSS := getThemeCSS()
	io.WriteString(w, CSS)
	// Paragraph Logic
	if para, ok := node.(*ast.Paragraph); ok {
		renderParagraph(w, para, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func newCustomizedRender() *html.Renderer {
	Header := `<title>Hello world</title>`
	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: myRenderHook,
		CSS:            "./Content/styles.css",
		Head:           []byte(Header),
	}
	return html.NewRenderer(opts)
}
