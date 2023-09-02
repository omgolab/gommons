package gcfile

import (
	"bufio"
	"os"
	"strings"
)

// Contains: scans a file line by line
// and check if it contains the given list of strings
// sLine and eLine are the start and end lines to scan
// it starts from 1 where 0 means all lines
func ContainsText(path string, sLine, eLine uint64, args ...string) ([]bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	arr := make([]bool, len(args))
	cLine := uint64(1)
	if eLine > 0 && sLine > eLine {
		eLine = sLine
	}

	for scanner.Scan() {
		t := scanner.Text()
		if cLine >= sLine {
			for i, arg := range args {
				if !arr[i] {
					arr[i] = strings.Contains(t, arg)
				}
			}
		}
		if cLine >= eLine && eLine > 0 {
			break
		}
		cLine++
	}

	if err := scanner.Err(); err != nil {
		return arr, err
	}

	return arr, nil
}

// ContainsAllText: scans a file line by line for the given list of strings
// and returns true if all the strings are found
// and false if any of the strings is not found
func ContainsAllTexts(path string, sLine, eLine uint64, args ...string) (bool, error) {
	b, err := ContainsText(path, sLine, eLine, args...)
	if err != nil {
		return false, err
	}
	for _, v := range b {
		if !v {
			return false, nil
		}
	}
	return true, nil
}
