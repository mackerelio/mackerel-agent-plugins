package mpgunicorn

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"os"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

// GunicornPlugin mackerel plugin for Gunicorn
type GunicornPlugin struct {
	FilePath    string
	Prefix      string
	LabelPrefix string
}

// {
//   "BusyWorkers": "0",
//   "IdleWorkers": "4",
//   "TotalAccesses": "440"
//   "stats": [
//     {
//       "uri": "/",
//       "method": "GET",
//       "host": "localhost:8000",
//       "pid": 26727,
//       "status": "_"
//     },
//     {
//       "uri": "/",
//       "method": "GET",
//       "host": "localhost:8000",
//       "pid": 26725,
//       "status": "_"
//     },
//     {
//       "uri": "/",
//       "method": "GET",
//       "host": "localhost:8000",
//       "pid": 26726,
//       "status": "_"
//     },
//     {
//       "host": "localhost:8000",
//       "method": "GET",
//       "uri": "/",
//       "pid": 26724,
//       "status": "_"
//     }
//   ],
// }

// field types vary between versions

// GunicornWorkerStatus struct
type GunicornWorkerStatus struct{}

// GunicornStatus sturct for file json
type GunicornStatus struct {
	TotalAccesses string                 `json:"TotalAccesses"`
	BusyWorkers   string                 `json:"BusyWorkers"`
	IdleWorkers   string                 `json:"IdleWorkers"`
	Stats         []GunicornWorkerStatus `json:"stats"`
}

// FetchMetrics interface for mackerelplugin
func (p GunicornPlugin) FetchMetrics() (map[string]interface{}, error) {
	file, err := os.Open(p.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return p.parseStats(file)
}

func (p GunicornPlugin) parseStats(body io.Reader) (map[string]interface{}, error) {
	stat := make(map[string]interface{})
	decoder := json.NewDecoder(body)

	var s GunicornStatus
	err := decoder.Decode(&s)
	if err != nil {
		return nil, err
	}

	stat["busy_workers"], err = strconv.ParseFloat(s.BusyWorkers, 64)
	if err != nil {
		return nil, errors.New("cannot get values")
	}

	stat["idle_workers"], err = strconv.ParseFloat(s.IdleWorkers, 64)
	if err != nil {
		return nil, errors.New("cannot get values")
	}

	stat["requests"], err = strconv.ParseUint(s.TotalAccesses, 10, 64)
	if err != nil {
		return nil, errors.New("cannot get values")
	}

	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (p GunicornPlugin) GraphDefinition() map[string]mp.Graphs {
	var graphdef = map[string]mp.Graphs{
		(p.Prefix + ".workers"): {
			Label: p.LabelPrefix + " Workers",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "busy_workers", Label: "Busy Workers", Diff: false, Stacked: true},
				{Name: "idle_workers", Label: "Idle Workers", Diff: false, Stacked: true},
			},
		},
		(p.Prefix + ".req"): {
			Label: p.LabelPrefix + " Requests",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "requests", Label: "Requests", Diff: true, Type: "uint64"},
			},
		},
	}

	return graphdef
}

// Do the plugin
func Do() {
	optFilePath := flag.String("status-file", "", "FilePath")
	optPrefix := flag.String("metric-key-prefix", "gunicorn", "Prefix")
	optLabelPrefix := flag.String("metric-label-prefix", "", "Label Prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	gunicorn := GunicornPlugin{FilePath: *optFilePath, Prefix: *optPrefix, LabelPrefix: *optLabelPrefix}
	if gunicorn.LabelPrefix == "" {
		gunicorn.LabelPrefix = strings.Title(gunicorn.Prefix)
	}

	helper := mp.NewMackerelPlugin(gunicorn)
	helper.Tempfile = *optTempfile

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
