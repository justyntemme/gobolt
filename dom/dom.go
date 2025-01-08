package dom

import (
	"sync"

	"github.com/gomarkdown/markdown"
	"github.com/sirupsen/logrus"
)

// Page represents a single page with Markdown content and generated HTML.
type Page struct {
	Markdown string
	HTML     string
}

type DOM struct {
	Pages           map[string]*Page
	Logger          logrus.Logger // Assuming Logger is defined elsewhere
	pagesUpdateChan chan pageUpdate
}

type pageUpdate struct {
	uri     string
	content string
}

// NewDOM creates a new DOM instance.
func NewDOM(logger *logrus.Logger) *DOM {
	return &DOM{
		Pages:           make(map[string]*Page),
		Logger:          *logger,
		pagesUpdateChan: make(chan pageUpdate, 10),
	}
}

// htmlWorker processes the Markdown content for a specific URI and generates the HTML.
func (d *DOM) htmlWorker(taskChan <-chan string, wg *sync.WaitGroup) {
	for uri := range taskChan {
		// Log the URI being processed
		d.Logger.Infof("Processing URI: %s", uri)

		// Get the page from Pages map
		page, exists := d.Pages[uri]

		if !exists {
			d.Logger.Warnf("Page with URI %s does not exist in DOM", uri)
			wg.Done()
			continue
		}

		// Log before starting the conversion
		d.Logger.Infof("Converting Markdown to HTML for URI: %s", uri)

		// Convert Markdown to HTML
		// Using default renderer
		// html := markdown.ToHTML([]byte(page.Markdown), nil, nil)
		renderer := newCustomizedRender()
		d.Logger.Debug(renderer.Opts.CSS)
		html := markdown.ToHTML([]byte(page.Markdown), nil, renderer)

		// Log after conversion is completed
		d.Logger.Infof("Conversion complete for URI: %s", uri)

		// Update the HTML field of the page
		page.HTML = string(html)

		// Log when HTML is successfully written to the Page
		d.Logger.Infof("HTML written to DOM for URI: %s", uri)

		// Mark this task as complete
		wg.Done()
	}
}
