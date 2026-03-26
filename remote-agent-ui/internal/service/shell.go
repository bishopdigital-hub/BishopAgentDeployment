package service

import (
	"io"
	"os/exec"
)

type ShellService struct{}

func NewShellService() *ShellService {
	return &ShellService{}
}

func (s *ShellService) RunCommand(command string, args []string, stdout, stderr io.Writer) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}
