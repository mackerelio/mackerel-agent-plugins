package mpunicorn

import "github.com/mattn/go-pipeline"

// PipedCommands interface
type PipedCommands interface {
	Output(...[]string) ([]byte, error)
}

// RealPipedCommands struct
type RealPipedCommands struct{}

var pipedCommands PipedCommands

// Output for RealPipedCommands
func (r RealPipedCommands) Output(commands ...[]string) ([]byte, error) {
	out, err := pipeline.Output(commands...)
	return out, err
}
