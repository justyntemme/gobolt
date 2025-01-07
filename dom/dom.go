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

// DOM represents a collection of pages.
type DOM struct {
	Pages  map[string]*Page
	mu     sync.Mutex
	Logger logrus.Logger // Assuming Logger is defined elsewhere
}

// NewDOM creates a new DOM instance.
func NewDOM(logger *logrus.Logger) *DOM {
	return &DOM{
		Pages:  make(map[string]*Page),
		Logger: *logger,
	}
}

// htmlWorker processes the Markdown content for a specific URI and generates the HTML.
func (d *DOM) htmlWorker(taskChan <-chan string, wg *sync.WaitGroup) {
	for uri := range taskChan {
		// Log the URI being processed
		d.Logger.Infof("Processing URI: %s", uri)

		// Lock to safely access the Pages map
		d.mu.Lock()
		page, exists := d.Pages[uri]
		d.mu.Unlock()

		if !exists {
			d.Logger.Warnf("Page with URI %s does not exist in DOM", uri)
			wg.Done()
			continue
		}

		// Log before starting the conversion
		d.Logger.Infof("Converting Markdown to HTML for URI: %s", uri)

		// Convert Markdown to HTML
		html := markdown.ToHTML([]byte(page.Markdown), nil, nil)

		// Log after conversion is completed
		d.Logger.Infof("Conversion complete for URI: %s", uri)

		// Lock to safely write the HTML back into the DOM
		d.mu.Lock()
		page.HTML = string(html)
		d.mu.Unlock()

		// Log when HTML is successfully written to the Page
		d.Logger.Infof("HTML written to DOM for URI: %s", uri)

		// Mark this task as complete
		wg.Done()
	}
}
