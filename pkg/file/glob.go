package gcfile

import (
	"os"
	"path/filepath"
	"strings"
)

// DeleteGlobPatternedPaths: Delete the files/folders as per the patterns mentioned
func DeleteGlobPatternedPaths(patterns string) {
	ps := strings.Fields(patterns)
	for _, pattern := range ps {
		files, _ := filepath.Glob(pattern)
		for _, file := range files {
			os.RemoveAll(file)
		}
	}
}
