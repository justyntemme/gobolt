package dom

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	cssContent string
	once       sync.Once
)

// LoadCSS loads the CSS file content into memory
func LoadCSS(filePath string) {
	once.Do(func() {
		data, err := os.ReadFile(filePath) // Using os.ReadFile instead of ioutil.ReadFile
		if err != nil {
			log.Fatalf("Failed to load CSS file: %v", err)
		}
		cssContent = string(data)
		fmt.Println("CSS file loaded successfully.")
	})
}

// GetThemeCSS returns the loaded CSS content as a string

func GetThemeCSS() string {
	return cssContent
}

func main() {
	// Load the CSS file at startup
	LoadCSS("styles.css")

	// Retrieve the CSS content
	fmt.Println(getThemeCSS())
}

