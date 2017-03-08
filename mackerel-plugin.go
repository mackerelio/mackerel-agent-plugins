package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func main() {
	os.Exit(run(os.Args))
}

const (
	exitOK = iota
	exitError
)

var helpReg = regexp.MustCompile(`--?h(?:elp)?`)

//go:generate sh -c "perl tool/gen_mackerel_plugin.pl > mackerel-plugin_gen.go"
func run(args []string) int {
	var plug string
	f, err := exec.LookPath(args[0])
	if err != nil {
		log.Println(err)
		return exitError
	}
	fi, err := os.Lstat(f)
	if err != nil {
		log.Println(err)
		return exitError
	}
	base := filepath.Base(f)
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink && strings.HasPrefix(base, "mackerel-plugin-") {
		// if mackerel-plugin is symbolic linked from mackerel-plugin-memcached, run the memcached plugin
		plug = strings.TrimPrefix(base, "mackerel-plugin-")
	} else {
		if len(args) < 2 {
			printHelp()
			return exitError
		}
		plug = args[1]
		if helpReg.MatchString(plug) {
			printHelp()
			return exitOK
		}
		os.Args = append([]string{f}, args[2:]...)
	}

	err = runPlugin(plug)

	if err != nil {
		return exitError
	}
	return exitOK
}

var version, gitcommit string

func printHelp() {
	fmt.Printf(`mackerel-plugin %s (rev %s) [%s %s %s]

Usage: mackerel-plugin <plugin> [<args>]

Following plugins are available:
    %s

See `+"`mackerel-plugin <plugin> -h` "+`for more information on a specific plugin
`, version, gitcommit, runtime.GOOS, runtime.GOARCH, runtime.Version(), strings.Join(plugins, "\n    "))
}
