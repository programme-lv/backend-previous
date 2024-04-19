package objects

type ProgrammingLanguage struct {
	ID                string
	Name              string
	CodeFilename      string
	CompileCommand    *string
	ExecuteCommand    string
	EnvVersionCommand *string
	HelloWorldCode    *string
	MonacoID          *string
	Enabled           bool
}
