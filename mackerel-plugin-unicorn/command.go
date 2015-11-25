package main

import "os/exec"

type Command interface {
	Output(string, ...string) ([]byte, error)
}

type RealCommand struct{}

var command Command

func (r RealCommand) Output(command string, args ...string) ([]byte, error) {
	out, err := exec.Command(command, args...).Output()
	return out, err
}
