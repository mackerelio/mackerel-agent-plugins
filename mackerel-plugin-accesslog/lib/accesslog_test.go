package mpaccesslog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Songmu/axslogparser"
)

var fetchMetricsTests = []struct {
	Name   string
	InFile string
	Output map[string]float64
}{
	{
		Name:   "Apache log",
		InFile: "testdata/sample-apache.log",
		Output: map[string]float64{
			"total_count":    10,
			"2xx_count":      7,
			"3xx_count":      0,
			"4xx_count":      2,
			"5xx_count":      1,
			"2xx_percentage": 70,
			"3xx_percentage": 0,
			"4xx_percentage": 20,
			"5xx_percentage": 10,
		},
	},
	{
		Name:   "Apache log (loose)",
		InFile: "testdata/sample-apache-loose.log",
		Output: map[string]float64{
			"total_count":    2,
			"2xx_count":      2,
			"3xx_count":      0,
			"4xx_count":      0,
			"5xx_count":      0,
			"2xx_percentage": 100,
			"3xx_percentage": 0,
			"4xx_percentage": 0,
			"5xx_percentage": 0,
		},
	},
	{
		Name:   "LTSV log",
		InFile: "testdata/sample-ltsv.tsv",
		Output: map[string]float64{
			"2xx_count":      7,
			"3xx_count":      1,
			"4xx_count":      1,
			"5xx_count":      1,
			"total_count":    10,
			"2xx_percentage": 70,
			"3xx_percentage": 10,
			"4xx_percentage": 10,
			"5xx_percentage": 10,
			"average":        0.7603999999999999,
			"90_percentile":  2.018,
			"95_percentile":  3.018,
			"99_percentile":  3.018,
		},
	},
	{
		Name:   "LTSV long line log",
		InFile: "testdata/sample-ltsv-long.tsv",
		Output: map[string]float64{
			"2xx_count":      2,
			"3xx_count":      0,
			"4xx_count":      0,
			"5xx_count":      0,
			"total_count":    2,
			"2xx_percentage": 100,
			"3xx_percentage": 0,
			"4xx_percentage": 0,
			"5xx_percentage": 0,
			"average":        0.015,
			"90_percentile":  0.015,
			"95_percentile":  0.015,
			"99_percentile":  0.015,
		},
	},
	{
		Name:   "LTSV log (loose)",
		InFile: "testdata/sample-ltsv-loose.tsv",
		Output: map[string]float64{
			"2xx_count":      3,
			"3xx_count":      0,
			"4xx_count":      0,
			"5xx_count":      0,
			"total_count":    3,
			"2xx_percentage": 100,
			"3xx_percentage": 0,
			"4xx_percentage": 0,
			"5xx_percentage": 0,
			"average":        0.020,
			"90_percentile":  0.025,
			"95_percentile":  0.025,
			"99_percentile":  0.025,
		},
	},
}

func TestFetchMetrics(t *testing.T) {
	for _, tt := range fetchMetricsTests {
		t.Logf("testing: %s", tt.Name)
		p := &AccesslogPlugin{
			file:      tt.InFile,
			noPosFile: true,
		}
		out, err := p.FetchMetrics()
		if err != nil {
			t.Errorf("%s(err): error should be nil but: %+v", tt.Name, err)
			continue
		}
		if !reflect.DeepEqual(out, tt.Output) {
			t.Errorf("%s: \n out:  %#v\n want: %#v", tt.Name, out, tt.Output)
		}
	}
}

func TestFetchMetricsWithCustomParser(t *testing.T) {
	// OK case
	p := &AccesslogPlugin{
		file:      "testdata/sample-ltsv.tsv",
		noPosFile: true,
		parser:    &axslogparser.LTSV{},
	}
	out, err := p.FetchMetrics()
	if err != nil {
		t.Errorf("error should be nil but: %+v", err)
		return
	}

	expected := map[string]float64{
		"2xx_count":      7,
		"3xx_count":      1,
		"4xx_count":      1,
		"5xx_count":      1,
		"total_count":    10,
		"2xx_percentage": 70,
		"3xx_percentage": 10,
		"4xx_percentage": 10,
		"5xx_percentage": 10,
		"average":        0.7603999999999999,
		"90_percentile":  2.018,
		"95_percentile":  3.018,
		"99_percentile":  3.018,
	}
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("out:  %#v\n want: %#v", out, expected)
	}

	// NG case (should not detect log format by log line)
	p = &AccesslogPlugin{
		file:      "testdata/sample-apache.log",
		noPosFile: true,
		parser:    &axslogparser.LTSV{},
	}
	out, err = p.FetchMetrics()
	if err != nil {
		t.Errorf("error should be nil but: %+v", err)
		return
	}

	expected = map[string]float64{
		"2xx_count":   0,
		"3xx_count":   0,
		"4xx_count":   0,
		"5xx_count":   0,
		"total_count": 0,
	}
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("out:  %#v\n want: %#v", out, expected)
	}
}

func TestSkipLogOnceIfNoPos(t *testing.T) {
	dir := t.TempDir()
	posFile := filepath.Join(dir, "plugin-accesslog.test.pos")
	p := &AccesslogPlugin{
		file:    "testdata/sample-ltsv.tsv",
		posFile: posFile,
	}
	out, err := p.FetchMetrics()
	if err != nil {
		t.Errorf("error should be nil but: %+v", err)
		return
	}
	if n := len(out); n != 0 {
		t.Errorf("got %d metrics; but want 0", n)
	}

	// see https://github.com/Songmu/postailer/blob/master/postailer.go#L27-L30
	var pos struct {
		Pos int64 `json:"pos"`
	}
	b, err := os.ReadFile(posFile)
	if err != nil {
		t.Errorf("ReadFile(%s): %v", posFile, err)
		return
	}
	if err := json.Unmarshal(b, &pos); err != nil {
		t.Fatal(err)
	}
	var want int64 = 1247
	if pos.Pos != want {
		t.Errorf("saved position = %d; want %d", pos.Pos, want)
	}
}
