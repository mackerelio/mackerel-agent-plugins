package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	//"io/ioutil"
	"errors"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"nginx.connections": mp.Graphs{
		Label: "Nginx Connections",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "connections", Label: "Active connections", Diff: false},
		},
	},
	"nginx.requests": mp.Graphs{
		Label: "Nginx requests",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "accepts", Label: "Accepted connections", Diff: true},
			mp.Metrics{Name: "handled", Label: "Handled connections", Diff: true},
			mp.Metrics{Name: "requests", Label: "Handled requests", Diff: true},
		},
	},
	"nginx.queue": mp.Graphs{
		Label: "Nginx connection status",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "reading", Label: "Reading", Diff: false},
			mp.Metrics{Name: "writing", Label: "Writing", Diff: false},
			mp.Metrics{Name: "waiting", Label: "Waiting", Diff: false},
		},
	},
}

type stringSlice []string

func (s *stringSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func (s *stringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

type NginxPlugin struct {
	Uri    string
	Header stringSlice
}

// % wget -qO- http://localhost:8080/nginx_status
// Active connections: 123
// server accepts handled requests
//  1693613501 1693613501 7996986318
// Reading: 66 Writing: 16 Waiting: 41

func (n NginxPlugin) FetchMetrics() (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", n.Uri, nil)
	if err != nil {
		return nil, err
	}
	for _, h := range n.Header {
		kv := strings.SplitN(h, ":", 2)
		var k, v string
		k = strings.TrimSpace(kv[0])
		if len(kv) == 2 {
			v = strings.TrimSpace(kv[1])
		}
		if http.CanonicalHeaderKey(k) == "Host" {
			req.Host = v
		} else {
			req.Header.Set(k, v)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return n.ParseStats(resp.Body)
}

func (n NginxPlugin) ParseStats(body io.Reader) (map[string]interface{}, error) {
	stat := make(map[string]interface{})

	r := bufio.NewReader(body)
	line, _, err := r.ReadLine()
	if err != nil {
		return nil, errors.New("cannot get values")
	}
	re := regexp.MustCompile("Active connections: ([0-9]+)")
	res := re.FindStringSubmatch(string(line))
	if res == nil || len(res) != 2 {
		return nil, errors.New("cannot get values")
	}
	stat["connections"], err = strconv.ParseFloat(res[1], 64)
	if err != nil {
		return nil, errors.New("cannot get values")
	}

	line, _, err = r.ReadLine()
	if err != nil {
		return nil, errors.New("cannot get values")
	}

	line, _, err = r.ReadLine()
	if err != nil {
		return nil, errors.New("cannot get values")
	}
	re = regexp.MustCompile("([0-9]+) ([0-9]+) ([0-9]+)")
	res = re.FindStringSubmatch(string(line))
	if res == nil || len(res) != 4 {
		return nil, errors.New("cannot get values")
	}
	stat["accepts"], err = strconv.ParseFloat(res[1], 64)
	if err != nil {
		return nil, errors.New("cannot get values")
	}
	stat["handled"], err = strconv.ParseFloat(res[2], 64)
	if err != nil {
		return nil, errors.New("cannot get values")
	}
	stat["requests"], err = strconv.ParseFloat(res[3], 64)
	if err != nil {
		return nil, errors.New("cannot get values")
	}

	line, _, err = r.ReadLine()
	if err != nil {
		return nil, errors.New("cannot get values")
	}
	re = regexp.MustCompile("Reading: ([0-9]+) Writing: ([0-9]+) Waiting: ([0-9]+)")
	res = re.FindStringSubmatch(string(line))
	if res == nil || len(res) != 4 {
		return nil, errors.New("cannot get values")
	}
	stat["reading"], err = strconv.ParseFloat(res[1], 64)
	if err != nil {
		return nil, errors.New("cannot get values")
	}
	stat["writing"], err = strconv.ParseFloat(res[2], 64)
	if err != nil {
		return nil, errors.New("cannot get values")
	}
	stat["wating"], err = strconv.ParseFloat(res[3], 64)
	if err != nil {
		return nil, errors.New("cannot get values")
	}
	return stat, nil
}

func (n NginxPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optUri := flag.String("uri", "", "URI")
	optScheme := flag.String("scheme", "http", "Scheme")
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "8080", "Port")
	optPath := flag.String("path", "/nginx_status", "Path")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optHeader := &stringSlice{}
	flag.Var(optHeader, "header", "Set http header (e.g. \"Host: servername\")")
	flag.Parse()

	var nginx NginxPlugin
	if *optUri != "" {
		nginx.Uri = *optUri
	} else {
		nginx.Uri = fmt.Sprintf("%s://%s:%s%s", *optScheme, *optHost, *optPort, *optPath)
	}
	nginx.Header = *optHeader

	helper := mp.NewMackerelPlugin(nginx)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-nginx")
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
