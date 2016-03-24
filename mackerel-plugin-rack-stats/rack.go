package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var sock string

func parseAddress(uri string) (scheme, path, port string, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return scheme, path, port, err
	}

	scheme = u.Scheme
	path = u.Path
	if path == "" {
		path = u.Host
	}

	host := strings.Split(u.Host, ":")
	if len(host) == 1 {
		port = "80"
	} else {
		port = host[1]
	}

	return scheme, path, port, err
}

func parseBody(r io.Reader, index string) (stats map[string]interface{}, err error) {
	scanner := bufio.NewScanner(r)
	stats = make(map[string]interface{})
	for scanner.Scan() {
		p := strings.Split(scanner.Text(), " ")
		if len(p) == 2 {
			stats[strings.Trim(p[0], ":")], err = strconv.ParseFloat(p[1], 64)
		} else {
			re := regexp.MustCompile(fmt.Sprintf("%s$", index))
			if ok := re.Match([]byte(p[0])); ok && err == nil {
				stats[strings.Trim(p[len(p)-2], ":")], err = strconv.ParseFloat(p[len(p)-1], 64)
			}
		}
	}
	stats["active"] = stats["active"].(float64) - 1

	return stats, err
}

func parseBodyHTTP(uri, port string) (stats map[string]interface{}, err error) {
	var req *http.Request
	req, err = http.NewRequest("GET", uri, nil)
	if err != nil {
		return stats, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return stats, err
	}
	defer resp.Body.Close()

	stats, err = parseBody(resp.Body, ":"+port)

	return stats, err
}

func fakeDial(proto, addr string) (conn net.Conn, err error) {
	return net.Dial("unix", sock)
}

func parseBodyUnix(path string) (stats map[string]interface{}, err error) {
	tr := &http.Transport{
		Dial: fakeDial,
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(fmt.Sprintf("http://dummy/%s", strings.TrimLeft(path, "/")))

	if err != nil {
		return stats, err
	}
	defer resp.Body.Close()

	stats, err = parseBody(resp.Body, sock)

	return stats, err
}

// RackStatsPlugin mackerel plugin for Rack servers
type RackStatsPlugin struct {
	Address   string
	Path      string
	MetricKey string
}

// FetchMetrics interface for mackerelplugin
func (u RackStatsPlugin) FetchMetrics() (stats map[string]interface{}, err error) {
	stats, err = u.parseStats()
	return stats, err
}

func (u RackStatsPlugin) parseStats() (stats map[string]interface{}, err error) {
	scheme, path, port, err := parseAddress(u.Address)

	switch scheme {
	case "http":
		stats, err = parseBodyHTTP(fmt.Sprintf("%s/%s", u.Address, strings.TrimLeft(u.Path, "/")), port)
	case "unix":
		sock = path
		stats, err = parseBodyUnix(u.Path)
	}

	return stats, err
}

// GraphDefinition interface for mackerelplugin
func (u RackStatsPlugin) GraphDefinition() map[string](mp.Graphs) {
	scheme, path, port, err := parseAddress(u.Address)
	if err != nil {
		log.Fatal(err)
	}

	var label string
	if u.MetricKey == "" {
		switch scheme {
		case "http":
			u.MetricKey = port
			label = fmt.Sprintf("Rack Port %s Stats", port)
		case "unix":
			u.MetricKey = strings.Replace(strings.Replace(path, "/", "_", -1), ".", "_", -1)
			label = fmt.Sprintf("Rack %s Stats", path)
		}
	} else {
		label = fmt.Sprintf("Rack %s Stats", u.MetricKey)
	}

	return map[string](mp.Graphs){
		fmt.Sprintf("rack.%s.stats", u.MetricKey): mp.Graphs{
			Label: label,
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "queued", Label: "Queued", Diff: false},
				mp.Metrics{Name: "active", Label: "Active", Diff: false},
				mp.Metrics{Name: "writing", Label: "Writing", Diff: false},
				mp.Metrics{Name: "calling", Label: "Calling", Diff: false},
			},
		},
	}
}

func main() {
	optAddress := flag.String("address", "http://localhost:8080", "URL or Unix Domain Socket")
	optPath := flag.String("path", "/_raindrops", "Path")
	optMetricKey := flag.String("metric-key-prefix", "", "Metric Key Prefix")
	optVersion := flag.Bool("version", false, "Version")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	if *optVersion {
		fmt.Println("0.3")
		os.Exit(0)
	}

	var rack RackStatsPlugin
	rack.Address = *optAddress
	rack.Path = *optPath
	rack.MetricKey = *optMetricKey

	helper := mp.NewMackerelPlugin(rack)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-rack-stats")
	}

	helper.Run()
}
