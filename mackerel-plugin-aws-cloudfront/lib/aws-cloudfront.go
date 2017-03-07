package mpawscloudfront

import (
	"errors"
	"flag"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

const (
	namespace          = "AWS/CloudFront"
	region             = "us-east-1"
	metricsTypeAverage = "Average"
	metricsTypeSum     = "Sum"
)

var graphdef = map[string]mp.Graphs{
	"cloudfront.Requests": {
		Label: "CloudFront Requests",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "Requests", Label: "Requests"},
		},
	},
	"cloudfront.Transfer": {
		Label: "CloudFront Transfer",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "BytesDownloaded", Label: "Download", Stacked: true},
			{Name: "BytesUploaded", Label: "Upload", Stacked: true},
		},
	},
	"cloudfront.ErrorRate": {
		Label: "CloudFront ErrorRate",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "4xxErrorRate", Label: "4xx", Stacked: true},
			{Name: "5xxErrorRate", Label: "5xx", Stacked: true},
		},
	},
}

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
}

func (p *CloudFrontPlugin) prepare() error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	config := aws.NewConfig()
	if p.AccessKeyID != "" && p.SecretAccessKey != "" {
		config = config.WithCredentials(credentials.NewStaticCredentials(p.AccessKeyID, p.SecretAccessKey, ""))
	}
	config = config.WithRegion(region)

	p.CloudWatch = cloudwatch.New(sess, config)

	return nil
}

func (p CloudFrontPlugin) getLastPoint(metric metrics) (float64, error) {
	now := time.Now()

	dimensions := []*cloudwatch.Dimension{
		{
			Name:  aws.String("DistributionId"),
			Value: aws.String(p.Name),
		},
		{
			Name:  aws.String("Region"),
			Value: aws.String("Global"),
		},
	}

	response, err := p.CloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Dimensions: dimensions,
		StartTime:  aws.Time(now.Add(time.Duration(180) * time.Second * -1)), // 3 min (to fetch at least 1 data-point)
		EndTime:    aws.Time(now),
		MetricName: aws.String(metric.Name),
		Period:     aws.Int64(60),
		Statistics: []*string{aws.String(metric.Type)},
		Namespace:  aws.String(namespace),
	})
	if err != nil {
		return 0, err
	}

	datapoints := response.Datapoints
	if len(datapoints) == 0 {
		return 0, errors.New("fetched no datapoints")
	}

	// get a least recently datapoint
	// because a most recently datapoint is not stable.
	least := time.Now()
	var latestVal float64
	for _, dp := range datapoints {
		if dp.Timestamp.Before(least) {
			least = *dp.Timestamp
			if metric.Type == metricsTypeAverage {
				latestVal = *dp.Average
			} else if metric.Type == metricsTypeSum {
				latestVal = *dp.Sum
			}
		}
	}

	return latestVal, nil
}

// FetchMetrics fetch the metrics
func (p CloudFrontPlugin) FetchMetrics() (map[string]float64, error) {
	stat := make(map[string]float64)

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
func (p CloudFrontPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// Do the plugin
func Do() {
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optIdentifier := flag.String("identifier", "", "Distribution ID")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var plugin CloudFrontPlugin

	plugin.AccessKeyID = *optAccessKeyID
	plugin.SecretAccessKey = *optSecretAccessKey
	plugin.Name = *optIdentifier

	err := plugin.prepare()
	if err != nil {
		log.Fatalln(err)
	}

	helper := mp.NewMackerelPlugin(plugin)
	helper.Tempfile = *optTempfile

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
