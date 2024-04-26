package gprocess

import (
	"bytes"
	"os/exec"
)

// CommandExecutor is an interface to abstract command execution
//
//go:generate mockgen -source=executor.go -destination=executor_mock_test.go -package=gcprocess_test
type CommandExecutor interface {
	ExecuteCommand(stdout, stderr *bytes.Buffer, command string, args ...string) error
}

// RealCommandExecutor is a real implementation of the CommandExecutor interface.
type RealCommandExecutor struct{}

func (e RealCommandExecutor) ExecuteCommand(stdout, stderr *bytes.Buffer, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}
