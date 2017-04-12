package mpawslambda

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/stretchr/testify/assert"
)

func defaultLambda() LambdaPlugin {
	var defaultLambda LambdaPlugin
	defaultLambda.Prefix = "lambda"
	return defaultLambda
}

func ExampleLambdaPlugin_GraphDefinition() {
	helperForDefault := mp.NewMackerelPlugin(defaultLambda())
	helperForDefault.OutputDefinitions()

	// Output:
	// # mackerel-agent-plugin
	// {"graphs":{"lambda.dead_letters":{"label":"Lambda Dead Letter","unit":"integer","metrics":[{"name":"dead_letter_errors","label":"Errors","stacked":false}]},"lambda.duration":{"label":"Lambda Duration","unit":"float","metrics":[{"name":"duration_avg","label":"Average","stacked":false},{"name":"duration_max","label":"Maximum","stacked":false},{"name":"duration_min","label":"Minimum","stacked":false}]},"lambda.invocations":{"label":"Lambda Invocations","unit":"integer","metrics":[{"name":"invocations_success","label":"Success","stacked":false},{"name":"invocations_error","label":"Error","stacked":false},{"name":"invocations_throttles","label":"Throttles","stacked":false}]}}}
}

func TestPrepare(t *testing.T) {
	p2 := defaultLambda()
	p2.Region = "MySuperRegion"
	p2.prepare()
	assert.Equal(t, "MySuperRegion", *p2.CloudWatch.Config.Region, "Specified region is used")

	// XXX Maybe we should test around AccesKeyID?
}

func TestTransformMetrics(t *testing.T) {
	regularStats := map[string]interface{}{
		"invocations_total": 150.0,
		"invocations_error": 30.0,
		"durations_avg":     250.3,
	}
	assert.Equal(t,
		map[string]interface{}{
			"invocations_success": 120.0,
			"invocations_error":   30.0,
			"durations_avg":       250.3,
		},
		transformMetrics(regularStats),
		"On regular cases values are transformed properly")

	noInvokeStats := map[string]interface{}{
		"durations_avg": 250.3,
	}
	assert.Equal(t,
		map[string]interface{}{
			"durations_avg": 250.3,
		},
		transformMetrics(noInvokeStats),
		"nothing happens when invocations_success is not present")

	// I don't know this case may happen in practice, but anyway I test.
	nonErrorStats := map[string]interface{}{
		"invocations_total": 150.0,
		"durations_avg":     250.3,
	}
	assert.Equal(t,
		map[string]interface{}{
			"invocations_success": 150.0,
			"durations_avg":       250.3,
		},
		transformMetrics(nonErrorStats),
		"Success will be calculated even if invocations_error is not present")
}

type mockCloudWatchClient struct {
	cloudwatchiface.CloudWatchAPI
	RequestedCount int
}

func (m *mockCloudWatchClient) GetMetricStatistics(input *cloudwatch.GetMetricStatisticsInput) (*cloudwatch.GetMetricStatisticsOutput, error) {
	m.RequestedCount++
	// Returns error unless expected payload

	// Check `Dimensions`
	expectedDimensions := []*cloudwatch.Dimension{
		{
			Name:  aws.String("FunctionName"),
			Value: aws.String("myFunction"),
		},
	}
	if input.Dimensions != nil && !assert.ObjectsAreEqual(expectedDimensions, input.Dimensions) {
		return nil, errors.New("Unexpected Dimension")
	}

	// Check `Statistics` for given `MetricName`
	var expectedStatistics []*string
	switch *input.MetricName {
	case "Duration":
		expectedStatistics = []*string{aws.String("Average"), aws.String("Maximum"), aws.String("Minimum")}
	default:
		expectedStatistics = []*string{aws.String("Sum")}
	}
	if !assert.ObjectsAreEqual(expectedStatistics, input.Statistics) {
		return nil, errors.New("Wrong Statistics")
	}

	// Construct Mock Response
	now := time.Now()
	output := new(cloudwatch.GetMetricStatisticsOutput)
	output.Label = input.MetricName
	switch *output.Label {
	case "Duration":
		output.Datapoints = []*cloudwatch.Datapoint{
			{Average: aws.Float64(30.0), Maximum: aws.Float64(50.0), Minimum: aws.Float64(10.0), Timestamp: aws.Time(now)},
			{Average: aws.Float64(25.0), Maximum: aws.Float64(45.0), Minimum: aws.Float64(5.0), Timestamp: aws.Time(now.Add(time.Duration(60) * time.Second * +1))},
			{Average: aws.Float64(35.0), Maximum: aws.Float64(55.0), Minimum: aws.Float64(15.0), Timestamp: aws.Time(now.Add(time.Duration(60) * time.Second * -1))},
		}
	default:
		output.Datapoints = []*cloudwatch.Datapoint{
			{Sum: aws.Float64(30.0), Timestamp: aws.Time(now)},
			{Sum: aws.Float64(25.0), Timestamp: aws.Time(now.Add(time.Duration(60) * time.Second * +1))},
			{Sum: aws.Float64(35.0), Timestamp: aws.Time(now.Add(time.Duration(60) * time.Second * -1))},
		}
	}
	return output, nil
}

func TestGetLastPointFromCloudWatch(t *testing.T) {
	mockCw := &mockCloudWatchClient{}

	dp0, err := getLastPointFromCloudWatch(mockCw, "myFunction",
		metricsGroup{CloudWatchName: "Throttles", Metrics: []metric{
			{MackerelName: "invocations_throttles", Type: metricsTypeSum},
		}})
	if err != nil {
		t.Errorf("getLastPointFromCloudWatch fails: %s", err)
	} else {
		assert.Equal(t,
			&cloudwatch.Datapoint{Sum: aws.Float64(25.0), Timestamp: dp0.Timestamp},
			dp0,
			"Can request Single statistics")
	}

	dp1, err := getLastPointFromCloudWatch(mockCw, "myFunction",
		metricsGroup{CloudWatchName: "Duration", Metrics: []metric{
			{MackerelName: "duration_avg", Type: metricsTypeAverage},
			{MackerelName: "duration_max", Type: metricsTypeMaximum},
			{MackerelName: "duration_min", Type: metricsTypeMinimum},
		}})
	if err != nil {
		t.Errorf("getLastPointFromCloudWatch fails: %s", err)
	} else {
		assert.Equal(t,
			&cloudwatch.Datapoint{Average: aws.Float64(25.0), Maximum: aws.Float64(45.0), Minimum: aws.Float64(5.0), Timestamp: dp1.Timestamp},
			dp1,
			"Can request multiple statistics at once")
	}

	dp2, err := getLastPointFromCloudWatch(mockCw, "",
		metricsGroup{CloudWatchName: "Throttles", Metrics: []metric{
			{MackerelName: "invocations_throttles", Type: metricsTypeSum},
		}})
	if err != nil {
		t.Errorf("getLastPointFromCloudWatch fails: %s", err)
	} else {
		assert.Equal(t,
			&cloudwatch.Datapoint{Sum: aws.Float64(25.0), Timestamp: dp2.Timestamp},
			dp2,
			"works even if function-name is empty")
	}

	assert.Equal(t, 3, mockCw.RequestedCount, "CloudWatch request is done once per call")
}

func TestMergeStatsFromDatapoint(t *testing.T) {
	stats := make(map[string]interface{})
	dp := cloudwatch.Datapoint{
		Average:   aws.Float64(25.0),
		Maximum:   aws.Float64(45.0),
		Minimum:   aws.Float64(5.0),
		Sum:       aws.Float64(500.0),
		Timestamp: aws.Time(time.Now()),
	}

	stats = mergeStatsFromDatapoint(stats,
		&dp,
		metricsGroup{CloudWatchName: "Invocations", Metrics: []metric{
			{MackerelName: "invocations_total", Type: metricsTypeSum},
		}})

	assert.Equal(t,
		map[string]interface{}{
			"invocations_total": 500.0,
		},
		stats,
		"Can merge single stat",
	)

	stats = mergeStatsFromDatapoint(stats,
		nil,
		metricsGroup{CloudWatchName: "DeadLetterErrors", Metrics: []metric{
			{MackerelName: "dead_letter_errors", Type: metricsTypeSum},
		}})

	assert.Equal(t,
		map[string]interface{}{
			"invocations_total": 500.0,
		},
		stats,
		"Does not panic on nil",
	)

	stats = mergeStatsFromDatapoint(stats,
		&dp,
		metricsGroup{CloudWatchName: "Duration", Metrics: []metric{
			{MackerelName: "duration_avg", Type: metricsTypeAverage},
			{MackerelName: "duration_max", Type: metricsTypeMaximum},
			{MackerelName: "duration_min", Type: metricsTypeMinimum},
		}})

	assert.Equal(t,
		map[string]interface{}{
			"invocations_total": 500.0,
			"duration_avg":      25.0,
			"duration_max":      45.0,
			"duration_min":      5.0,
		},
		stats,
		"Can merge already existing stats / can merge multiple stats at once",
	)
}
