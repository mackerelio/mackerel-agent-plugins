package mpawss3requests

import (
	"errors"
	"flag"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

const (
	namespace          = "AWS/S3"
	metricsTypeAverage = "Average"
	metricsTypeSum     = "Sum"
	metricsTypeMaximum = "Maximum"
	metricsTypeMinimum = "Minimum"
)

// has 1 CloudWatch MetricName and corresponding N Mackerel Metrics
type metricsGroup struct {
	CloudWatchName string
	Metrics        []metric
}

type metric struct {
	MackerelName string
	Type         string
}

// S3RequestsPlugin is mackerel plugin for aws s3 with metric configuration
type S3RequestsPlugin struct {
	BucketName  string
	FilterID    string
	KeyPrefix   string
	LabelPrefix string

	AccessKeyID     string
	SecretAccessKey string
	Region          string

	CloudWatch *cloudwatch.CloudWatch
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p S3RequestsPlugin) MetricKeyPrefix() string {
	if p.KeyPrefix == "" {
		return "s3-requests"
	}
	return p.KeyPrefix
}

// MetricLabelPrefix
func (p S3RequestsPlugin) MetricLabelPrefix() string {
	if p.LabelPrefix == "" {
		return "S3"
	}
	return p.LabelPrefix
}

// prepare creates CloudWatch instance
func (p *S3RequestsPlugin) prepare() error {
	// validate params
	// apparently we need BucketName and FilterID
	if p.BucketName == "" || p.FilterID == "" {
		return errors.New("Both --bucket-name and --filter-id are necessary")
	}

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
func getLastPointFromCloudWatch(cw cloudwatchiface.CloudWatchAPI, bucketName string, filterID string, metric metricsGroup) (*cloudwatch.Datapoint, error) {
	now := time.Now()
	statsInput := make([]*string, len(metric.Metrics))
	for i, typ := range metric.Metrics {
		statsInput[i] = aws.String(typ.Type)
	}
	input := &cloudwatch.GetMetricStatisticsInput{
		StartTime:  aws.Time(now.Add(time.Duration(180) * time.Second * -1)), // 3 min
		EndTime:    aws.Time(now),
		MetricName: aws.String(metric.CloudWatchName),
		Period:     aws.Int64(600),
		Statistics: statsInput,
		Namespace:  aws.String(namespace),
	}
	input.Dimensions = []*cloudwatch.Dimension{
		{
			Name:  aws.String("BucketName"),
			Value: aws.String(bucketName),
		},
		{
			Name:  aws.String("FilterId"),
			Value: aws.String(filterID),
		},
	}
	response, err := cw.GetMetricStatistics(input)
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

func mergeStatsFromDatapoint(stats map[string]float64, dp *cloudwatch.Datapoint, mg metricsGroup) map[string]float64 {
	for _, met := range mg.Metrics {
		switch met.Type {
		case metricsTypeAverage:
			stats[met.MackerelName] = *dp.Average
		case metricsTypeSum:
			stats[met.MackerelName] = *dp.Sum
		case metricsTypeMaximum:
			stats[met.MackerelName] = *dp.Maximum
		case metricsTypeMinimum:
			stats[met.MackerelName] = *dp.Minimum
		}
	}
	return stats
}

var s3RequestMetricsGroup = []metricsGroup{
	{CloudWatchName: "AllRequests", Metrics: []metric{
		{MackerelName: "AllRequests", Type: metricsTypeSum},
	}},
	{CloudWatchName: "GetRequests", Metrics: []metric{
		{MackerelName: "GetRequests", Type: metricsTypeSum},
	}},
	{CloudWatchName: "PutRequests", Metrics: []metric{
		{MackerelName: "PutRequests", Type: metricsTypeSum},
	}},
	{CloudWatchName: "DeleteRequests", Metrics: []metric{
		{MackerelName: "DeleteRequests", Type: metricsTypeSum},
	}},
	{CloudWatchName: "HeadRequests", Metrics: []metric{
		{MackerelName: "HeadRequests", Type: metricsTypeSum},
	}},
	{CloudWatchName: "PostRequests", Metrics: []metric{
		{MackerelName: "PostRequests", Type: metricsTypeSum},
	}},
	{CloudWatchName: "ListRequests", Metrics: []metric{
		{MackerelName: "ListRequests", Type: metricsTypeSum},
	}},
	{CloudWatchName: "4xxErrors", Metrics: []metric{
		{MackerelName: "4xxErrors", Type: metricsTypeSum},
	}},
	{CloudWatchName: "5xxErrors", Metrics: []metric{
		{MackerelName: "5xxErrors", Type: metricsTypeSum},
	}},
	{CloudWatchName: "BytesDownloaded", Metrics: []metric{
		{MackerelName: "BytesDownloaded", Type: metricsTypeSum},
	}},
	{CloudWatchName: "BytesUploaded", Metrics: []metric{
		{MackerelName: "BytesUploaded", Type: metricsTypeSum},
	}},
	{CloudWatchName: "FirstByteLatency", Metrics: []metric{
		{MackerelName: "FirstByteLatencyAvg", Type: metricsTypeAverage},
		{MackerelName: "FirstByteLatencyMax", Type: metricsTypeMaximum},
		{MackerelName: "FirstByteLatencyMin", Type: metricsTypeMinimum},
	}},
	{CloudWatchName: "TotalRequestLatency", Metrics: []metric{
		{MackerelName: "TotalRequestLatencyAvg", Type: metricsTypeAverage},
		{MackerelName: "TotalRequestLatencyMax", Type: metricsTypeMaximum},
		{MackerelName: "TotalRequestLatencyMin", Type: metricsTypeMinimum},
	}},
}

// FetchMetrics fetch the metrics
func (p S3RequestsPlugin) FetchMetrics() (map[string]float64, error) {
	stats := make(map[string]float64)

	for _, met := range s3RequestMetricsGroup {
		v, err := getLastPointFromCloudWatch(p.CloudWatch, p.BucketName, p.FilterID, met)
		if err != nil {
			log.Printf("%s: %s", met, err)
		} else if v != nil {
			stats = mergeStatsFromDatapoint(stats, v, met)
		}
	}
	return stats, nil
}

// GraphDefinition of S3RequestsPlugin
func (p S3RequestsPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := p.MetricLabelPrefix()

	graphdef := map[string]mp.Graphs{
		"requests": {
			Label: (labelPrefix + " Requests"),
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "AllRequests", Label: "All", Stacked: false},
				{Name: "GetRequests", Label: "Get", Stacked: true},
				{Name: "PutRequests", Label: "Put", Stacked: true},
				{Name: "DeleteRequests", Label: "Delete", Stacked: true},
				{Name: "HeadRequests", Label: "Head", Stacked: true},
				{Name: "PostRequests", Label: "Post", Stacked: true},
				{Name: "ListRequests", Label: "List", Stacked: true},
			},
		},
		"errors": {
			Label: (labelPrefix + " Errors"),
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "4xxErrors", Label: "4xx"},
				{Name: "5xxErrors", Label: "5xx"},
			},
		},
		"bytes": {
			Label: (labelPrefix + " Bytes"),
			Unit:  mp.UnitBytes,
			Metrics: []mp.Metrics{
				{Name: "BytesDownloaded", Label: "Downloaded"},
				{Name: "BytesUploaded", Label: "Uploaded"},
			},
		},
		"first_byte_latency": {
			Label: (labelPrefix + " FirstByteLatency [ms]"),
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "FirstByteLatencyAvg", Label: "Average"},
				{Name: "FirstByteLatencyMax", Label: "Maximum"},
				{Name: "FirstByteLatencyMin", Label: "Minimum"},
			},
		},
		"total_request_latency": {
			Label: (labelPrefix + " TotalRequestLatency [ms]"),
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "TotalRequestLatencyAvg", Label: "Average"},
				{Name: "TotalRequestLatencyMax", Label: "Maximum"},
				{Name: "TotalRequestLatencyMin", Label: "Minimum"},
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
	optBucketName := flag.String("bucket-name", "", "S3 bucket Name")
	optFilterID := flag.String("filter-id", "", "S3 FilterId in metrics configuration")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optKeyPrefix := flag.String("metric-key-prefix", "s3-requests", "Metric key prefix")
	optLabelPrefix := flag.String("metric-label-prefix", "S3", "Metric label prefix")
	flag.Parse()

	var plugin S3RequestsPlugin

	plugin.AccessKeyID = *optAccessKeyID
	plugin.SecretAccessKey = *optSecretAccessKey
	plugin.Region = *optRegion

	plugin.BucketName = *optBucketName
	plugin.FilterID = *optFilterID
	plugin.KeyPrefix = *optKeyPrefix
	plugin.LabelPrefix = *optLabelPrefix

	err := plugin.prepare()
	if err != nil {
		log.Fatalln(err)
	}

	helper := mp.NewMackerelPlugin(plugin)
	helper.Tempfile = *optTempfile

	helper.Run()
}
