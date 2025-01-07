package dom

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func (d *DOM) LoadMarkdown(baseDir string) error {
	taskChan := make(chan string, 10) // Buffered channel
	wg := &sync.WaitGroup{}

	// Start workers to process Markdown into HTML
	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		go d.htmlWorker(taskChan, wg)
	}

	// Walk through the directory and enqueue Markdown files
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			// Generate URI from file path
			uri := strings.TrimPrefix(path, baseDir)
			uri = strings.TrimSuffix(uri, ".md")
			uri = strings.ReplaceAll(uri, string(os.PathSeparator), "/")

			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			// Add the page to the DOM with its Markdown content
			d.mu.Lock()
			d.Pages[uri] = &Page{Markdown: string(content)}
			d.mu.Unlock()

			// Enqueue the URI for HTML generation
			wg.Add(1)
			taskChan <- uri
		}
		return nil
	})

	// Close the task channel and wait for workers to finish
	close(taskChan)
	wg.Wait()

	return err
}
