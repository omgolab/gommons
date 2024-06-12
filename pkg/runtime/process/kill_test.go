package gprocess

import "testing"

func TestKillByNamePattern(t *testing.T) {

	err := KillByNamePattern("", true)
	if err != nil {
		t.Errorf("ProcessFilter.ListProcesses() error = %v", err)
	}
}
