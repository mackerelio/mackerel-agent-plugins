package mpmailq

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

type mailq struct {
	command string
	args    []string
	line    int
	pattern string
}

var mailqFormats = map[string]mailq{
	"postfix": {
		command: "postqueue",
		args:    []string{"-p"},
		line:    -1,
		pattern: `-- \d+ Kbytes in (\d+) Requests\.`,
	},
	"qmail": {
		command: "qmail-qstat",
		pattern: `messages in queue: (\d+)`,
	},
	"exim": {
		command: "exim",
		args:    []string{"-bpc"},
		pattern: `(\d+)`,
	},
}

type plugin struct {
	path                   string
	mailq                  mailq
	keyPrefix, labelPrefix string
}

func (format *mailq) parse(rd io.Reader) (count uint64, err error) {
	l, err := getNthLine(rd, format.line)
	if err != nil {
		return
	}

	m := regexp.MustCompile(format.pattern).FindStringSubmatch(l)
	if m != nil {
		count, err = strconv.ParseUint(m[1], 10, 64)
	}

	return
}

func (p *plugin) fetchMailqCount() (count uint64, err error) {

	var path string
	if p.path != "" {
		path = p.path
	} else {
		path, err = exec.LookPath(p.mailq.command)
		if err != nil {
			return
		}
	}

	cmd := exec.Cmd{
		Path: path,
		Args: append([]string{p.mailq.command}, p.mailq.args...),
	}

	if p.path != "" {
		cmd.Path = p.path
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	count, err = p.mailq.parse(stdout)
	if err != nil {
		cmd.Wait()
		return
	}

	err = cmd.Wait()
	return
}

func (p *plugin) FetchMetrics() (map[string]interface{}, error) {
	count, err := p.fetchMailqCount()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"count": count}, nil
}

func (p *plugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		p.keyPrefix: {
			Label: p.labelPrefix + " Count",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "count", Label: "count", Type: "uint64"},
			},
		},
	}
}

func getNthLine(rd io.Reader, nth int) (string, error) {
	var result string

	brd := bufio.NewReader(rd)

	if nth >= 0 { // get nth-first line (0 is the first line)
		for i := 0; ; i++ {
			s, err := brd.ReadString('\n')
			if err == nil || err == io.EOF {
				if i == nth {
					result = strings.TrimRight(s, "\n")
				}
				if err == io.EOF {
					break
				}
			} else {
				return "", err
			}
		}
	} else { // get nth-last line (-1 is the last line, -2 is the second last line)
		buffer := make([]string, -nth)
		i := 0
		for {
			s, err := brd.ReadString('\n')
			if err == nil || err == io.EOF {
				if s != "" {
					buffer[i%-nth] = strings.TrimRight(s, "\n")
					i++
				}

				if err == io.EOF {
					result = buffer[i%-nth]
					break
				}
			} else {
				return "", err
			}
		}
	}

	return result, nil
}

// Do the plugin
func Do() {
	var mtas []string
	for k := range mailqFormats {
		mtas = append(mtas, k)
	}

	mta := flag.String("mta", "", fmt.Sprintf("type of MTA (one of %v)", mtas))
	flag.StringVar(mta, "M", "", "shorthand for -mta")
	command := flag.String("command", "", "path to queue-printing command (guessed by -M flag if not given)")
	flag.StringVar(command, "c", "", "shorthand for -command")
	tempfile := flag.String("tempfile", "", "path to tempfile")
	keyPrefix := flag.String("metric-key-prefix", "mailq", "prefix to metric key")
	labelPrefix := flag.String("metric-label-prefix", "Mailq", "prefix to metric label")

	flag.Parse()

	if format, ok := mailqFormats[*mta]; *mta == "" || !ok {
		fmt.Fprintf(os.Stderr, "Unknown MTA: %s\n", *mta)
		flag.PrintDefaults()
		os.Exit(1)
	} else {
		plugin := &plugin{
			path:        *command,
			mailq:       format,
			keyPrefix:   *keyPrefix,
			labelPrefix: *labelPrefix,
		}
		helper := mp.NewMackerelPlugin(plugin)
		helper.Tempfile = *tempfile
		helper.Run()
	}
}
