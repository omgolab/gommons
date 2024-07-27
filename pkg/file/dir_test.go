package gfile_test

import (
	"os"
	"path/filepath"
	"testing"

	gfile "github.com/omgolab/go-commons/pkg/file"
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
	actualSize, err := gfile.GetDirSize(dir)
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

func TestCopyFile(t *testing.T) {
	srcFile, err := os.CreateTemp("", "srcFile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(srcFile.Name())

	content := []byte("Hello, World!")
	if _, err := srcFile.Write(content); err != nil {
		t.Fatal(err)
	}
	srcFile.Close()

	dstFile, err := os.CreateTemp("", "dstFile")
	if err != nil {
		t.Fatal(err)
	}
	dstFile.Close()
	defer os.Remove(dstFile.Name())

	if err := gfile.CopyFile(srcFile.Name(), dstFile.Name(), 0644); err != nil {
		t.Fatal(err)
	}

	dstContent, err := os.ReadFile(dstFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if string(dstContent) != string(content) {
		t.Errorf("expected %s but got %s", content, dstContent)
	}
}

func TestCopyDir(t *testing.T) {
	srcDir, err := os.MkdirTemp("", "srcDir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(srcDir)

	subDir := filepath.Join(srcDir, "subDir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	srcFile := filepath.Join(subDir, "file.txt")
	if err := os.WriteFile(srcFile, []byte("Hello, World!"), 0644); err != nil {
		t.Fatal(err)
	}

	dstDir, err := os.MkdirTemp("", "dstDir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dstDir)

	if err := gfile.CopyDir(srcDir, dstDir); err != nil {
		t.Fatal(err)
	}

	dstFile := filepath.Join(dstDir, "subDir", "file.txt")
	dstContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(dstContent) != "Hello, World!" {
		t.Errorf("expected %s but got %s", "Hello, World!", dstContent)
	}
}
