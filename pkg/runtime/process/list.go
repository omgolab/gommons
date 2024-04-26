package gprocess

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// ProcessFilter represents a map-based configuration for filtering the list of processes.
type ProcessFilter struct {
	Attrs           map[string]string
	CommandExecutor CommandExecutor
}

func NewLiveProcessFilter(attrs map[string]string) ProcessFilter {
	return ProcessFilter{
		Attrs:           attrs,
		CommandExecutor: RealCommandExecutor{},
	}
}

// ListProcesses fetches and filters processes based on the underlying OS and the provided filter configuration.
// It returns a slice of maps, each representing a process, and an error if any occurred.
func (pf ProcessFilter) ListProcesses() ([]map[string]string, error) {
	var command string
	var args []string
	var stdout, stderr bytes.Buffer

	if runtime.GOOS == "windows" {
		command = "cmd"
		args = []string{"/C", "tasklist"}
	} else {
		command = "sh"
		args = []string{"-c", "ps aux"}
	}

	// execute the command
	err := pf.CommandExecutor.ExecuteCommand(&stdout, &stderr, command, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing command: %w, stderr: %v", err, stderr)
	}

	rawOutput := stdout.Bytes()
	lines := bytes.Split(rawOutput, []byte{'\n'})
	if len(lines) < 2 {
		return nil, errors.New("Invalid output format")
	}

	rawHeader := bytes.Fields(lines[0])
	header := sanitizeHeaders(rawHeader)

	processes := make([]map[string]string, 0, len(lines)-1)

	for _, line := range lines[1:] {
		if len(line) == 0 {
			continue
		}
		values := bytes.Fields(line)
		if len(values) < len(header) {
			continue
		}

		processObj := make(map[string]string)
		match := true

		for index, field := range header {
			valStr := string(values[index])
			if v, ok := pf.Attrs[field]; ok && !strings.Contains(valStr, v) {
				match = false
				break
			}
			processObj[field] = valStr
		}

		if match {
			processes = append(processes, processObj)
		}
	}
	return processes, nil
}

// sanitizeHeaders removes special characters from header fields
func sanitizeHeaders(rawHeader [][]byte) []string {
	header := make([]string, len(rawHeader))
	for i, h := range rawHeader {
		header[i] = strings.Trim(string(h), "%")
	}
	return header
}
