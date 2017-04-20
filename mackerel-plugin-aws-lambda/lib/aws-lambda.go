package mpawslambda

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
	namespace          = "AWS/Lambda"
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

// LambdaPlugin mackerel plugin for aws Lambda
type LambdaPlugin struct {
	FunctionName string
	Prefix       string

	AccessKeyID     string
	SecretAccessKey string
	Region          string

	CloudWatch *cloudwatch.CloudWatch
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p LambdaPlugin) MetricKeyPrefix() string {
	return p.Prefix
}

// prepare creates CloudWatch instance
func (p *LambdaPlugin) prepare() error {

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
func getLastPointFromCloudWatch(cw cloudwatchiface.CloudWatchAPI, functionName string, metric metricsGroup) (*cloudwatch.Datapoint, error) {
	now := time.Now()
	statsInput := make([]*string, len(metric.Metrics))
	for i, typ := range metric.Metrics {
		statsInput[i] = aws.String(typ.Type)
	}
	input := &cloudwatch.GetMetricStatisticsInput{
		// Usually Cloudwatch datapoints delays about 2 mins, so retrieve last 3 mins (with 1 min buffer)
		StartTime:  aws.Time(now.Add(time.Duration(180) * time.Second * -1)),
		EndTime:    aws.Time(now),
		MetricName: aws.String(metric.CloudWatchName),
		Period:     aws.Int64(60),
		Statistics: statsInput,
		Namespace:  aws.String(namespace),
	}
	if functionName != "" {
		input.Dimensions = []*cloudwatch.Dimension{{
			Name:  aws.String("FunctionName"),
			Value: aws.String(functionName),
		}}
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

// TransformMetrics converts some of datapoints to post differences of two metrics
func transformMetrics(stats map[string]interface{}) map[string]interface{} {
	// Although stats are interface{}, those values from cloudwatch.Datapoint are guaranteed to be float64.
	if totalCount, ok := stats["invocations_total"].(float64); ok {
		if errorCount, ok := stats["invocations_error"].(float64); ok {
			stats["invocations_success"] = totalCount - errorCount
		} else {
			stats["invocations_success"] = totalCount
		}
		delete(stats, "invocations_total")
	}
	return stats
}

func mergeStatsFromDatapoint(stats map[string]interface{}, dp *cloudwatch.Datapoint, mg metricsGroup) map[string]interface{} {
	if dp != nil {
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
	}
	return stats
}

var lambdaMetricsGroup = []metricsGroup{
	{CloudWatchName: "Invocations", Metrics: []metric{
		{MackerelName: "invocations_total", Type: metricsTypeSum},
	}},
	{CloudWatchName: "Errors", Metrics: []metric{
		{MackerelName: "invocations_error", Type: metricsTypeSum},
	}},
	{CloudWatchName: "DeadLetterErrors", Metrics: []metric{
		{MackerelName: "dead_letter_errors", Type: metricsTypeSum},
	}},
	{CloudWatchName: "Throttles", Metrics: []metric{
		{MackerelName: "invocations_throttles", Type: metricsTypeSum},
	}},
	{CloudWatchName: "Duration", Metrics: []metric{
		{MackerelName: "duration_avg", Type: metricsTypeAverage},
		{MackerelName: "duration_max", Type: metricsTypeMaximum},
		{MackerelName: "duration_min", Type: metricsTypeMinimum},
	}},
}

// FetchMetrics fetch the metrics
func (p LambdaPlugin) FetchMetrics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	for _, met := range lambdaMetricsGroup {
		v, err := getLastPointFromCloudWatch(p.CloudWatch, p.FunctionName, met)
		if err == nil {
			stats = mergeStatsFromDatapoint(stats, v, met)
		} else {
			log.Printf("%s: %s", met, err)
		}
	}
	return transformMetrics(stats), nil
}

// GraphDefinition of LambdaPlugin
func (p LambdaPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(p.Prefix)

	graphdef := map[string]mp.Graphs{
		"invocations": {
			Label: (labelPrefix + " Invocations"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "invocations_success", Label: "Success"},
				{Name: "invocations_error", Label: "Error"},
				{Name: "invocations_throttles", Label: "Throttles"},
			},
		},
		"dead_letters": {
			Label: (labelPrefix + " Dead Letter"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "dead_letter_errors", Label: "Errors"},
			},
		},
		"duration": {
			Label: (labelPrefix + " Duration"),
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "duration_avg", Label: "Average"},
				{Name: "duration_max", Label: "Maximum"},
				{Name: "duration_min", Label: "Minimum"},
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
	optFunctionName := flag.String("function-name", "", "Function Name")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optPrefix := flag.String("metric-key-prefix", "lambda", "Metric key prefix")
	flag.Parse()

	var plugin LambdaPlugin

	plugin.AccessKeyID = *optAccessKeyID
	plugin.SecretAccessKey = *optSecretAccessKey
	plugin.Region = *optRegion

	plugin.FunctionName = *optFunctionName
	plugin.Prefix = *optPrefix

	err := plugin.prepare()
	if err != nil {
		log.Fatalln(err)
	}

	helper := mp.NewMackerelPlugin(plugin)
	helper.Tempfile = *optTempfile

	helper.Run()
}
