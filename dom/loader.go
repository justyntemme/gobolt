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

	// Log the base directory
	// d.Logger.Infof("Starting to load Markdown files from base directory: %s", baseDir)

	// Start workers to process Markdown into HTML
	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		go d.htmlWorker(taskChan, wg)
		// d.Logger.Infof("Started worker %d for HTML generation", i)
	}

	// Walk through the directory and enqueue Markdown files
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// d.Logger.Errorf("Error walking the directory: %v", err)
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			// Generate URI from file path
			uri := strings.TrimPrefix(path, baseDir)                     // Remove baseDir prefix
			uri = strings.TrimSuffix(uri, ".md")                         // Remove file extension
			uri = strings.ReplaceAll(uri, string(os.PathSeparator), "/") // Normalize slashes

			// Log the filtered URI
			// d.Logger.Infof("Filtered URI: %s", uri)

			// Read the Markdown content from the file
			content, err := os.ReadFile(path)
			if err != nil {
				// d.Logger.Errorf("Failed to read file %s: %v", path, err)
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			// Log the content being added
			// d.Logger.Infof("Adding page to DOM with URI: %s", uri)

			// Add the page to the DOM with its Markdown content
			d.Pages[uri] = &Page{Markdown: string(content)}

			// Enqueue the URI for HTML generation
			wg.Add(1)
			taskChan <- uri
		}
		return nil
	})

	// Close the task channel and wait for workers to finish
	close(taskChan)
	wg.Wait()

	// Log the completion of the loading process
	// d.Logger.Infof("Finished loading all Markdown files from %s", baseDir)

	return err
}
