package gfile

import (
	"os"
	"path/filepath"
	"strings"
)

// TODO: verify with test cases

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

func GetAllMatches(patterns []string, optionalCommonRoot ...string) ([]string, error) {
	var commonRoot string
	if len(optionalCommonRoot) > 0 {
		commonRoot = optionalCommonRoot[0]
	} else {
		commonRoot = findCommonRoot(patterns)
	}

	// Pre-calculate pattern info
	patternInfos := make([]patternInfo, 0, len(patterns))
	for _, pattern := range patterns {
		isNegation := false
		if pattern[0] == '!' {
			isNegation = true
			pattern = pattern[1:]
		}
		patternSegments := strings.Split(pattern, "/")
		patternInfos = append(patternInfos, patternInfo{segments: patternSegments, isNegation: isNegation, minLength: len(patternSegments)})
	}

	matches := make(map[string]bool)
	if err := matchGlob(commonRoot, patternInfos, matches); err != nil {
		return nil, err
	}

	var allMatches []string
	for match := range matches {
		allMatches = append(allMatches, match)
	}

	return allMatches, nil
}

type patternInfo struct {
	segments   []string
	isNegation bool
	minLength  int
}

func findCommonRoot(patterns []string) string {
	if len(patterns) == 0 {
		return "./"
	}

	minPath := patterns[0]
	maxPath := patterns[0]
	for _, path := range patterns[1:] {
		if path < minPath {
			minPath = path
		}
		if path > maxPath {
			maxPath = path
		}
	}

	commonRoot := []string{}
	for i, minSeg := range strings.Split(minPath, "/") {
		maxSeg := strings.Split(maxPath, "/")[i]

		if minSeg == maxSeg {
			commonRoot = append(commonRoot, minSeg)
		} else {
			break
		}
	}

	return strings.Join(commonRoot, "/")
}

func matchGlob(root string, patterns []patternInfo, matches map[string]bool) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // ignore permissions error
		}

		relativePath, _ := filepath.Rel(root, path)
		if relativePath == "." {
			return nil // Skip the root itself
		}

		// Normalize path to use forward slashes on all platforms
		relativePath = filepath.ToSlash(relativePath)
		pathSegments := strings.Split(relativePath, "/")

		for _, patternInfo := range patterns {
			if len(pathSegments) < patternInfo.minLength {
				continue
			}

			matchesPattern := true
			j := 0

			for i, segment := range patternInfo.segments {
				if segment == "**" {
					j = len(pathSegments) - len(patternInfo.segments) + i + 1
					continue
				}

				if j >= len(pathSegments) {
					matchesPattern = false
					break
				}

				if match, _ := filepath.Match(segment, pathSegments[j]); !match {
					matchesPattern = false
					break
				}

				j++
			}

			if matchesPattern {
				if patternInfo.isNegation {
					delete(matches, path)
				} else {
					matches[path] = true
				}
			}
		}

		return nil
	})
}
