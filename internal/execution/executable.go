package execution

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

type Executable struct {
	directory  string
	executeCmd string
}

func (s *Executable) Execute() (ExecutionResult, error) {
	cmd := exec.Command(strings.Split(s.executeCmd, " ")[0], strings.Split(s.executeCmd, " ")[1:]...)
	cmd.Dir = s.directory
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return ExecutionResult{
		Stdout:   string(stdoutStderr),
		Stderr:   string(stdoutStderr),
		ExitCode: cmd.ProcessState.ExitCode(),
	}, nil
}

func (s *Executable) Cleanup() {
	os.RemoveAll(s.directory)
}
