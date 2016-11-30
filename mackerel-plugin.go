package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	os.Exit(run(os.Args))
}

const (
	exitOK = iota
	exitError
)

var helpReg = regexp.MustCompile(`--?h(?:elp)?`)

func run(args []string) int {
	if len(args) < 2 {
		printHelp()
		return exitError
	}
	plug := args[1]
	if helpReg.MatchString(plug) {
		printHelp()
		return exitOK
	}
	osargs := []string{args[0]}
	osargs = append(osargs, args[2:]...)
	os.Args = osargs
	err := runPlugin(plug)

	if err != nil {
		return exitError
	}
	return exitOK
}

func printHelp() {
	fmt.Println("please specify the plugin by argument")
}
