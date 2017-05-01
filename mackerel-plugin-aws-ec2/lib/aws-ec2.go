package mpawsec2

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

// EC2Plugin is a mackerel plugin for ec2
type EC2Plugin struct {
	InstanceID  string
	Region      string
	Credentials *credentials.Credentials
	CloudWatch  *cloudwatch.CloudWatch
}

var graphdef = map[string]mp.Graphs{
	"ec2.CPUUtilization": {
		Label: "CPU Utilization",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "CPUUtilization", Label: "CPUUtilization"},
		},
	},
	"ec2.DiskBytes": {
		Label: "Disk Bytes",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "DiskReadBytes", Label: "DiskReadBytes"},
			{Name: "DiskWriteBytes", Label: "DiskWriteBytes"},
		},
	},
	"ec2.DiskOps": {
		Label: "Disk Ops",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "DiskReadOps", Label: "DiskReadOps"},
			{Name: "DiskWriteOps", Label: "DiskWriteOps"},
		},
	},
	"ec2.Network": {
		Label: "Network",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "NetworkIn", Label: "NetworkIn"},
			{Name: "NetworkOut", Label: "NetworkOut"},
		},
	},
	"ec2.NetworkPackets": {
		Label: "Network",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "NetworkPacketsIn", Label: "NetworkPacketsIn"},
			{Name: "NetworkPacketsOut", Label: "NetworkPacketsOut"},
		},
	},
	"ec2.StatusCheckFailed": {
		Label: "StatusCheckFailed",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "StatusCheckFailed", Label: "StatusCheckFailed"},
			{Name: "StatusCheckFailed_Instance", Label: "StatusCheckFailed_Instance"},
			{Name: "StatusCheckFailed_System", Label: "StatusCheckFailed_System"},
		},
	},
}

// GraphDefinition returns graphdef
func (p EC2Plugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

func getLastPoint(cloudWatch *cloudwatch.CloudWatch, dimension *cloudwatch.Dimension, metricName string) (float64, error) {
	now := time.Now()
	response, err := cloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		StartTime:  aws.Time(now.Add(time.Duration(600) * time.Second * -1)),
		EndTime:    aws.Time(now),
		Period:     aws.Int64(60),
		Namespace:  aws.String("AWS/EC2"),
		Dimensions: []*cloudwatch.Dimension{dimension},
		MetricName: aws.String(metricName),
		Statistics: []*string{aws.String("Average")},
	})
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

// FetchMetrics fetches metrics from CloudWatch
func (p EC2Plugin) FetchMetrics() (map[string]float64, error) {
	stat := make(map[string]float64)
	p.CloudWatch = cloudwatch.New(session.New(
		&aws.Config{
			Credentials: p.Credentials,
			Region:      &p.Region,
		}))
	dimension := &cloudwatch.Dimension{
		Name:  aws.String("InstanceId"),
		Value: aws.String(p.InstanceID),
	}

	for _, met := range [...]string{
		"CPUUtilization",
		"DiskReadBytes",
		"DiskReadOps",
		"DiskWriteBytes",
		"DiskWriteOps",
		"NetworkIn",
		"NetworkOut",
		"NetworkPacketsIn",
		"NetworkPacketsOut",
		"StatusCheckFailed",
		"StatusCheckFailed_Instance",
		"StatusCheckFailed_System",
	} {
		v, err := getLastPoint(p.CloudWatch, dimension, met)
		if err == nil {
			stat[met] = v
		} else {
			log.Printf("%s: %s", met, err)
		}
	}

	return stat, nil
}

// Do the plugin
func Do() {
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optRegion := flag.String("region", "", "AWS Region")
	optInstanceID := flag.String("instance-id", "", "Instance ID")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var ec2 EC2Plugin

	// use credentials from option
	if *optAccessKeyID != "" && *optSecretAccessKey != "" {
		ec2.Credentials = credentials.NewStaticCredentials(*optAccessKeyID, *optSecretAccessKey, "")
	}

	// get metadata in ec2 instance
	ec2MC := ec2metadata.New(session.New())
	ec2.Region = *optRegion
	if *optRegion == "" {
		ec2.Region, _ = ec2MC.Region()
	}
	ec2.InstanceID = *optInstanceID
	if *optInstanceID == "" {
		ec2.InstanceID, _ = ec2MC.GetMetadata("instance-id")
	}

	helper := mp.NewMackerelPlugin(ec2)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.SetTempfileByBasename(fmt.Sprintf("mackerel-plugin-aws-ec2-%s", ec2.InstanceID))
	}

	helper.Run()
}
