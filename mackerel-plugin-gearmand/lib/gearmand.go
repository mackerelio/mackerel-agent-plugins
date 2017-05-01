package mpgearmand

import (
	"bufio"
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

const nullPrefix = "-"

var graphdef = map[string]mp.Graphs{
	"gearmand.queue.#": {
		Label: "Gearmand Queue",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "available", Label: "Available", Diff: false, Stacked: false},
			{Name: "running", Label: "Running", Diff: false, Stacked: false},
			{Name: "total", Label: "Total", Diff: false, Stacked: false},
		},
	},
}

type gearmandFunction struct {
	function  string
	available uint32
	running   uint32
	total     uint32
}

func (f *gearmandFunction) key(key string) string {
	return "gearmand.queue." + f.name() + "." + key
}

func (f *gearmandFunction) name() string {
	name := f.function
	// XXX: escape invalid characters
	// Characters which can be used in custom metric names include any alphanumeric characters as well as hyphens (-), underscores (_), and dots (.).
	name = strings.Replace(name, "\t", "-", -1)
	name = strings.Replace(name, ":", "-", -1)
	name = strings.Replace(name, "+", "-", -1)
	name = strings.Replace(name, "=", "-", -1)
	name = strings.Replace(name, "#", "-", -1)
	name = strings.Replace(name, "$", "-", -1)
	return name
}

// GearmandPlugin mackerel plugin for gearmand
type GearmandPlugin struct {
	Target   string
	Socket   string
	Tempfile string
}

func (m GearmandPlugin) connect() (net.Conn, error) {
	network := "tcp"
	target := m.Target
	if m.Socket != "" {
		network = "unix"
		target = m.Socket
	}
	return net.Dial(network, target)
}

// FetchMetrics interface for mackerelplugin
func (m GearmandPlugin) FetchMetrics() (map[string]interface{}, error) {
	conn, err := m.connect()
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(conn, "status")
	return m.parseStats(conn)
}

func (m GearmandPlugin) parseStats(conn io.Reader) (map[string]interface{}, error) {
	scanner := bufio.NewScanner(conn)
	stat := make(map[string]interface{})

	for scanner.Scan() {
		line := scanner.Text()
		if line == "." {
			return stat, nil
		}

		function, err := parseLine(line)
		if err != nil {
			return nil, err
		}

		stat[function.key("available")] = function.available
		stat[function.key("running")] = function.running
		stat[function.key("total")] = function.total
	}
	if err := scanner.Err(); err != nil {
		return stat, err
	}
	return nil, nil
}

func parseLine(line string) (*gearmandFunction, error) {
	// format: FUNCTION\tTOTAL\tRUNNING\tAVAILABLE_WORKERS
	// XXX: function may include some tab characters
	res := reverse(strings.Split(line, "\t"))
	if len(res) < 4 {
		return nil, errors.New("Invalid format: " + line)
	}

	available, err := strconv.ParseUint(res[0], 10, 32)
	if err != nil {
		return nil, err
	}

	running, err := strconv.ParseUint(res[1], 10, 32)
	if err != nil {
		return nil, err
	}

	total, err := strconv.ParseUint(res[2], 10, 32)
	if err != nil {
		return nil, err
	}

	name := strings.Join(reverse(res[3:]), "\t")
	return &gearmandFunction{
		function:  name,
		available: uint32(available),
		running:   uint32(running),
		total:     uint32(total),
	}, nil
}

func reverse(src []string) []string {
	length := len(src)
	dest := make([]string, length)
	for i, v := range src {
		dest[length-i-1] = v
	}
	return dest
}

// GraphDefinition interface for mackerelplugin
func (m GearmandPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "7003", "Port")
	optSocket := flag.String("socket", "", "Server socket (overrides hosts and port)")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var gearmand GearmandPlugin
	if *optSocket != "" {
		gearmand.Socket = *optSocket
	} else {
		gearmand.Target = fmt.Sprintf("%s:%s", *optHost, *optPort)
	}
	helper := mp.NewMackerelPlugin(gearmand)

	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		if gearmand.Socket != "" {
			helper.SetTempfileByBasename(fmt.Sprintf("mackerel-plugin-gearmand-%s", fmt.Sprintf("%x", md5.Sum([]byte(gearmand.Socket)))))
		} else {
			helper.SetTempfileByBasename(fmt.Sprintf("mackerel-plugin-gearmand-%s-%s", *optHost, *optPort))
		}
	}
	helper.Run()
}
