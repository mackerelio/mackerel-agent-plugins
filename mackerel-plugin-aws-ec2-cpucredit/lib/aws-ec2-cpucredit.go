package mpawsec2cpucredit

import (
	"errors"
	"flag"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

var graphdef = map[string]mp.Graphs{
	"ec2.cpucredit": {
		Label: "EC2 CPU Credit",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "usage", Label: "Usage", Diff: false},
			{Name: "balance", Label: "Balance", Diff: false},
			{Name: "surplus_balance", Label: "Surplus Usage", Diff: false},
			{Name: "surplus_charged", Label: "Surplus Charged", Diff: false},
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

	input := &cloudwatch.GetMetricStatisticsInput{
		Dimensions: []*cloudwatch.Dimension{dimension},
		EndTime:    aws.Time(now),
		StartTime:  aws.Time(prev),
		MetricName: aws.String(metricName),
		Period:     aws.Int64(60),
		Statistics: []*string{aws.String("Average")},
		Namespace:  aws.String(namespace),
	}

	response, err := cw.GetMetricStatistics(input)
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
		latestVal = *dp.Average
	}

	return latestVal, nil
}

// FetchMetrics fetch the metrics
func (p CPUCreditPlugin) FetchMetrics() (map[string]float64, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	config := aws.NewConfig()
	if p.AccessKeyID != "" && p.SecretAccessKey != "" {
		config = config.WithCredentials(credentials.NewStaticCredentials(p.AccessKeyID, p.SecretAccessKey, ""))
	}
	if p.Region != "" {
		config = config.WithRegion(p.Region)
	}

	cw := cloudwatch.New(sess, config)

	dimension := &cloudwatch.Dimension{
		Name:  aws.String("InstanceId"),
		Value: aws.String(p.InstanceID),
	}

	stat := make(map[string]float64)

	stat["usage"], err = getLastPointAverage(cw, dimension, "CPUCreditUsage")
	if err != nil {
		return nil, err
	}

	stat["balance"], err = getLastPointAverage(cw, dimension, "CPUCreditBalance")
	if err != nil {
		return nil, err
	}

	stat["surplus_balance"], err = getLastPointAverage(cw, dimension, "CPUSurplusCreditBalance")
	if err != nil {
		return nil, err
	}

	stat["surplus_charged"], err = getLastPointAverage(cw, dimension, "CPUSurplusCreditsCharged")
	if err != nil {
		return nil, err
	}

	return stat, nil
}

// GraphDefinition for plugin
func (p CPUCreditPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// Do the plugin
func Do() {
	optRegion := flag.String("region", "", "AWS Region")
	optInstanceID := flag.String("instance-id", "", "Instance ID")
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var cpucredit CPUCreditPlugin

	if *optRegion == "" || *optInstanceID == "" {
		sess, err := session.NewSession()
		if err != nil {
			log.Fatalln(err)
		}
		ec2metadata := ec2metadata.New(sess)
		if ec2metadata.Available() {
			cpucredit.Region, _ = ec2metadata.Region()
			cpucredit.InstanceID, _ = ec2metadata.GetMetadata("instance-id")
		}
	} else {
		cpucredit.Region = *optRegion
		cpucredit.InstanceID = *optInstanceID
	}

	cpucredit.AccessKeyID = *optAccessKeyID
	cpucredit.SecretAccessKey = *optSecretAccessKey

	helper := mp.NewMackerelPlugin(cpucredit)
	helper.Tempfile = *optTempfile

	helper.Run()
}
