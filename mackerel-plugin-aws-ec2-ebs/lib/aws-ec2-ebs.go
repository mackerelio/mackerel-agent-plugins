package mpawsec2ebs

import (
	"errors"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var metricPeriodDefault = 300
var metricPeriodByVolumeType = map[string]int{
	"io1": 60,
}

var defaultGraphs = []string{
	"ec2.ebs.bandwidth.#",
	"ec2.ebs.throughput.#",
	"ec2.ebs.size_per_op.#",
	"ec2.ebs.latency.#",
	"ec2.ebs.queue_length.#",
	"ec2.ebs.idle_time.#",
}

var allGraphs = append([]string{
	"ec2.ebs.throughput_delivered.#",
	"ec2.ebs.consumed_ops.#",
}, defaultGraphs...)

type cloudWatchSetting struct {
	MetricName string
	Statistics string
	CalcFunc   func(float64, float64) float64
}

// http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/monitoring-volume-status.html
var cloudwatchdefs = map[string](cloudWatchSetting){
	"ec2.ebs.bandwidth.#.read": cloudWatchSetting{
		MetricName: "VolumeReadBytes", Statistics: "Sum",
		CalcFunc: func(val float64, period float64) float64 { return val / period },
	},
	"ec2.ebs.bandwidth.#.write": cloudWatchSetting{
		MetricName: "VolumeWriteBytes", Statistics: "Sum",
		CalcFunc: func(val float64, period float64) float64 { return val / period },
	},
	"ec2.ebs.throughput.#.read": cloudWatchSetting{
		MetricName: "VolumeReadOps", Statistics: "Sum",
		CalcFunc: func(val float64, period float64) float64 { return val / period },
	},
	"ec2.ebs.throughput.#.write": cloudWatchSetting{
		MetricName: "VolumeWriteOps", Statistics: "Sum",
		CalcFunc: func(val float64, period float64) float64 { return val / period },
	},
	"ec2.ebs.size_per_op.#.read": cloudWatchSetting{
		MetricName: "VolumeReadBytes", Statistics: "Average",
		CalcFunc: func(val float64, period float64) float64 { return val },
	},
	"ec2.ebs.size_per_op.#.write": cloudWatchSetting{
		MetricName: "VolumeWriteBytes", Statistics: "Average",
		CalcFunc: func(val float64, period float64) float64 { return val },
	},
	"ec2.ebs.latency.#.read": cloudWatchSetting{
		MetricName: "VolumeTotalReadTime", Statistics: "Average",
		CalcFunc: func(val float64, period float64) float64 { return val * 1000 },
	},
	"ec2.ebs.latency.#.write": cloudWatchSetting{
		MetricName: "VolumeTotalWriteTime", Statistics: "Average",
		CalcFunc: func(val float64, period float64) float64 { return val * 1000 },
	},
	"ec2.ebs.queue_length.#.queue_length": cloudWatchSetting{
		MetricName: "VolumeQueueLength", Statistics: "Average",
		CalcFunc: func(val float64, period float64) float64 { return val },
	},
	"ec2.ebs.idle_time.#.idle_time": cloudWatchSetting{
		MetricName: "VolumeIdleTime", Statistics: "Sum",
		CalcFunc: func(val float64, period float64) float64 { return val / period * 100 },
	},
	"ec2.ebs.throughput_delivered.#.throughput_delivered": cloudWatchSetting{
		MetricName: "VolumeThroughputPercentage", Statistics: "Average",
		CalcFunc: func(val float64, period float64) float64 { return val },
	},
	"ec2.ebs.consumed_ops.#.consumed_ops": cloudWatchSetting{
		MetricName: "VolumeConsumedReadWriteOps", Statistics: "Sum",
		CalcFunc: func(val float64, period float64) float64 { return val },
	},
}

var graphdef = map[string]mp.Graphs{
	"ec2.ebs.bandwidth.#": {
		Label: "EBS Bandwidth",
		Unit:  "bytes/sec",
		Metrics: []mp.Metrics{
			{Name: "read", Label: "Read", Diff: false},
			{Name: "write", Label: "Write", Diff: false},
		},
	},
	"ec2.ebs.throughput.#": {
		Label: "EBS Throughput (op/s)",
		Unit:  "iops",
		Metrics: []mp.Metrics{
			{Name: "read", Label: "Read", Diff: false},
			{Name: "write", Label: "Write", Diff: false},
		},
	},
	"ec2.ebs.size_per_op.#": {
		Label: "EBS Avg Op Size (Bytes/op)",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "read", Label: "Read", Diff: false},
			{Name: "write", Label: "Write", Diff: false},
		},
	},
	"ec2.ebs.latency.#": {
		Label: "EBS Avg Latency (ms/op)",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "read", Label: "Read", Diff: false},
			{Name: "write", Label: "Write", Diff: false},
		},
	},
	"ec2.ebs.queue_length.#": {
		Label: "EBS Avg Queue Length (ops)",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "queue_length", Label: "Queue Length", Diff: false},
		},
	},
	"ec2.ebs.idle_time.#": {
		Label: "EBS Time Spent Idle",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "idle_time", Label: "Idle Time", Diff: false},
		},
	},
	"ec2.ebs.throughput_delivered.#": {
		Label: "EBS Throughput of Provisioned IOPS",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "throughput_delivered", Label: "Throughput", Diff: false},
		},
	},
	"ec2.ebs.consumed_ops.#": {
		Label: "EBS Consumed Ops of Provisioned IOPS",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "consumed_ops", Label: "Consumed Ops", Diff: false},
		},
	},
}

var stderrLogger *log.Logger

// EBSPlugin mackerel plugin for ebs
type EBSPlugin struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	InstanceID      string
	Credentials     *credentials.Credentials
	EC2             *ec2.EC2
	CloudWatch      *cloudwatch.CloudWatch
	Volumes         []*ec2.Volume
}

func (p *EBSPlugin) prepare() error {
	if p.AccessKeyID != "" && p.SecretAccessKey != "" {
		p.Credentials = credentials.NewStaticCredentials(p.AccessKeyID, p.SecretAccessKey, "")
	}

	p.EC2 = ec2.New(session.New(&aws.Config{Credentials: p.Credentials, Region: &p.Region}))
	resp, err := p.EC2.DescribeVolumes(&ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("attachment.instance-id"),
				Values: []*string{
					&p.InstanceID,
				},
			},
		},
	})
	if err != nil {
		return err
	}
	if resp.NextToken != nil {
		return errors.New("DescribeVolumes response has NextToken")
	}

	p.Volumes = resp.Volumes
	if len(p.Volumes) == 0 {
		return errors.New("DescribeVolumes response has no volumes")
	}

	return nil
}

func (p EBSPlugin) getLastPoint(vol *ec2.Volume, metricName string, statType string) (float64, int, error) {
	now := time.Now()

	period := metricPeriodDefault
	if tmp, ok := metricPeriodByVolumeType[*vol.VolumeType]; ok {
		period = tmp
	}
	start := now.Add(time.Duration(period) * 3 * time.Second * -1)

	resp, err := p.CloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("VolumeId"),
				Value: vol.VolumeId,
			},
		},
		StartTime:  &start,
		EndTime:    &now,
		MetricName: &metricName,
		Period:     aws.Int64(60),
		Statistics: []*string{&statType},
		Namespace:  aws.String("AWS/EBS"),
	})
	if err != nil {
		return 0, 0, err
	}

	datapoints := resp.Datapoints
	if len(datapoints) == 0 {
		return 0, 0, errors.New("fetched no datapoints")
	}

	latest := time.Unix(0, 0)
	var latestVal float64
	for _, dp := range datapoints {
		if dp.Timestamp.Before(latest) {
			continue
		}

		latest = *dp.Timestamp
		switch statType {
		case "Average":
			latestVal = *dp.Average
		case "Sum":
			latestVal = *dp.Sum
		}
	}

	return latestVal, period, nil
}

// FetchMetrics fetch the metrics
func (p EBSPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})
	p.CloudWatch = cloudwatch.New(session.New(&aws.Config{Credentials: p.Credentials, Region: &p.Region}))
	for _, vol := range p.Volumes {
		volumeID := normalizeVolumeID(*vol.VolumeId)
		graphs := defaultGraphs
		if *vol.VolumeType == "io1" {
			graphs = allGraphs
		}
		for _, graphName := range graphs {
			for _, metric := range graphdef[graphName].Metrics {
				metricKey := graphName + "." + metric.Name
				cloudwatchdef := cloudwatchdefs[metricKey]
				val, period, err := p.getLastPoint(vol, cloudwatchdef.MetricName, cloudwatchdef.Statistics)
				if err != nil {
					retErr := errors.New(volumeID + " " + err.Error() + ":" + cloudwatchdef.MetricName)
					if err.Error() == "fetched no datapoints" {
						getStderrLogger().Println(retErr)
					} else {
						return nil, retErr
					}
				} else {
					stat[strings.Replace(metricKey, "#", volumeID, -1)] = cloudwatchdef.CalcFunc(val, float64(period))
				}
			}
		}
	}
	return stat, nil
}

// GraphDefinition for plugin
func (p EBSPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

func getStderrLogger() *log.Logger {
	if stderrLogger == nil {
		stderrLogger = log.New(os.Stderr, "", log.LstdFlags)
	}
	return stderrLogger
}

func normalizeVolumeID(volumeID string) string {
	return strings.Replace(volumeID, ".", "_", -1)
}

// Do the plugin
func Do() {
	optRegion := flag.String("region", "", "AWS Region")
	optInstanceID := flag.String("instance-id", "", "Instance ID")
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var ebs EBSPlugin

	ebs.Region = *optRegion
	ebs.InstanceID = *optInstanceID

	// get metadata in ec2 instance
	ec2MC := ec2metadata.New(session.New())
	if *optRegion == "" {
		ebs.Region, _ = ec2MC.Region()
	}
	if *optInstanceID == "" {
		ebs.InstanceID, _ = ec2MC.GetMetadata("instance-id")
	}

	ebs.AccessKeyID = *optAccessKeyID
	ebs.SecretAccessKey = *optSecretAccessKey

	if err := ebs.prepare(); err != nil {
		log.Fatalln(err)
	}

	helper := mp.NewMackerelPlugin(ebs)
	helper.Tempfile = *optTempfile

	helper.Run()
}
