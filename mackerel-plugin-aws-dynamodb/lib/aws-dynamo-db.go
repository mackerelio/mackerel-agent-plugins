package mpawsdynamodb

import (
	"flag"
	"log"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

const (
	namespace              = "AWS/DynamoDB"
	metricsTypeAverage     = "Average"
	metricsTypeSum         = "Sum"
	metricsTypeMaximum     = "Maximum"
	metricsTypeMinimum     = "Minimum"
	metricsTypeSampleCount = "SampleCount"
)

// has 1 CloudWatch MetricName and corresponding N Mackerel Metrics
type metricsGroup struct {
	CloudWatchName string
	Metrics        []metric
}

type metric struct {
	MackerelName string
	Type         string
	FillZero     bool
}

// DynamoDBPlugin mackerel plugin for aws kinesis
type DynamoDBPlugin struct {
	TableName string
	Prefix    string

	AccessKeyID     string
	SecretAccessKey string
	Region          string
	CloudWatch      *cloudwatch.CloudWatch
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p DynamoDBPlugin) MetricKeyPrefix() string {
	return p.Prefix
}

// prepare creates CloudWatch instance
func (p *DynamoDBPlugin) prepare() error {
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

func transformAndAppendDatapoint(stats map[string]interface{}, dp *cloudwatch.Datapoint, dataType string, label string, fillZero bool) map[string]interface{} {
	if dp != nil {
		switch dataType {
		case metricsTypeAverage:
			stats[label] = *dp.Average
		case metricsTypeSum:
			stats[label] = *dp.Sum
		case metricsTypeMaximum:
			stats[label] = *dp.Maximum
		case metricsTypeMinimum:
			stats[label] = *dp.Minimum
		case metricsTypeSampleCount:
			stats[label] = *dp.SampleCount
		}
	} else if fillZero {
		stats[label] = 0.0
	}
	return stats
}

// fetch metrics which takes "Operation" dimensions querying both ListMetrics and GetMetricsStatistics
func fetchOperationWildcardMetrics(cw cloudwatchiface.CloudWatchAPI, mg metricsGroup, baseDimensions []*cloudwatch.Dimension) (map[string]interface{}, error) {
	// get available dimensions
	dimensionFilters := make([]*cloudwatch.DimensionFilter, len(baseDimensions))
	for i, dimension := range baseDimensions {
		dimensionFilters[i] = &cloudwatch.DimensionFilter{
			Name:  dimension.Name,
			Value: dimension.Value,
		}
	}
	input := &cloudwatch.ListMetricsInput{
		Dimensions: dimensionFilters,
		Namespace:  aws.String(namespace),
		MetricName: aws.String(mg.CloudWatchName),
	}
	// ListMetrics can retrieve up to 500 metrics, but DynamoDB Operations are apparently less than 500
	res, err := cw.ListMetrics(input)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]interface{})

	// get datapoints with retrieved dimensions
	for _, cwMetric := range res.Metrics {
		dimensions := cwMetric.Dimensions
		// extract operation name
		var operation *string
		for _, d := range dimensions {
			if *d.Name == "Operation" {
				operation = d.Value
				break
			}
		}
		if operation == nil {
			log.Printf("Unexpected dimension, skip: %s", dimensions)
			continue
		}

		dp, err := getLastPointFromCloudWatch(cw, mg, dimensions)
		if err != nil {
			return nil, nil
		}
		if dp != nil {
			for _, met := range mg.Metrics {
				label := strings.Replace(met.MackerelName, "#", *operation, 1)
				stats = transformAndAppendDatapoint(stats, dp, met.Type, label, met.FillZero)
			}
		}
	}

	return stats, nil
}

// getLastPoint fetches a CloudWatch metric and parse
func getLastPointFromCloudWatch(cw cloudwatchiface.CloudWatchAPI, metric metricsGroup, dimensions []*cloudwatch.Dimension) (*cloudwatch.Datapoint, error) {
	now := time.Now()
	statsInput := make([]*string, len(metric.Metrics))
	for i, typ := range metric.Metrics {
		statsInput[i] = aws.String(typ.Type)
	}
	input := &cloudwatch.GetMetricStatisticsInput{
		// 8 min, since some metrics are aggregated over 5 min
		StartTime:  aws.Time(now.Add(time.Duration(480) * time.Second * -1)),
		EndTime:    aws.Time(now),
		MetricName: aws.String(metric.CloudWatchName),
		Period:     aws.Int64(60),
		Statistics: statsInput,
		Namespace:  aws.String(namespace),
		Dimensions: dimensions,
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

var defaultMetricsGroup = []metricsGroup{
	{CloudWatchName: "ConditionalCheckFailedRequests", Metrics: []metric{
		{MackerelName: "ConditionalCheckFailedRequests", Type: metricsTypeSum},
	}},
	{CloudWatchName: "ConsumedReadCapacityUnits", Metrics: []metric{
		{MackerelName: "ConsumedReadCapacityUnitsSum", Type: metricsTypeSum},
		{MackerelName: "ConsumedReadCapacityUnitsAverage", Type: metricsTypeAverage},
	}},
	{CloudWatchName: "ConsumedWriteCapacityUnits", Metrics: []metric{
		{MackerelName: "ConsumedWriteCapacityUnitsSum", Type: metricsTypeSum},
		{MackerelName: "ConsumedWriteCapacityUnitsAverage", Type: metricsTypeAverage},
	}},
	{CloudWatchName: "ProvisionedReadCapacityUnits", Metrics: []metric{
		{MackerelName: "ProvisionedReadCapacityUnits", Type: metricsTypeMinimum},
	}},
	{CloudWatchName: "ProvisionedWriteCapacityUnits", Metrics: []metric{
		{MackerelName: "ProvisionedWriteCapacityUnits", Type: metricsTypeMinimum},
	}},
	{CloudWatchName: "SystemErrors", Metrics: []metric{
		{MackerelName: "SystemErrors", Type: metricsTypeSum},
	}},
	{CloudWatchName: "UserErrors", Metrics: []metric{
		{MackerelName: "UserErrors", Type: metricsTypeSum},
	}},
	{CloudWatchName: "ReadThrottleEvents", Metrics: []metric{
		{MackerelName: "ReadThrottleEvents", Type: metricsTypeSum, FillZero: true},
	}},
	{CloudWatchName: "WriteThrottleEvents", Metrics: []metric{
		{MackerelName: "WriteThrottleEvents", Type: metricsTypeSum, FillZero: true},
	}},
	{CloudWatchName: "TimeToLiveDeletedItemCount", Metrics: []metric{
		{MackerelName: "TimeToLiveDeletedItemCount", Type: metricsTypeSum},
	}},
}

var operationalMetricsGroup = []metricsGroup{
	{CloudWatchName: "SuccessfulRequestLatency", Metrics: []metric{
		{MackerelName: "SuccessfulRequests.#", Type: metricsTypeSampleCount},
		{MackerelName: "SuccessfulRequestLatency.#.Minimum", Type: metricsTypeMinimum},
		{MackerelName: "SuccessfulRequestLatency.#.Maximum", Type: metricsTypeMaximum},
		{MackerelName: "SuccessfulRequestLatency.#.Average", Type: metricsTypeAverage},
	}},
	{CloudWatchName: "ThrottledRequests", Metrics: []metric{
		{MackerelName: "ThrottledRequests.#", Type: metricsTypeSampleCount},
	}},
	{CloudWatchName: "SystemErrors", Metrics: []metric{
		{MackerelName: "SystemErrors.#", Type: metricsTypeSampleCount},
	}},
	{CloudWatchName: "UserErrors", Metrics: []metric{
		{MackerelName: "UserErrors.#", Type: metricsTypeSampleCount},
	}},
	{CloudWatchName: "ReturnedItemCount", Metrics: []metric{
		{MackerelName: "ReturnedItemCount.#", Type: metricsTypeAverage},
	}},
}

// FetchMetrics fetch the metrics
func (p DynamoDBPlugin) FetchMetrics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	tableDimensions := []*cloudwatch.Dimension{{
		Name:  aws.String("TableName"),
		Value: aws.String(p.TableName),
	}}
	var mu sync.Mutex
	var eg errgroup.Group
	for _, met := range defaultMetricsGroup {
		met := met
		eg.Go(func() error {
			dp, err := getLastPointFromCloudWatch(p.CloudWatch, met, tableDimensions)
			if err != nil {
				return err
			}
			for _, m := range met.Metrics {
				mu.Lock()
				stats = transformAndAppendDatapoint(stats, dp, m.Type, m.MackerelName, m.FillZero)
				mu.Unlock()
			}
			return nil
		})
	}

	for _, met := range operationalMetricsGroup {
		met := met
		eg.Go(func() error {
			operationalStats, err := fetchOperationWildcardMetrics(p.CloudWatch, met, tableDimensions)
			if err != nil {
				return err
			}

			for name, s := range operationalStats {
				mu.Lock()
				stats[name] = s
				mu.Unlock()
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return transformMetrics(stats), nil
}

// TransformMetrics converts some of datapoints to post differences of two metrics
func transformMetrics(stats map[string]interface{}) map[string]interface{} {
	// Although stats are interface{}, those values from cloudwatch.Datapoint are guaranteed to be numerical
	if consumedReadCapacitySum, ok := stats["ConsumedReadCapacityUnitsSum"].(float64); ok {
		stats["ConsumedReadCapacityUnitsNormalized"] = consumedReadCapacitySum / 60.0
	}
	if consumedWriteCapacitySum, ok := stats["ConsumedWriteCapacityUnitsSum"].(float64); ok {
		stats["ConsumedWriteCapacityUnitsNormalized"] = consumedWriteCapacitySum / 60.0
	}
	return stats
}

// GraphDefinition of DynamoDBPlugin
func (p DynamoDBPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := cases.Title(language.Und, cases.NoLower).String(p.Prefix)
	labelPrefix = strings.Replace(labelPrefix, "-", " ", -1)
	labelPrefix = strings.Replace(labelPrefix, "Dynamodb", "DynamoDB", -1)

	var graphdef = map[string]mp.Graphs{
		"ReadCapacity": {
			Label: (labelPrefix + " Read Capacity Units"),
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "ProvisionedReadCapacityUnits", Label: "Provisioned"},
				{Name: "ConsumedReadCapacityUnitsNormalized", Label: "Consumed"},
				{Name: "ConsumedReadCapacityUnitsAverage", Label: "Consumed (Average per request)"},
			},
		},
		"WriteCapacity": {
			Label: (labelPrefix + " Write Capacity Units"),
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "ProvisionedWriteCapacityUnits", Label: "Provisioned"},
				{Name: "ConsumedWriteCapacityUnitsNormalized", Label: "Consumed"},
				{Name: "ConsumedWriteCapacityUnitsAverage", Label: "Consumed (Average per request)"},
			},
		},
		"TimeToLiveDeletedItemCount": {
			Label: (labelPrefix + " TimeToLiveDeletedItemCount"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "TimeToLiveDeletedItemCount", Label: "Count"},
			},
		},
		"ThrottledEvents": {
			Label: (labelPrefix + " Throttle Events"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "ReadThrottleEvents", Label: "Read"},
				{Name: "WriteThrottleEvents", Label: "Write"},
			},
		},
		"ConditionalCheckFailedRequests": {
			Label: (labelPrefix + " ConditionalCheckFailedRequests"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "ConditionalCheckFailedRequests", Label: "Counts"},
			},
		},
		"ThrottledRequests": {
			Label: (labelPrefix + " ThrottledRequests"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "PutItem", Label: "PutItem", Stacked: true, AbsoluteName: true},
				{Name: "DeleteItem", Label: "DeleteItem", Stacked: true, AbsoluteName: true},
				{Name: "UpdateItem", Label: "UpdateItem", Stacked: true, AbsoluteName: true},
				{Name: "GetItem", Label: "GetItem", Stacked: true, AbsoluteName: true},
				{Name: "BatchGetItem", Label: "BatchGetItem", Stacked: true, AbsoluteName: true},
				{Name: "Scan", Label: "Scan", Stacked: true, AbsoluteName: true},
				{Name: "Query", Label: "Query", Stacked: true, AbsoluteName: true},
				{Name: "BatchWriteItem", Label: "BatchWriteItem", Stacked: true, AbsoluteName: true},
			},
		},
		"SystemErrors": {
			Label: (labelPrefix + " SystemErrors"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "PutItem", Label: "PutItem", Stacked: true, AbsoluteName: true},
				{Name: "DeleteItem", Label: "DeleteItem", Stacked: true, AbsoluteName: true},
				{Name: "UpdateItem", Label: "UpdateItem", Stacked: true, AbsoluteName: true},
				{Name: "GetItem", Label: "GetItem", Stacked: true, AbsoluteName: true},
				{Name: "BatchGetItem", Label: "BatchGetItem", Stacked: true, AbsoluteName: true},
				{Name: "Scan", Label: "Scan", Stacked: true, AbsoluteName: true},
				{Name: "Query", Label: "Query", Stacked: true, AbsoluteName: true},
				{Name: "BatchWriteItem", Label: "BatchWriteItem", Stacked: true, AbsoluteName: true},
			},
		},
		"UserErrors": {
			Label: (labelPrefix + " UserErrors"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "PutItem", Label: "PutItem", Stacked: true, AbsoluteName: true},
				{Name: "DeleteItem", Label: "DeleteItem", Stacked: true, AbsoluteName: true},
				{Name: "UpdateItem", Label: "UpdateItem", Stacked: true, AbsoluteName: true},
				{Name: "GetItem", Label: "GetItem", Stacked: true, AbsoluteName: true},
				{Name: "BatchGetItem", Label: "BatchGetItem", Stacked: true, AbsoluteName: true},
				{Name: "Scan", Label: "Scan", Stacked: true, AbsoluteName: true},
				{Name: "Query", Label: "Query", Stacked: true, AbsoluteName: true},
				{Name: "BatchWriteItem", Label: "BatchWriteItem", Stacked: true, AbsoluteName: true},
			},
		},
		"ReturnedItemCount": {
			Label: (labelPrefix + " ReturnedItemCount"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "PutItem", Label: "PutItem", AbsoluteName: true},
				{Name: "DeleteItem", Label: "DeleteItem", AbsoluteName: true},
				{Name: "UpdateItem", Label: "UpdateItem", AbsoluteName: true},
				{Name: "GetItem", Label: "GetItem", AbsoluteName: true},
				{Name: "BatchGetItem", Label: "BatchGetItem", AbsoluteName: true},
				{Name: "Scan", Label: "Scan", AbsoluteName: true},
				{Name: "Query", Label: "Query", AbsoluteName: true},
				{Name: "BatchWriteItem", Label: "BatchWriteItem", AbsoluteName: true},
			},
		},
		"SuccessfulRequests": {
			Label: (labelPrefix + " SuccessfulRequests"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "PutItem", Label: "PutItem", AbsoluteName: true},
				{Name: "DeleteItem", Label: "DeleteItem", AbsoluteName: true},
				{Name: "UpdateItem", Label: "UpdateItem", AbsoluteName: true},
				{Name: "GetItem", Label: "GetItem", AbsoluteName: true},
				{Name: "BatchGetItem", Label: "BatchGetItem", AbsoluteName: true},
				{Name: "Scan", Label: "Scan", AbsoluteName: true},
				{Name: "Query", Label: "Query", AbsoluteName: true},
				{Name: "BatchWriteItem", Label: "BatchWriteItem", AbsoluteName: true},
			},
		},
		"SuccessfulRequestLatency.#": {
			Label: (labelPrefix + " SuccessfulRequestLatency"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "Minimum", Label: "Min"},
				{Name: "Maximum", Label: "Max"},
				{Name: "Average", Label: "Average"},
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
	optTableName := flag.String("table-name", "", "DynamoDB Table Name")
	optPrefix := flag.String("metric-key-prefix", "dynamodb", "Metric key prefix")
	flag.Parse()

	var plugin DynamoDBPlugin

	plugin.AccessKeyID = *optAccessKeyID
	plugin.SecretAccessKey = *optSecretAccessKey
	plugin.Region = *optRegion
	plugin.TableName = *optTableName
	plugin.Prefix = *optPrefix

	err := plugin.prepare()
	if err != nil {
		log.Fatalln(err)
	}

	helper := mp.NewMackerelPlugin(plugin)

	helper.Run()
}
