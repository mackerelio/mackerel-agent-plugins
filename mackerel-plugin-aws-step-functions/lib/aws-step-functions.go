package mpawsstepfunctions

import (
	"flag"
	"time"

	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"log"
)

const (
	namespace          = "AWS/States"
	metricsTypeAverage = "Average"
	metricsTypeSum     = "Sum"
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

// StepFunctionsPlugin is a mackerel plugin
type StepFunctionsPlugin struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	StateMachineArn string
	Prefix          string

	CloudWatch *cloudwatch.CloudWatch
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p StepFunctionsPlugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		return "step-functions"
	}
	return p.Prefix
}

func mergeStatsFromDatapoint(stats map[string]float64, dp *cloudwatch.Datapoint, mg metricsGroup) map[string]float64 {
	for _, met := range mg.Metrics {
		switch met.Type {
		case metricsTypeAverage:
			stats[met.MackerelName] = *dp.Average
		case metricsTypeSum:
			stats[met.MackerelName] = *dp.Sum
		}
	}
	return stats
}

var stepFunctionsMetricsGroup = []metricsGroup{
	{CloudWatchName: "ExecutionTime", Metrics: []metric{
		{MackerelName: "ExecutionTime", Type: metricsTypeAverage},
	}},
	{CloudWatchName: "ExecutionThrottled", Metrics: []metric{
		{MackerelName: "ExecutionThrottled", Type: metricsTypeSum},
	}},
	{CloudWatchName: "ExecutionsAborted", Metrics: []metric{
		{MackerelName: "ExecutionsAborted", Type: metricsTypeSum},
	}},
	{CloudWatchName: "ExecutionsFailed", Metrics: []metric{
		{MackerelName: "ExecutionsFailed", Type: metricsTypeSum},
	}},
	{CloudWatchName: "ExecutionsStarted", Metrics: []metric{
		{MackerelName: "ExecutionsStarted", Type: metricsTypeSum},
	}},
	{CloudWatchName: "ExecutionsSucceeded", Metrics: []metric{
		{MackerelName: "ExecutionsSucceeded", Type: metricsTypeSum},
	}},
	{CloudWatchName: "ExecutionsTimedOut", Metrics: []metric{
		{MackerelName: "ExecutionsTimedOut", Type: metricsTypeSum},
	}},
}

// prepare creates CloudWatch instance
func (p *StepFunctionsPlugin) prepare() error {
	// validate params
	if p.StateMachineArn == "" {
		return errors.New("-state-machine-arn is necessary")
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

func getLastPointFromCloudWatch(cw *cloudwatch.CloudWatch, stateMachineArn string, metric metricsGroup) (*cloudwatch.Datapoint, error) {
	statsInput := make([]*string, len(metric.Metrics))
	for i, typ := range metric.Metrics {
		statsInput[i] = aws.String(typ.Type)
	}
	now := time.Now()
	prev := now.Add(time.Duration(180) * time.Second * -1)

	input := &cloudwatch.GetMetricStatisticsInput{
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("StateMachineArn"),
				Value: aws.String(stateMachineArn),
			},
		},
		EndTime:    aws.Time(now),
		StartTime:  aws.Time(prev),
		MetricName: aws.String(metric.CloudWatchName),
		Period:     aws.Int64(60),
		Statistics: statsInput,
		Namespace:  aws.String(namespace),
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

// FetchMetrics fetch the metrics
func (p StepFunctionsPlugin) FetchMetrics() (map[string]float64, error) {
	stats := make(map[string]float64)

	for _, met := range stepFunctionsMetricsGroup {
		v, err := getLastPointFromCloudWatch(p.CloudWatch, p.StateMachineArn, met)
		if err != nil {
			log.Printf("%s: %s", met, err)
		} else if v != nil {
			stats = mergeStatsFromDatapoint(stats, v, met)
		}
	}
	return stats, nil
}

// GraphDefinition for plugin
func (p StepFunctionsPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := p.MetricKeyPrefix()

	return map[string]mp.Graphs{
		"Executions": {
			Label: labelPrefix + " Executions",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "ExecutionsAborted", Label: "Aborted"},
				{Name: "ExecutionsFailed", Label: "Failed"},
				{Name: "ExecutionsStarted", Label: "Started"},
				{Name: "ExecutionsSucceeded", Label: "Succeeded"},
				{Name: "ExecutionsTimedOut", Label: "TimedOut"},
			},
		},
		"ExecutionThrottled": {
			Label: labelPrefix + " ExecutionThrottled",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "ExecutionThrottled", Label: "ExecutionThrottled"},
			},
		},
		"ExecutionTime": {
			Label: labelPrefix + " ExecutionTime",
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "ExecutionTime", Label: "ExecutionTime"},
			},
		},
	}
}

// Do the plugin
func Do() {
	optRegion := flag.String("region", "", "AWS Region")
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optStateMachineArn := flag.String("state-machine-arn", "", "AWS Step Functions State Machine Arn")
	optPrefix := flag.String("metric-key-prefix", "step-functions", "Metric key prefix")
	flag.Parse()

	var plugin StepFunctionsPlugin
	plugin.AccessKeyID = *optAccessKeyID
	plugin.SecretAccessKey = *optSecretAccessKey
	plugin.Region = *optRegion
	plugin.StateMachineArn = *optStateMachineArn
	plugin.Prefix = *optPrefix

	err := plugin.prepare()
	if err != nil {
		log.Fatalln(err)
	}

	mp.NewMackerelPlugin(plugin).Run()
}
