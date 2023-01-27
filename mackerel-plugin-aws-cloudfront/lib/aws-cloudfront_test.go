package mpawscloudfront

import (
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
	p := CloudFrontPlugin{
		Name:   "AWSCloudFront",
		Prefix: "aws-cloudfront",
	}
	graphdef := p.GraphDefinition()

	expectedLabels := map[string](string){
		"Requests":  "AWS CloudFront Requests",
		"Transfer":  "AWS CloudFront Transfer",
		"ErrorRate": "AWS CloudFront ErrorRate",
	}

	expectedMetricNames := map[string]([]string){
		"Requests": {
			"Requests",
		},
		"Transfer": {
			"BytesDownloaded",
			"BytesUploaded",
		},
		"ErrorRate": {
			"4xxErrorRate",
			"5xxErrorRate",
		},
	}

	expectedMetricLabels := map[string]([]string){
		"Requests": {
			"Requests",
		},
		"Transfer": {
			"Download",
			"Upload",
		},
		"ErrorRate": {
			"4xx",
			"5xx",
		},
	}

	for k, labels := range expectedLabels {
		value, ok := graphdef[k]
		if !ok {
			t.Errorf("graphdef of %s cannot be fetched", k)
			continue
		}

		assert.Equal(t, labels, value.Label)

		var metricNames []string
		var metricLabels []string
		for _, metric := range value.Metrics {
			metricNames = append(metricNames, metric.Name)
			metricLabels = append(metricLabels, metric.Label)
		}

		names := expectedMetricNames[k]
		sort.Strings(names)
		sort.Strings(metricNames)
		if !reflect.DeepEqual(names, metricNames) {
			t.Errorf("graphdef of %s should contain names %v, but %v",
				k, names, metricNames)
		}

		labels := expectedMetricLabels[k]
		sort.Strings(labels)
		sort.Strings(metricLabels)
		if !reflect.DeepEqual(labels, metricLabels) {
			t.Errorf("graphdef of %s should contain labels %v, but %v",
				k, labels, metricLabels)
		}
	}
}
