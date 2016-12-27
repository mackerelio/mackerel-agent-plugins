package mpunicorn

import "os/exec"

// Command interface
type Command interface {
	Output(string, ...string) ([]byte, error)
}

// RealCommand struct
type RealCommand struct{}

var command Command

// Output for RealCommand
func (r RealCommand) Output(command string, args ...string) ([]byte, error) {
	out, err := exec.Command(command, args...).Output()
	return out, err
}
