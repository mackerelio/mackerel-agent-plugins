package mpaccesslog

import (
	"reflect"
	"testing"
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
			"90_percentile":  3.018,
			"95_percentile":  4.018,
			"99_percentile":  4.018,
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
			t.Errorf("%s(err): error should be nil but: %+v", err)
			continue
		}
		if !reflect.DeepEqual(out, tt.Output) {
			t.Errorf("%s: \n out:  %#v\n want: %#v", tt.Name, out, tt.Output)
		}
	}
}
