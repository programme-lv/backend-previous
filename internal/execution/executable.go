package execution

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

type Executable struct {
	directory  string
	executeCmd string
}

func (s *Executable) Execute() (*ExecutionResult, error) {
	cmd := exec.Command(strings.Split(s.executeCmd, " ")[0], strings.Split(s.executeCmd, " ")[1:]...)
	cmd.Dir = s.directory

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	slurpOut, _ := io.ReadAll(stdout)
	slurpErr, _ := io.ReadAll(stderr)

	_ = cmd.Wait()

	return &ExecutionResult{
		Stdout:   string(slurpOut),
		Stderr:   string(slurpErr),
		ExitCode: cmd.ProcessState.ExitCode(),
	}, nil
}

func (s *Executable) Cleanup() {
	os.RemoveAll(s.directory)
}
