package dom

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestDOM_LoadMarkdown(t *testing.T) {
	// Create a temporary directory for test files
	tempValidDir, err := os.MkdirTemp("", "markdown-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempValidDir)
	tempInvalidDir, err := os.MkdirTemp("", "markdown-test-fail")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempInvalidDir)

	// Create test markdown files
	validMarkdown := "# Test Header\nThis is valid markdown."
	invalidMarkdown := string([]byte{0xFF, 0xFE, 0xFD}) // Invalid UTF-8

	validPath := filepath.Join(tempValidDir, "valid.md")
	invalidPath := filepath.Join(tempInvalidDir, "invalid.md")

	if err := os.WriteFile(validPath, []byte(validMarkdown), 0644); err != nil {
		t.Fatalf("Failed to create valid test file: %v", err)
	}
	if err := os.WriteFile(invalidPath, []byte(invalidMarkdown), 0644); err != nil {
		t.Fatalf("Failed to create invalid test file: %v", err)
	}

	type fields struct {
		Pages           map[string]*Page
		pagesUpdateChan chan pageUpdate
	}
	type args struct {
		baseDir string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func() error
		cleanup func()
	}{
		{
			name: "Valid directory with markdown files",
			fields: fields{
				Pages:           make(map[string]*Page),
				pagesUpdateChan: make(chan pageUpdate, 10),
			},
			args: args{
				baseDir: tempValidDir,
			},
			wantErr: false,
		},
		{
			name: "Directory with invalid markdown files",
			fields: fields{
				Pages:           make(map[string]*Page),
				pagesUpdateChan: make(chan pageUpdate, 10),
			},
			args: args{
				baseDir: tempInvalidDir,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Fatalf("Test setup failed: %v", err)
				}
			}

			d := &DOM{
				Pages:           tt.fields.Pages,
				pagesUpdateChan: tt.fields.pagesUpdateChan,
			}

			err := d.LoadMarkdown(tt.args.baseDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("DOM.LoadMarkdown() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.cleanup != nil {
				tt.cleanup()
			}

			//Additional assertions for successful cases // Also should add content check for invalid cases
			if !tt.wantErr && err == nil {
				// Check if valid files were properly loaded
				if tt.args.baseDir == tempValidDir {
					expectedPaths := []string{"/valid"}
					for _, path := range expectedPaths {
						if _, exists := d.Pages[path]; !exists {
							t.Errorf("Expected page at path %s to exist in DOM", path)
						}
					}
				}
				// TODO add tests for specific rendering objects {paragraph, header, CSS}
				if tt.args.baseDir == tempInvalidDir {
					html := tt.fields.Pages["/invalid"].HTML
					hash := md5.Sum([]byte(html))
					hashString := hex.EncodeToString(hash[:])
					if hashString != "d91a8480cda772fd34e02fb61e1a226d" {
						t.Errorf("Expected Specific hash value for invalid string")
					}

				}
			}
		})
	}
}
