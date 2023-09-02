package gcfile_test

import (
	"os"
	"path/filepath"
	"testing"

	fu "github.com/omar391/go-commons/pkg/file"
)

func TestDeleteGlobePatterns(t *testing.T) {
	// Create temporary files for testing
	createTempFiles()

	// Test case 1: Delete files with pattern "*.txt"
	patterns := "*.txt" // Replace with your specific patterns
	fu.DeleteGlobPatternedPaths(patterns)

	// Validate that the files matching the pattern are deleted
	files, _ := filepath.Glob(patterns)
	for _, file := range files {
		t.Errorf("Expected file %s to be deleted, but it still exists", file)
	}

	// Test case 2: Delete files with pattern "*.doc"
	patterns = "*.doc" // Replace with your specific patterns
	fu.DeleteGlobPatternedPaths(patterns)

	// Validate that the files matching the pattern are deleted
	files, _ = filepath.Glob(patterns)
	for _, file := range files {
		t.Errorf("Expected file %s to be deleted, but it still exists", file)
	}

	// Cleanup: Delete temporary files
	deleteTempFiles()

}

// Helper function to create temporary files for testing
func createTempFiles() {
	files := []string{"file1.txt", "file2.txt", "file3.doc", "file4.docx", "file5.pdf"}

	for _, file := range files {
		content := []byte("This is a temporary file.")
		_ = os.WriteFile(file, content, 0644)
	}

}

// Helper function to delete temporary files
func deleteTempFiles() {
	files, _ := filepath.Glob("*.txt")
	for _, file := range files {
		_ = os.Remove(file)
	}

	files, _ = filepath.Glob("*.doc")
	for _, file := range files {
		_ = os.Remove(file)
	}

	files, _ = filepath.Glob("*.docx")
	for _, file := range files {
		_ = os.Remove(file)
	}

	files, _ = filepath.Glob("*.pdf")
	for _, file := range files {
		_ = os.Remove(file)
	}
}
