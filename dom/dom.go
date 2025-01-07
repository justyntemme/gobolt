package dom

import (
	"sync"
	"github.com/gomarkdown/markdown"
)

type Page struct {
	Markdown string
	HTML     string
}

type DOM struct {
	Pages map[string]*Page
	mu    sync.Mutex
}

func NewDOM() *DOM {
	return &DOM{
		Pages: make(map[string]*Page),
	}
}

func (d *DOM) htmlWorker(taskChan <-chan string, wg *sync.WaitGroup) {
	for uri := range taskChan {
		d.mu.Lock()
		page, exists := d.Pages[uri]
		d.mu.Unlock()

		if !exists {
			wg.Done()
			continue
		}

		// Convert Markdown to HTML
		html := markdown.ToHTML([]byte(page.Markdown), nil, nil)

		d.mu.Lock()
		page.HTML = string(html)
		d.mu.Unlock()

		wg.Done()
	}
}
