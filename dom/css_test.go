package dom

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCSS(t *testing.T) {
	// Test cases struct to organize our tests
	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "Valid CSS file",
			filePath: "testdata/valid.css",
			wantErr:  false,
		},
		{
			name:     "Non-existent file",
			filePath: "testdata/nonexistent.css",
			wantErr:  true,
		},
		{
			name:     "Empty file path",
			filePath: "",
			wantErr:  true,
		},
	}

	// Create test directory and file
	testDir := "testdata"
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Create a valid CSS file for testing
	validCSS := "body { background-color: #fff; }"
	err = os.WriteFile(filepath.Join(testDir, "valid.css"), []byte(validCSS), 0644)
	if err != nil {
		t.Fatalf("Failed to create test CSS file: %v", err)
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset cssContent before each test
			cssContent = ""

			// Execute the function
			err := LoadCSS(tt.filePath)

			// Check if error matches expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadCSS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// For successful cases, verify content
			if !tt.wantErr {
				if got := GetThemeCSS(); got != validCSS {
					t.Errorf("GetThemeCSS() = %v, want %v", got, validCSS)
				}
			}
		})
	}
}
