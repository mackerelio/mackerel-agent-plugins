package mpjson

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// JSONPlugin plugin for JSON
type JSONPlugin struct {
	Target             string
	Tempfile           string
	URL                string
	Prefix             string
	InsecureSkipVerify bool
	ShowOnlyNum        bool
	ExcludeExp         *regexp.Regexp
	IncludeExp         *regexp.Regexp
}

func (p JSONPlugin) traverseMap(content interface{}, path []string) (map[string]float64, error) {
	stat := make(map[string]float64)
	if reflect.TypeOf(content).Kind() == reflect.Slice {
		for i, c := range content.([]interface{}) {
			ts, _ := p.traverseMap(c, append(path, strconv.Itoa(i)))
			for tk, tv := range ts {
				stat[tk] = tv
			}
		}
	} else {
		for k, v := range content.(map[string]interface{}) {
			if v != nil {
				if reflect.TypeOf(v).Kind() == reflect.Map {
					ts, _ := p.traverseMap(v, append(path, k))
					for tk, tv := range ts {
						stat[tk] = tv
					}
				} else if reflect.TypeOf(v).Kind() == reflect.Slice {
					for i, c := range v.([]interface{}) {
						ts, _ := p.traverseMap(c, append(path, strconv.Itoa(i)))
						for tk, tv := range ts {
							stat[tk] = tv
						}
					}
				} else {
					tk, tv := p.outputMetric(strings.Join(append(path, k), "."), v)
					if tk != "" {
						stat[tk] = tv
					}
				}
			}
		}
	}

	return stat, nil
}

func (p JSONPlugin) outputMetric(path string, value interface{}) (string, float64) {
	if p.IncludeExp.MatchString(path) && !p.ExcludeExp.MatchString(path) {
		if reflect.TypeOf(value).Kind() == reflect.Float64 {
			return path, value.(float64)
		}
	}

	return "", 0
}

// FetchMetrics interface for mackerel-plugin
func (p JSONPlugin) FetchMetrics() (map[string]float64, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: p.InsecureSkipVerify},
	}
	client := &http.Client{Transport: tr}
	response, err := client.Get(p.URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var content interface{}
	if err := json.Unmarshal(bytes, &content); err != nil {
		return nil, err
	}

	return p.traverseMap(content, []string{p.Prefix})
}

// Do do doo
func Do() {
	url := flag.String("url", "", "URL to get a JSON")
	prefix := flag.String("prefix", "custom", "Prefix for metric names")
	insecure := flag.Bool("insecure", false, "Skip certificate verifications")
	exclude := flag.String("exclude", `^$`, "Exclude metrics that matches the expression")
	include := flag.String("include", ``, "Output metrics that matches the expression")
	flag.Parse()

	if *url == "" {
		fmt.Println("-url is mandatory")
		os.Exit(1)
	}

	var jsonplugin JSONPlugin

	jsonplugin.URL = *url
	jsonplugin.Prefix = *prefix
	jsonplugin.InsecureSkipVerify = *insecure
	var err error
	jsonplugin.ExcludeExp, err = regexp.Compile(*exclude)
	if err != nil {
		fmt.Printf("exclude expression is invalid: %s", err)
		os.Exit(1)
	}
	jsonplugin.IncludeExp, err = regexp.Compile(*include)
	if err != nil {
		fmt.Printf("include expression is invalid: %s", err)
		os.Exit(1)
	}

	metrics, _ := jsonplugin.FetchMetrics()
	for k, v := range metrics {
		fmt.Printf("%s\t%f\t%d\n", k, v, time.Now().Unix())
	}
}
