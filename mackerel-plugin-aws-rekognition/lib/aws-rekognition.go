package mpawsrekognition

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
	"rekognition.RequestCount": {
		Label: "AWS Rekognition RequestCount",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "SuccessfulRequestCount", Label: "SuccessfulRequestCount"},
			{Name: "ThrottledCount", Label: "ThrottledCount"},
			{Name: "ServerErrorCount", Label: "ServerErrorCount"},
			{Name: "UserErrorCount", Label: "UserErrorCount"},
		},
	},
	"rekognition.ResponseTime": {
		Label: "AWS Rekognition ResponseTime",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "ResponseTime", Label: "ResponseTime"},
		},
	},
	"rekognition.DetectedCount": {
		Label: "AWS Rekognition DetectedCount",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "DetectedFaceCount", Label: "DetectedFaceCount"},
			{Name: "DetectedLabelCount", Label: "DetectedLabelCount"},
		},
	},
}

// RekognitionPlugin mackerel plugin for aws rekognition
type RekognitionPlugin struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Operation       string
	CloudWatch      *cloudwatch.CloudWatch
}

func (p *RekognitionPlugin) prepare() error {
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

func (p RekognitionPlugin) getLastPoint(metricName string) (float64, error) {
	now := time.Now()

	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Rekognition"),
		MetricName: aws.String(metricName),
		StartTime:  aws.Time(now.Add(time.Duration(180) * time.Second * -1)),
		EndTime:    aws.Time(now),
		Period:     aws.Int64(60),
		Statistics: []*string{aws.String("Average")},
	}
	if p.Operation != "" {
		input.Dimensions = []*cloudwatch.Dimension{{
			Name:  aws.String("Operation"),
			Value: aws.String(p.Operation),
		}}
	}

	response, err := p.CloudWatch.GetMetricStatistics(input)

	if err != nil {
		return 0, err
	}

	datapoints := response.Datapoints
	if len(datapoints) == 0 {
		return 0, errors.New("fetched no datapoints")
	}

	latest := time.Unix(0, 0)
	var latestVal float64
	for _, dp := range datapoints {
		if dp.Timestamp.Before(latest) {
			continue
		}

		latest = *dp.Timestamp
		latestVal = *dp.Average
	}

	return latestVal, nil
}

// FetchMetrics interface for mackerelplugin
func (p RekognitionPlugin) FetchMetrics() (map[string]float64, error) {
	stat := make(map[string]float64)

	for _, met := range [...]string{
		"SuccessfulRequestCount",
		"ThrottledCount",
		"ResponseTime",
		"DetectedFaceCount",
		"DetectedLabelCount",
		"ServerErrorCount",
		"UserErrorCount",
	} {
		v, err := p.getLastPoint(met)
		if err == nil {
			stat[met] = v
		} else {
			log.Printf("%s: %s", met, err)
		}
	}

	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (p RekognitionPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// Do the plugin
func Do() {
	optRegion := flag.String("region", "", "AWS Region")
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optOperation := flag.String("operation", "", "AWS Rekognition Operation")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var rekognition RekognitionPlugin

	if *optRegion == "" {
		ec2metadata := ec2metadata.New(session.New())
		if ec2metadata.Available() {
			rekognition.Region, _ = ec2metadata.Region()
		}
	} else {
		rekognition.Region = *optRegion
	}

	rekognition.Region = *optRegion
	rekognition.AccessKeyID = *optAccessKeyID
	rekognition.SecretAccessKey = *optSecretAccessKey
	rekognition.Operation = *optOperation

	err := rekognition.prepare()
	if err != nil {
		log.Fatalln(err)
	}

	helper := mp.NewMackerelPlugin(rekognition)
	helper.Tempfile = *optTempfile
	helper.Run()
}
