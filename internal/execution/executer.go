package execution

type ExecutionResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

type Executer interface {
	Execute() (*ExecutionResult, error)
	Cleanup()
}
