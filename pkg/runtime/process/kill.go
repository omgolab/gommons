package gprocess

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// KillByNamePattern kills processes by name regex pattern.
// It returns an error if any occurred.
func KillByNamePattern(pattern string, killSelf bool) error {
	var cmd *exec.Cmd
	var output []byte
	var err error

	if pattern == "" {
		return errors.New("process kill failed. name is empty")
	}

	var pIDs []string
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s*", pattern))
		output, err = cmd.Output()
		if err != nil {
			return fmt.Errorf("error executing tasklist: %w", err)
		}

		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) > 1 && fields[1] != "PID" {
				pIDs = append(pIDs, fields[1])
			}
		}
	} else {
		cmd = exec.Command("pgrep", "-f", pattern)
		output, err = cmd.Output()
		if err != nil {
			return fmt.Errorf("error executing pgrep: %w", err)
		}

		pIDs = strings.Split(strings.TrimSpace(string(output)), "\n")
	}

	return KillByPID(killSelf, pIDs...)
}

func KillByPID(killSelf bool, pIDs ...string) error {
	currentPID := strconv.Itoa(os.Getpid())
	foundSelf := false

	for _, PID := range pIDs {
		if PID == currentPID {
			foundSelf = true
			continue
		}

		err := kill(PID)
		if err != nil {
			return err
		}
	}

	// if we're killing itself then we need to return the response first then kill
	if foundSelf {
		go func() {
			time.Sleep(time.Second * 2)
			_ = kill(currentPID)
		}()
		return nil
	}

	return nil
}

func kill(PID string) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("taskkill", "/F", "/PID", PID)
	} else {
		cmd = exec.Command("kill", "-9", PID)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error killing process with PID %s: %w", PID, err)
	}

	return nil
}
