package mpawskinesisstreams

import (
	"errors"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

const (
	namespace          = "AWS/Kinesis"
	metricsTypeAverage = "Average"
	metricsTypeMaximum = "Maximum"
	metricsTypeMinimum = "Minimum"
)

type metrics struct {
	CloudWatchName string
	MackerelName   string
	Type           string
}

// KinesisStreamsPlugin mackerel plugin for aws kinesis
type KinesisStreamsPlugin struct {
	Name   string
	Prefix string

	AccessKeyID     string
	SecretAccessKey string
	Region          string
	CloudWatch      *cloudwatch.CloudWatch
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p KinesisStreamsPlugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "kinesis-streams"
	}
	return p.Prefix
}

// prepare creates CloudWatch instance
func (p *KinesisStreamsPlugin) prepare() error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	config := aws.NewConfig()
	if p.AccessKeyID != "" && p.SecretAccessKey != "" {
		config = config.WithCredentials(credentials.NewStaticCredentials(p.AccessKeyID, p.SecretAccessKey, ""))
	}
	if p.Region != "" {
		config = config.WithRegion(p.Region)
	}

	p.CloudWatch = cloudwatch.New(sess, config)

	return nil
}

// getLastPoint fetches a CloudWatch metric and parse
func (p KinesisStreamsPlugin) getLastPoint(metric metrics) (float64, error) {
	now := time.Now()

	dimensions := []*cloudwatch.Dimension{
		{
			Name:  aws.String("StreamName"),
			Value: aws.String(p.Name),
		},
	}

	response, err := p.CloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Dimensions: dimensions,
		StartTime:  aws.Time(now.Add(time.Duration(180) * time.Second * -1)), // 3 min
		EndTime:    aws.Time(now),
		MetricName: aws.String(metric.CloudWatchName),
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

	latest := new(time.Time)
	var latestVal float64
	for _, dp := range datapoints {
		if dp.Timestamp.Before(*latest) {
			continue
		}

		latest = dp.Timestamp
		switch metric.Type {
		case metricsTypeAverage:
			latestVal = *dp.Average
		case metricsTypeMaximum:
			latestVal = *dp.Maximum
		case metricsTypeMinimum:
			latestVal = *dp.Minimum
		}
	}

	return latestVal, nil
}

// FetchMetrics fetch the metrics
func (p KinesisStreamsPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})

	for _, met := range [...]metrics{
		{CloudWatchName: "GetRecords.Bytes", MackerelName: "GetRecordsBytes", Type: metricsTypeAverage},
		// Max of IteratorAgeMilliseconds is useful especially when few of iterators are in trouble
		{CloudWatchName: "GetRecords.IteratorAgeMilliseconds", MackerelName: "GetRecordsDelayMaxMilliseconds", Type: metricsTypeMaximum},
		{CloudWatchName: "GetRecords.IteratorAgeMilliseconds", MackerelName: "GetRecordsDelayMinMilliseconds", Type: metricsTypeMinimum},
		{CloudWatchName: "GetRecords.IteratorAgeMilliseconds", MackerelName: "GetRecordsDelayAverageMilliseconds", Type: metricsTypeAverage},
		{CloudWatchName: "GetRecords.Latency", MackerelName: "GetRecordsLatency", Type: metricsTypeAverage},
		{CloudWatchName: "GetRecords.Records", MackerelName: "GetRecordsRecords", Type: metricsTypeAverage},
		{CloudWatchName: "GetRecords.Success", MackerelName: "GetRecordsSuccess", Type: metricsTypeAverage},
		{CloudWatchName: "IncomingBytes", MackerelName: "IncomingBytes", Type: metricsTypeAverage},
		{CloudWatchName: "IncomingRecords", MackerelName: "IncomingRecords", Type: metricsTypeAverage},
		{CloudWatchName: "PutRecord.Bytes", MackerelName: "PutRecordBytes", Type: metricsTypeAverage},
		{CloudWatchName: "PutRecord.Latency", MackerelName: "PutRecordLatency", Type: metricsTypeAverage},
		{CloudWatchName: "PutRecord.Success", MackerelName: "PutRecordSuccess", Type: metricsTypeAverage},
		{CloudWatchName: "PutRecords.Bytes", MackerelName: "PutRecordsBytes", Type: metricsTypeAverage},
		{CloudWatchName: "PutRecords.Latency", MackerelName: "PutRecordsLatency", Type: metricsTypeAverage},
		{CloudWatchName: "PutRecords.Records", MackerelName: "PutRecordsRecords", Type: metricsTypeAverage},
		{CloudWatchName: "PutRecords.Success", MackerelName: "PutRecordsSuccess", Type: metricsTypeAverage},
		{CloudWatchName: "ReadProvidionedThroughputExceeded", MackerelName: "ReadThroughputExceeded", Type: metricsTypeAverage},
		{CloudWatchName: "WriteProvidionedThroughputExceeded", MackerelName: "WriteThroughputExceeded", Type: metricsTypeAverage},
	} {
		v, err := p.getLastPoint(met)
		if err == nil {
			stat[met.MackerelName] = v
		} else {
			log.Printf("%s: %s", met, err)
		}
	}
	return stat, nil
}

// GraphDefinition of KinesisStreamsPlugin
func (p KinesisStreamsPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(p.Prefix)
	labelPrefix = strings.Replace(labelPrefix, "-", " ", -1)

	var graphdef = map[string]mp.Graphs{
		"bytes": {
			Label: (labelPrefix + " Bytes"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "GetRecordsBytes", Label: "GetRecords"},
				{Name: "IncomingBytes", Label: "Total Incoming"},
				{Name: "PutRecordBytes", Label: "PutRecord"},
				{Name: "PutRecordsBytes", Label: "PutRecords"},
			},
		},
		"iteratorage": {
			Label: (labelPrefix + " Read Delay"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "GetRecordsDelayAverageMilliseconds", Label: "Average"},
				{Name: "GetRecordsDelayMaxMilliseconds", Label: "Max"},
				{Name: "GetRecordsDelayMinMilliseconds", Label: "Min"},
			},
		},
		"latency": {
			Label: (labelPrefix + " Operation Latency"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "GetRecordsLatency", Label: "GetRecords"},
				{Name: "PutRecordLatency", Label: "PutRecord"},
				{Name: "PutRecordsLatency", Label: "PutRecords"},
			},
		},
		"records": {
			Label: (labelPrefix + " Records"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "GetRecordsRecords", Label: "GetRecords"},
				{Name: "IncomingRecords", Label: "Total Incoming"},
				{Name: "PutRecordsRecords", Label: "PutRecords"},
			},
		},
		"success": {
			Label: (labelPrefix + " Operation Success"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "GetRecordsSuccess", Label: "GetRecords"},
				{Name: "PutRecordSuccess", Label: "PutRecord"},
				{Name: "PutRecordsSuccess", Label: "PutRecords"},
			},
		},
		"pending": {
			Label: (labelPrefix + " Pending Operations"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "ReadThroughputExceeded", Label: "Read"},
				{Name: "WriteThroughputExceeded", Label: "Write"},
			},
		},
	}
	return graphdef
}

// Do the plugin
func Do() {
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optRegion := flag.String("region", "", "AWS Region")
	optIdentifier := flag.String("identifier", "", "Stream Name")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optPrefix := flag.String("metric-key-prefix", "kinesis-streams", "Metric key prefix")
	flag.Parse()

	var plugin KinesisStreamsPlugin

	plugin.AccessKeyID = *optAccessKeyID
	plugin.SecretAccessKey = *optSecretAccessKey
	plugin.Region = *optRegion
	plugin.Name = *optIdentifier
	plugin.Prefix = *optPrefix

	err := plugin.prepare()
	if err != nil {
		log.Fatalln(err)
	}

	helper := mp.NewMackerelPlugin(plugin)
	helper.Tempfile = *optTempfile

	helper.Run()
}
