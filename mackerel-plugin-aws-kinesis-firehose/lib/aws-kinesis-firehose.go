package mpawskinesisfirehose

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

const (
	namespace = "AWS/Firehose"
)

type metrics struct {
	CloudWatchName string
	MackerelName   string
}

// KinesisFirehosePlugin mackerel plugin for aws kinesis firehose
type KinesisFirehosePlugin struct {
	Name   string
	Prefix string

	AccessKeyID     string
	SecretAccessKey string
	Region          string
	CloudWatch      *cloudwatch.CloudWatch
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p KinesisFirehosePlugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "kinesis-firehose"
	}
	return p.Prefix
}

func (p *KinesisFirehosePlugin) prepare() error {
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

func getLastPointFromCloudWatch(cw cloudwatchiface.CloudWatchAPI, deliveryStreamName string, metric metrics) (*cloudwatch.Datapoint, error) {
	now := time.Now()

	dimensions := []*cloudwatch.Dimension{
		{
			Name:  aws.String("DeliveryStreamName"),
			Value: aws.String(deliveryStreamName),
		},
	}

	response, err := cw.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Dimensions: dimensions,
		StartTime:  aws.Time(now.Add(time.Duration(480) * time.Second * -1)), // 8 min, since some metrics are aggregated over 5 min
		EndTime:    aws.Time(now),
		MetricName: aws.String(metric.CloudWatchName),
		Period:     aws.Int64(60),
		Statistics: []*string{aws.String("Average")},
		Namespace:  aws.String(namespace),
	})
	if err != nil {
		return nil, err
	}

	datapoints := response.Datapoints
	if len(datapoints) == 0 {
		return nil, nil
	}

	latest := new(time.Time)
	var latestDp *cloudwatch.Datapoint
	for _, dp := range datapoints {
		if dp.Timestamp.Before(*latest) {
			continue
		}

		latest = dp.Timestamp
		latestDp = dp
	}

	return latestDp, nil
}

func mergeStatsFromDatapoint(stats map[string]interface{}, dp *cloudwatch.Datapoint, metric metrics) map[string]interface{} {
	if dp != nil {
		stats[metric.MackerelName] = *dp.Average
	}
	return stats
}

// FetchMetrics fetch the metrics
func (p KinesisFirehosePlugin) FetchMetrics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	for _, met := range [...]metrics{
		{CloudWatchName: "DeliveryToElasticsearch.Bytes", MackerelName: "DeliveryToElasticsearchBytes"},
		{CloudWatchName: "DeliveryToElasticsearch.Records", MackerelName: "DeliveryToElasticsearchRecords"},
		{CloudWatchName: "DeliveryToElasticsearch.Success", MackerelName: "DeliveryToElasticsearchSuccess"},
		{CloudWatchName: "DeliveryToRedshift.Bytes", MackerelName: "DeliveryToRedshiftBytes"},
		{CloudWatchName: "DeliveryToRedshift.Records", MackerelName: "DeliveryToRedshiftRecords"},
		{CloudWatchName: "DeliveryToRedshift.Success", MackerelName: "DeliveryToRedshiftSuccess"},
		{CloudWatchName: "DeliveryToS3.Bytes", MackerelName: "DeliveryToS3Bytes"},
		{CloudWatchName: "DeliveryToS3.DataFreshness", MackerelName: "DeliveryToS3DataFreshness"},
		{CloudWatchName: "DeliveryToS3.Records", MackerelName: "DeliveryToS3Records"},
		{CloudWatchName: "DeliveryToS3.Success", MackerelName: "DeliveryToS3Success"},
		{CloudWatchName: "IncomingBytes", MackerelName: "IncomingBytes"},
		{CloudWatchName: "IncomingRecords", MackerelName: "IncomingRecords"},
		{CloudWatchName: "DescribeDeliveryStream.Latency", MackerelName: "DescribeDeliveryStreamLatency"},
		{CloudWatchName: "DescribeDeliveryStream.Requests", MackerelName: "DescribeDeliveryStreamRequests"},
		{CloudWatchName: "PutRecord.Bytes", MackerelName: "PutRecordBytes"},
		{CloudWatchName: "PutRecord.Latency", MackerelName: "PutRecordLatency"},
		{CloudWatchName: "PutRecord.Requests", MackerelName: "PutRecordRequests"},
		{CloudWatchName: "PutRecordBatch.Bytes", MackerelName: "PutRecordBatchBytes"},
		{CloudWatchName: "PutRecordBatch.Latency", MackerelName: "PutRecordBatchLatency"},
		{CloudWatchName: "PutRecordBatch.Records", MackerelName: "PutRecordBatchRecords"},
		{CloudWatchName: "PutRecordBatch.Requests", MackerelName: "PutRecordBatchRequests"},
		{CloudWatchName: "UpdateDeliveryStream.Latency", MackerelName: "UpdateDeliveryStreamLatency"},
		{CloudWatchName: "UpdateDeliveryStream.Requests", MackerelName: "UpdateDeliveryStreamRequests"},
	} {
		v, err := getLastPointFromCloudWatch(p.CloudWatch, p.Name, met)
		if err == nil {
			stats = mergeStatsFromDatapoint(stats, v, met)
		} else {
			log.Printf("%s: %s", met, err)
		}
	}

	return stats, nil
}

// GraphDefinition of KinesisFirehosePlugin
func (p KinesisFirehosePlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(p.Prefix)
	labelPrefix = strings.Replace(labelPrefix, "-", " ", -1)

	var graphdef = map[string]mp.Graphs{
		"bytes": {
			Label: (labelPrefix + " Bytes"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "DeliveryToElasticsearchBytes", Label: "DeliveryToElasticsearch"},
				{Name: "DeliveryToRedshiftBytes", Label: "DeliveryToRedshift"},
				{Name: "DeliveryToS3", Label: "DeliveryToS3"},
				{Name: "IncomingBytes", Label: "Total Incoming"},
				{Name: "PutRecordBytes", Label: "PutRecord"},
				{Name: "PutRecordBatchBytes", Label: "PutRecordBatch"},
			},
		},
		"records": {
			Label: (labelPrefix + " Records"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "DeliveryToElasticsearchRecords", Label: "DeliveryToElasticsearch"},
				{Name: "DeliveryToRedshiftRecords", Label: "DeliveryToRedshift"},
				{Name: "DeliveryToS3Records", Label: "DeliveryToS3"},
				{Name: "IncomingRecords", Label: "Total Incoming"},
				{Name: "PutRecordRecords", Label: "PutRecord"},
			},
		},
		"success": {
			Label: (labelPrefix + " Success"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "DeliveryToElasticsearchSuccess", Label: "DeliveryToElasticsearch"},
				{Name: "DeliveryToRedshiftSuccess", Label: "DeliveryToRedshift"},
				{Name: "DeliveryToS3Success", Label: "DeliveryToS3"},
			},
		},
		"datafreshness": {
			Label: (labelPrefix + " DataFreshness"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "DeliveryToS3DataFreshness", Label: "DeliveryToS3"},
			},
		},
		"latency": {
			Label: (labelPrefix + " Operation Latency"),
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "DescribeDeliveryStreamLatency", Label: "DescribeDeliveryStream"},
				{Name: "PutRecordLatency", Label: "PutRecord"},
				{Name: "PutRecordBatchLatency", Label: "PutRecordBatch"},
				{Name: "UpdateDeliveryStreamLatency", Label: "UpdateDeliveryStream"},
			},
		},
		"requests": {
			Label: (labelPrefix + " Requests"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "DescribeDeliveryStreamRequests", Label: "DescribeDeliveryStream"},
				{Name: "PutRecordRequests", Label: "PutRecord"},
				{Name: "PutRecordBatchRequests", Label: "PutRecordBatch"},
				{Name: "UpdateDeliveryStreamRequests", Label: "UpdateDeliveryStream"},
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
	optIdentifier := flag.String("identifier", "", "Delivery Stream Name")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optPrefix := flag.String("metric-key-prefix", "kinesis-firehose", "Metric key prefix")
	flag.Parse()

	var plugin KinesisFirehosePlugin

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
