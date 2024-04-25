package gcfile_test

import (
	"os"
	"path/filepath"
	"testing"

	file_utils "github.com/omgolab/go-commons/pkg/file"
)

func TestGetDirSize(t *testing.T) {
	// Create a temporary directory for testing
	dir := createTempDir()
	defer os.RemoveAll(dir)

	// Create files with different sizes in the temporary directory
	filePaths := []string{
		createTestFile(dir, "file1.txt", 100),
		createTestFile(dir, "file2.txt", 200),
		createTestFile(dir, "file3.txt", 300),
	}

	// Get the expected size of the directory
	expectedSize := int64(0)
	for _, filePath := range filePaths {
		fileInfo, _ := os.Stat(filePath)
		expectedSize += fileInfo.Size()
	}

	// Call the function to get the actual size of the directory
	actualSize, err := file_utils.GetDirSize(dir)
	if err != nil {
		t.Errorf("GetDirSize returned an error: %v", err)
	}

	// Check if the actual size matches the expected size
	if actualSize != expectedSize {
		t.Errorf("Actual size (%d) does not match expected size (%d)", actualSize, expectedSize)
	}
}

// Helper functions for creating temporary directory and files

func createTempDir() string {
	dir, err := os.MkdirTemp("", "test-dir")
	if err != nil {
		panic(err)
	}
	return dir
}

func createTestFile(dir, filename string, size int64) string {
	filePath := filepath.Join(dir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_ = file.Truncate(size)

	return filePath
}
