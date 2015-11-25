package main

import "github.com/mattn/go-pipeline"

type PipedCommands interface {
	Output(...[]string) ([]byte, error)
}

type RealPipedCommands struct{}

var pipedCommands PipedCommands

func (r RealPipedCommands) Output(commands ...[]string) ([]byte, error) {
	out, err := pipeline.Output(commands...)
	return out, err
}
