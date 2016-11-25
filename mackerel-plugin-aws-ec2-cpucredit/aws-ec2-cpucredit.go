package main

import (
	"errors"
	"flag"
	"os"
	"time"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

var graphdef = map[string]mp.Graphs{
	"ec2.cpucredit": {
		Label: "EC2 CPU Credit",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "usage", Label: "Usage", Diff: false},
			{Name: "balance", Label: "Balance", Diff: false},
		},
	},
}

// CPUCreditPlugin is a mackerel plugin
type CPUCreditPlugin struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	InstanceID      string
}

func getLastPointAverage(cw *cloudwatch.CloudWatch, dimension *cloudwatch.Dimension, metricName string) (float64, error) {
	namespace := "AWS/EC2"
	now := time.Now()
	prev := now.Add(time.Duration(600) * time.Second * -1) // 10 min (to fetch at least 1 data-point)

	request := &cloudwatch.GetMetricStatisticsRequest{
		Dimensions: []cloudwatch.Dimension{*dimension},
		EndTime:    now,
		StartTime:  prev,
		MetricName: metricName,
		Period:     60,
		Statistics: []string{"Average"},
		Namespace:  namespace,
	}

	response, err := cw.GetMetricStatistics(request)
	if err != nil {
		return 0, err
	}

	datapoints := response.GetMetricStatisticsResult.Datapoints
	if len(datapoints) == 0 {
		return 0, errors.New("fetched no datapoints")
	}

	latest := time.Unix(0, 0)
	var latestVal float64
	for _, dp := range datapoints {
		if dp.Timestamp.Before(latest) {
			continue
		}

		latest = dp.Timestamp
		latestVal = dp.Average
	}

	return latestVal, nil
}

// FetchMetrics fetch the metrics
func (p CPUCreditPlugin) FetchMetrics() (map[string]float64, error) {
	region := aws.Regions[p.Region]
	dimension := &cloudwatch.Dimension{
		Name:  "InstanceId",
		Value: p.InstanceID,
	}

	auth, err := aws.GetAuth(p.AccessKeyID, p.SecretAccessKey, "", time.Now())
	if err != nil {
		return nil, err
	}
	cw, err := cloudwatch.NewCloudWatch(auth, region.CloudWatchServicepoint)

	stat := make(map[string]float64)

	stat["usage"], err = getLastPointAverage(cw, dimension, "CPUCreditUsage")
	if err != nil {
		return nil, err
	}

	stat["balance"], err = getLastPointAverage(cw, dimension, "CPUCreditBalance")
	if err != nil {
		return nil, err
	}

	return stat, nil
}

// GraphDefinition for plugin
func (p CPUCreditPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

func main() {
	optRegion := flag.String("region", "", "AWS Region")
	optInstanceID := flag.String("instance-id", "", "Instance ID")
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var cpucredit CPUCreditPlugin

	if *optRegion == "" || *optInstanceID == "" {
		cpucredit.Region = aws.InstanceRegion()
		cpucredit.InstanceID = aws.InstanceId()
	} else {
		cpucredit.Region = *optRegion
		cpucredit.InstanceID = *optInstanceID
	}

	cpucredit.AccessKeyID = *optAccessKeyID
	cpucredit.SecretAccessKey = *optSecretAccessKey

	helper := mp.NewMackerelPlugin(cpucredit)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = "/tmp/mackerel-plugin-cpucredit"
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
