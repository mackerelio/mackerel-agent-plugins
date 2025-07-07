package mpawsec2ebs

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

const (
	metricPeriodDefault = 300
	aggregationPeriod   = 60
)

var metricPeriodByVolumeType = map[types.VolumeType]int{
	types.VolumeTypeIo1: 60,
}

var baseGraphs = []string{
	"ec2.ebs.bandwidth.#",
	"ec2.ebs.throughput.#",
	"ec2.ebs.size_per_op.#",
	"ec2.ebs.latency.#",
	"ec2.ebs.queue_length.#",
	"ec2.ebs.idle_time.#",
}

var defaultGraphs = append([]string{
	"ec2.ebs.burst_balance.#",
}, baseGraphs...)

var io1Graphs = append([]string{
	"ec2.ebs.throughput_delivered.#",
	"ec2.ebs.consumed_ops.#",
}, baseGraphs...)

type additionalCloudWatchSetting struct {
	MetricName string
	Statistics cloudwatchTypes.Statistic
	CalcFunc   func(float64, float64) float64
}

type cloudWatchSetting struct {
	MetricName string
	Statistics cloudwatchTypes.Statistic
	CalcFunc   func(float64) float64
	Additional *additionalCloudWatchSetting
}

func value(val float64) float64 {
	return val
}

func valuePerSec(val float64) float64 {
	return val / aggregationPeriod
}

func sec2msec(val float64) float64 {
	return val * 1000
}

func valPerOps(val, ops float64) float64 {
	return val / ops
}

// http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/monitoring-volume-status.html
var cloudwatchdefs = map[string](cloudWatchSetting){
	"ec2.ebs.bandwidth.#.read": cloudWatchSetting{
		MetricName: "VolumeReadBytes", Statistics: cloudwatchTypes.StatisticSum,
		CalcFunc: valuePerSec,
	},
	"ec2.ebs.bandwidth.#.write": cloudWatchSetting{
		MetricName: "VolumeWriteBytes", Statistics: cloudwatchTypes.StatisticSum,
		CalcFunc: valuePerSec,
	},
	"ec2.ebs.throughput.#.read": cloudWatchSetting{
		MetricName: "VolumeReadOps", Statistics: cloudwatchTypes.StatisticSum,
		CalcFunc: valuePerSec,
	},
	"ec2.ebs.throughput.#.write": cloudWatchSetting{
		MetricName: "VolumeWriteOps", Statistics: cloudwatchTypes.StatisticSum,
		CalcFunc: valuePerSec,
	},
	"ec2.ebs.size_per_op.#.read": cloudWatchSetting{
		MetricName: "VolumeReadBytes", Statistics: cloudwatchTypes.StatisticAverage,
		CalcFunc: value,
	},
	"ec2.ebs.size_per_op.#.write": cloudWatchSetting{
		MetricName: "VolumeWriteBytes", Statistics: cloudwatchTypes.StatisticAverage,
		CalcFunc: value,
	},
	"ec2.ebs.latency.#.read": cloudWatchSetting{
		MetricName: "VolumeTotalReadTime", Statistics: cloudwatchTypes.StatisticAverage,
		CalcFunc: sec2msec,
	},
	"ec2.ebs.latency.#.write": cloudWatchSetting{
		MetricName: "VolumeTotalWriteTime", Statistics: cloudwatchTypes.StatisticAverage,
		CalcFunc: sec2msec,
	},
	"ec2.ebs.queue_length.#.queue_length": cloudWatchSetting{
		MetricName: "VolumeQueueLength", Statistics: cloudwatchTypes.StatisticAverage,
		CalcFunc: value,
	},
	"ec2.ebs.idle_time.#.idle_time": cloudWatchSetting{
		MetricName: "VolumeIdleTime", Statistics: cloudwatchTypes.StatisticSum,
		CalcFunc: func(val float64) float64 { return val / aggregationPeriod * 100.0 },
	},
	"ec2.ebs.throughput_delivered.#.throughput_delivered": cloudWatchSetting{
		MetricName: "VolumeThroughputPercentage", Statistics: cloudwatchTypes.StatisticAverage,
		CalcFunc: value,
	},
	"ec2.ebs.consumed_ops.#.consumed_ops": cloudWatchSetting{
		MetricName: "VolumeConsumedReadWriteOps", Statistics: cloudwatchTypes.StatisticSum,
		CalcFunc: value,
	},
	"ec2.ebs.burst_balance.#.burst_balance": cloudWatchSetting{
		MetricName: "BurstBalance", Statistics: cloudwatchTypes.StatisticAverage,
		CalcFunc: value,
	},
}

var cloudwatchdefsNitro = map[string](cloudWatchSetting){
	"ec2.ebs.size_per_op.#.read": cloudWatchSetting{
		MetricName: "VolumeReadBytes", Statistics: cloudwatchTypes.StatisticSum,
		Additional: &additionalCloudWatchSetting{
			MetricName: "VolumeReadOps", Statistics: cloudwatchTypes.StatisticSum,
			CalcFunc: valPerOps,
		},
	},
	"ec2.ebs.size_per_op.#.write": cloudWatchSetting{
		MetricName: "VolumeWriteBytes", Statistics: cloudwatchTypes.StatisticSum,
		Additional: &additionalCloudWatchSetting{
			MetricName: "VolumeWriteOps", Statistics: cloudwatchTypes.StatisticSum,
			CalcFunc: valPerOps,
		},
	},
	"ec2.ebs.latency.#.read": cloudWatchSetting{
		MetricName: "VolumeTotalReadTime", Statistics: cloudwatchTypes.StatisticSum,
		Additional: &additionalCloudWatchSetting{
			MetricName: "VolumeReadOps", Statistics: cloudwatchTypes.StatisticSum,
			CalcFunc: valPerOps,
		},
	},
	"ec2.ebs.latency.#.write": cloudWatchSetting{
		MetricName: "VolumeTotalWriteTime", Statistics: cloudwatchTypes.StatisticSum,
		Additional: &additionalCloudWatchSetting{
			MetricName: "VolumeWriteOps", Statistics: cloudwatchTypes.StatisticSum,
			CalcFunc: valPerOps,
		},
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
	"ec2.ebs.burst_balance.#": {
		Label: "EBS Burst Balance",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "burst_balance", Label: "Burst Balance", Diff: false},
		},
	},
}

// EBSPlugin mackerel plugin for ebs
type EBSPlugin struct {
	// command line options
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	InstanceID      string

	// internal states
	EC2        *ec2.Client
	CloudWatch *cloudwatch.Client
	Volumes    []types.Volume
	Hypervisor types.InstanceTypeHypervisor
}

func (p *EBSPlugin) prepare(ctx context.Context) error {
	var opts []func(*config.LoadOptions) error
	if p.AccessKeyID != "" && p.SecretAccessKey != "" {
		opts = append(opts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(p.AccessKeyID, p.SecretAccessKey, "")))
	}
	if p.Region != "" {
		opts = append(opts, config.WithRegion(p.Region))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return err
	}

	p.EC2 = ec2.NewFromConfig(cfg)

	var instanceType types.InstanceType
	instance, err := p.EC2.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{p.InstanceID},
	})
	if err != nil {
		return err
	}
	if instance.NextToken != nil {
		return errors.New("DescribeInstances response has NextToken")
	}
	for i := range instance.Reservations {
		for j := range instance.Reservations[i].Instances {
			instanceType = instance.Reservations[i].Instances[j].InstanceType
		}
	}

	instanceDetail, err := p.EC2.DescribeInstanceTypes(ctx, &ec2.DescribeInstanceTypesInput{
		InstanceTypes: []types.InstanceType{instanceType},
	})
	if err != nil {
		return err
	}
	if instanceDetail.NextToken != nil {
		return errors.New("DescribeInstanceTypes response has NextToken")
	}
	for i := range instanceDetail.InstanceTypes {
		p.Hypervisor = instanceDetail.InstanceTypes[i].Hypervisor
	}

	resp, err := p.EC2.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("attachment.instance-id"),
				Values: []string{p.InstanceID},
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

	p.CloudWatch = cloudwatch.NewFromConfig(cfg)

	return nil
}

var errNoDataPoint = errors.New("fetched no datapoints")

func (p EBSPlugin) getLastPoint(ctx context.Context, vol types.Volume, metricName string, statType cloudwatchTypes.Statistic) (float64, error) {
	now := time.Now()

	period := metricPeriodDefault
	if tmp, ok := metricPeriodByVolumeType[vol.VolumeType]; ok {
		period = tmp
	}
	start := now.Add(time.Duration(period) * 3 * time.Second * -1)

	resp, err := p.CloudWatch.GetMetricStatistics(ctx, &cloudwatch.GetMetricStatisticsInput{
		Dimensions: []cloudwatchTypes.Dimension{
			{
				Name:  aws.String("VolumeId"),
				Value: vol.VolumeId,
			},
		},
		StartTime:  &start,
		EndTime:    &now,
		MetricName: &metricName,
		Period:     aws.Int32(aggregationPeriod),
		Statistics: []cloudwatchTypes.Statistic{statType},
		Namespace:  aws.String("AWS/EBS"),
	})
	if err != nil {
		return 0, err
	}

	datapoints := resp.Datapoints
	if len(datapoints) == 0 {
		return 0, errNoDataPoint
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

	return latestVal, nil
}

func (p EBSPlugin) fetch(ctx context.Context, volume types.Volume, setting cloudWatchSetting) (float64, error) {
	val, err := p.getLastPoint(ctx, volume, setting.MetricName, setting.Statistics)
	if err != nil {
		return 0, fmt.Errorf("%s %w : %s", *volume.VolumeId, err, setting.MetricName)
	}

	if setting.Additional == nil {
		return setting.CalcFunc(val), nil
	}

	val2, err := p.getLastPoint(ctx, volume, setting.Additional.MetricName, setting.Additional.Statistics)
	if err != nil {
		return 0, fmt.Errorf("%s %w : %s", *volume.VolumeId, err, setting.Additional.MetricName)
	}
	return setting.Additional.CalcFunc(val, val2), nil
}

// FetchMetrics fetch the metrics
func (p EBSPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})

	// Override when Nitro instance.
	if p.Hypervisor == types.InstanceTypeHypervisorNitro {
		for i := range cloudwatchdefsNitro {
			cloudwatchdefs[i] = cloudwatchdefsNitro[i]
		}
	}

	for _, vol := range p.Volumes {
		volumeID := normalizeVolumeID(*vol.VolumeId)
		var graphs []string
		if vol.VolumeType == types.VolumeTypeIo1 {
			graphs = io1Graphs
		} else {
			graphs = defaultGraphs
		}
		for _, graphName := range graphs {
			for _, metric := range graphdef[graphName].Metrics {
				metricKey := graphName + "." + metric.Name
				cloudwatchdef := cloudwatchdefs[metricKey]
				val, err := p.fetch(context.TODO(), vol, cloudwatchdef)
				if err != nil {
					if errors.Is(err, errNoDataPoint) {
						// nop
					} else {
						return nil, err
					}
				} else {
					stat[strings.Replace(metricKey, "#", volumeID, -1)] = val
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

func normalizeVolumeID(volumeID string) string {
	return strings.Replace(volumeID, ".", "_", -1)
}

// overwritten with syscall.SIGTERM on unix environment (see aws-ec2-ebs_unix.go)
var defaultSignal = os.Interrupt

// Do the plugin
func Do() {
	optRegion := flag.String("region", "", "AWS Region")
	optInstanceID := flag.String("instance-id", "", "Instance ID")
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), defaultSignal)
	defer stop()

	var ebs EBSPlugin

	ebs.Region = *optRegion
	ebs.InstanceID = *optInstanceID

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// get metadata in ec2 instance
	imdsClient := imds.NewFromConfig(cfg)
	if *optRegion == "" {
		out, err := imdsClient.GetRegion(ctx, nil)
		if err != nil {
			log.Fatalln(err)
		}
		ebs.Region = out.Region
	}
	if *optInstanceID == "" {
		metadata, err := imdsClient.GetMetadata(ctx, &imds.GetMetadataInput{
			Path: "instance-id",
		})
		if err != nil {
			log.Fatalln(err)
		}
		content, _ := io.ReadAll(metadata.Content)
		ebs.InstanceID = string(content)
	}

	ebs.AccessKeyID = *optAccessKeyID
	ebs.SecretAccessKey = *optSecretAccessKey

	if err := ebs.prepare(ctx); err != nil {
		log.Fatalln(err)
	}

	helper := mp.NewMackerelPlugin(ebs)
	helper.Tempfile = *optTempfile

	helper.Run()
}
