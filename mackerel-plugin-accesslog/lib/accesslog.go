package mpaccesslog

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Songmu/axslogparser"
	"github.com/Songmu/postailer"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/mackerelio/golib/pluginutil"
	"github.com/montanaflynn/stats"
)

// AccesslogPlugin mackerel plugin
type AccesslogPlugin struct {
	prefix    string
	file      string
	posFile   string
	parser    axslogparser.Parser
	noPosFile bool
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p *AccesslogPlugin) MetricKeyPrefix() string {
	if p.prefix == "" {
		p.prefix = "accesslog"
	}
	return p.prefix
}

// GraphDefinition interface for mackerelplugin
func (p *AccesslogPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(p.prefix)
	return map[string]mp.Graphs{
		"access_num": {
			Label: labelPrefix + " Access Num",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "total_count", Label: "Total Count"},
				{Name: "5xx_count", Label: "HTTP 5xx Count", Stacked: true},
				{Name: "4xx_count", Label: "HTTP 4xx Count", Stacked: true},
				{Name: "3xx_count", Label: "HTTP 3xx Count", Stacked: true},
				{Name: "2xx_count", Label: "HTTP 2xx Count", Stacked: true},
			},
		},
		"access_rate": {
			Label: labelPrefix + " Access Rate",
			Unit:  "percentage",
			Metrics: []mp.Metrics{
				{Name: "5xx_percentage", Label: "HTTP 5xx Percentage", Stacked: true},
				{Name: "4xx_percentage", Label: "HTTP 4xx Percentage", Stacked: true},
				{Name: "3xx_percentage", Label: "HTTP 3xx Percentage", Stacked: true},
				{Name: "2xx_percentage", Label: "HTTP 2xx Percentage", Stacked: true},
			},
		},
		"latency": {
			Label: labelPrefix + " Latency",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "99_percentile", Label: "99 Percentile"},
				{Name: "95_percentile", Label: "95 Percentile"},
				{Name: "90_percentile", Label: "90 Percentile"},
				{Name: "average", Label: "Average"},
			},
		},
	}
}

var posRe = regexp.MustCompile(`^([a-zA-Z]):[/\\]`)

func (p *AccesslogPlugin) getPosPath() string {
	base := p.file + ".pos.json"
	if p.posFile != "" {
		if filepath.IsAbs(p.posFile) {
			return p.posFile
		}
		base = p.posFile
	}
	return filepath.Join(
		pluginutil.PluginWorkDir(),
		"mackerel-plugin-accesslog.d",
		posRe.ReplaceAllString(base, `$1`+string(filepath.Separator)),
	)
}

func (p *AccesslogPlugin) getReadCloser() (io.ReadCloser, bool, error) {
	if p.noPosFile {
		rc, err := os.Open(p.file)
		return rc, true, err
	}
	posfile := p.getPosPath()
	fi, err := os.Stat(posfile)
	// don't output any metrics when the pos file doesn't exist or is too old
	takeMetrics := err == nil && fi.ModTime().After(time.Now().Add(-2*time.Minute))
	rc, err := postailer.Open(p.file, posfile)
	return rc, takeMetrics, err
}

// FetchMetrics interface for mackerelplugin
func (p *AccesslogPlugin) FetchMetrics() (map[string]float64, error) {
	rc, takeMetrics, err := p.getReadCloser()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	if !takeMetrics {
		// discard existing contents to seek position
		_, err := ioutil.ReadAll(rc)
		return map[string]float64{}, err
	}

	countMetrics := []string{"total_count", "2xx_count", "3xx_count", "4xx_count", "5xx_count"}
	ret := make(map[string]float64)
	for _, k := range countMetrics {
		ret[k] = 0
	}
	var reqtimes []float64
	r := bufio.NewReader(rc)
	for {
		var (
			l   *axslogparser.Log
			err error
			bb  bytes.Buffer
		)
		buf, isPrefix, err := r.ReadLine()
		bb.Write(buf)
		for isPrefix {
			buf, isPrefix, err = r.ReadLine()
			if err != nil {
				break
			}
			bb.Write(buf)
		}
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		line := bb.String()
		if p.parser == nil {
			p.parser, l, err = axslogparser.GuessParser(line)
		} else {
			l, err = p.parser.Parse(line)
		}
		if err != nil {
			log.Println(err)
			continue
		}
		ret[string(fmt.Sprintf("%d", l.Status)[0])+"xx_count"]++
		ret["total_count"]++

		if l.ReqTime != nil {
			reqtimes = append(reqtimes, *l.ReqTime)
		} else if l.TakenSec != nil {
			reqtimes = append(reqtimes, *l.TakenSec)
		}
	}
	if ret["total_count"] > 0 {
		for _, v := range []string{"2xx", "3xx", "4xx", "5xx"} {
			ret[v+"_percentage"] = ret[v+"_count"] * 100 / ret["total_count"]
		}
	}
	if len(reqtimes) > 0 {
		ret["average"], _ = stats.Mean(reqtimes)
		for _, v := range []int{90, 95, 99} {
			ret[fmt.Sprintf("%d", v)+"_percentile"], _ = stats.Percentile(reqtimes, float64(v))
		}
	}
	return ret, nil
}

// Do the plugin
func Do() {
	var (
		optPrefix    = flag.String("metric-key-prefix", "", "Metric key prefix")
		optFormat    = flag.String("format", "", "Access Log format ('ltsv' or 'apache')")
		optPosFile   = flag.String("posfile", "", "(not necessary to specify it in the usual use case) posfile")
		optNoPosFile = flag.Bool("no-posfile", false, "no position file")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION] /path/to/access.log\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	var parser axslogparser.Parser
	switch *optFormat {
	case "":
		parser = nil // guess format by log (default)
	case "ltsv":
		parser = &axslogparser.LTSV{}
	case "apache":
		parser = &axslogparser.Apache{}
	default:
		fmt.Fprintf(os.Stderr, "Error: '%s' is invalid format name\n", *optFormat)
		flag.Usage()
		os.Exit(1)
	}

	mp.NewMackerelPlugin(&AccesslogPlugin{
		prefix:    *optPrefix,
		file:      flag.Args()[0],
		posFile:   *optPosFile,
		noPosFile: *optNoPosFile,
		parser:    parser,
	}).Run()
}
