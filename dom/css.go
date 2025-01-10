package dom

import (
	"fmt"
	"os"
)

var cssContent string

// LoadCSS loads the CSS file content into memory
func LoadCSS(filePath string) error {
	data, err := os.ReadFile(filePath) // Using os.ReadFile instead of ioutil.ReadFile
	if err != nil {
		return err
	}
	cssContent = string(data)
	fmt.Println("CSS file loaded successfully.")

	return nil
}

func GetThemeCSS() string {
	return cssContent
}
