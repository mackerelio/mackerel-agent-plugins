package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"time"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

const (
	namespace          = "AWS/CloudFront"
	region             = "us-east-1"
	metricsTypeAverage = "Average"
	metricsTypeSum     = "Sum"
)

type metrics struct {
	Name string
	Type string
}

// CloudFrontPlugin mackerel plugin for cloudfront
type CloudFrontPlugin struct {
	AccessKeyID     string
	SecretAccessKey string
	CloudWatch      *cloudwatch.CloudWatch
	Name            string
	KeyPrefix       string
	LabelPrefix     string
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p CloudFrontPlugin) MetricKeyPrefix() string {
	if p.KeyPrefix == "" {
		p.KeyPrefix = "cloudfront"
	}
	return p.KeyPrefix
}

func (p *CloudFrontPlugin) prepare() error {
	auth, err := aws.GetAuth(p.AccessKeyID, p.SecretAccessKey, "", time.Now())
	if err != nil {
		return err
	}

	p.CloudWatch, err = cloudwatch.NewCloudWatch(auth, aws.Regions[region].CloudWatchServicepoint)
	if err != nil {
		return err
	}

	return nil
}

func (p CloudFrontPlugin) getLastPoint(metric metrics) (float64, error) {
	now := time.Now()

	dimensions := []cloudwatch.Dimension{
		{
			Name:  "DistributionId",
			Value: p.Name,
		},
		{
			Name:  "Region",
			Value: "Global",
		},
	}

	response, err := p.CloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsRequest{
		Dimensions: dimensions,
		StartTime:  now.Add(time.Duration(180) * time.Second * -1), // 3 min (to fetch at least 1 data-point)
		EndTime:    now,
		MetricName: metric.Name,
		Period:     60,
		Statistics: []string{metric.Type},
		Namespace:  namespace,
	})
	if err != nil {
		return 0, err
	}

	datapoints := response.GetMetricStatisticsResult.Datapoints
	if len(datapoints) == 0 {
		return 0, errors.New("fetched no datapoints")
	}

	// get a least recently datapoint
	// because a most recently datapoint is not stable.
	least := now
	var latestVal float64
	for _, dp := range datapoints {
		if dp.Timestamp.Before(least) {
			least = dp.Timestamp
			if metric.Type == metricsTypeAverage {
				latestVal = dp.Average
			} else if metric.Type == metricsTypeSum {
				latestVal = dp.Sum
			}
		}
	}

	return latestVal, nil
}

// FetchMetrics fetch the metrics
func (p CloudFrontPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})

	for _, met := range [...]metrics{
		{Name: "Requests", Type: metricsTypeSum},
		{Name: "BytesDownloaded", Type: metricsTypeSum},
		{Name: "BytesUploaded", Type: metricsTypeSum},
		{Name: "4xxErrorRate", Type: metricsTypeAverage},
		{Name: "5xxErrorRate", Type: metricsTypeAverage},
	} {
		v, err := p.getLastPoint(met)
		if err == nil {
			stat[met.Name] = v
		} else {
			log.Printf("%s: %s", met, err)
		}
	}

	return stat, nil
}

// GraphDefinition of CloudFrontPlugin
func (p CloudFrontPlugin) GraphDefinition() map[string](mp.Graphs) {
	labelPrefix := p.LabelPrefix

	var graphdef = map[string](mp.Graphs){
		"Requests": mp.Graphs{
			Label: (labelPrefix + " Requests"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "Requests", Label: "Requests"},
			},
		},
		"Transfer": mp.Graphs{
			Label: (labelPrefix + " Transfer"),
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "BytesDownloaded", Label: "Download", Stacked: true},
				mp.Metrics{Name: "BytesUploaded", Label: "Upload", Stacked: true},
			},
		},
		"ErrorRate": mp.Graphs{
			Label: (labelPrefix + " ErrorRate"),
			Unit:  "percentage",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "4xxErrorRate", Label: "4xx", Stacked: true},
				mp.Metrics{Name: "5xxErrorRate", Label: "5xx", Stacked: true},
			},
		},
	}

	return graphdef
}

func main() {
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optIdentifier := flag.String("identifier", "", "Distribution ID")
	optKeyPrefix := flag.String("metric-key-prefix", "cloudfront", "Metric key prefix")
	optLabelPrefix := flag.String("metric-label-prefix", "CloudFront", "Metric label prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var plugin CloudFrontPlugin

	plugin.AccessKeyID = *optAccessKeyID
	plugin.SecretAccessKey = *optSecretAccessKey
	plugin.Name = *optIdentifier
	plugin.LabelPrefix = *optLabelPrefix
	plugin.KeyPrefix = *optKeyPrefix

	helper := mp.NewMackerelPlugin(plugin)
	helper.Tempfile = *optTempfile

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		err := plugin.prepare()
		if err != nil {
			log.Fatalln(err)
		}
		helper.OutputValues()
	}
}
