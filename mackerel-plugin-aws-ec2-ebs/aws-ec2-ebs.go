package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

var metricPeriodDefault = 300
var metricPeriodByVolumeType = map[string]int{
	"io1": 60,
}

var volumeTypesHavingExtraMetrics = []string{
	"io1",
}

var defaultGraphs = []string{
	"ec2.ebs.bandwidth",
	"ec2.ebs.throughput",
	"ec2.ebs.size_per_op",
	"ec2.ebs.latency",
	"ec2.ebs.queue_length",
	"ec2.ebs.idle_time",
}

var extraGraphs = []string{
	"ec2.ebs.throughput_delivered",
	"ec2.ebs.consumed_ops",
}

type cloudWatchSetting struct {
	MetricName string
	Statistics string
	CalcFunc   func(float64, float64) float64
}

// http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/monitoring-volume-status.html
var cloudwatchdef = map[string](cloudWatchSetting){
	"bw_%s_read": cloudWatchSetting{
		MetricName: "VolumeReadBytes", Statistics: "Sum",
		CalcFunc: func(val float64, period float64) float64 { return val / period },
	},
	"bw_%s_write": cloudWatchSetting{
		MetricName: "VolumeWriteBytes", Statistics: "Sum",
		CalcFunc: func(val float64, period float64) float64 { return val / period },
	},
	"throughput_%s_read": cloudWatchSetting{
		MetricName: "VolumeReadOps", Statistics: "Sum",
		CalcFunc: func(val float64, period float64) float64 { return val / period },
	},
	"throughput_%s_write": cloudWatchSetting{
		MetricName: "VolumeWriteOps", Statistics: "Sum",
		CalcFunc: func(val float64, period float64) float64 { return val / period },
	},
	"size_per_op_%s_read": cloudWatchSetting{
		MetricName: "VolumeReadBytes", Statistics: "Average",
		CalcFunc: func(val float64, period float64) float64 { return val },
	},
	"size_per_op_%s_write": cloudWatchSetting{
		MetricName: "VolumeWriteBytes", Statistics: "Average",
		CalcFunc: func(val float64, period float64) float64 { return val },
	},
	"latency_%s_read": cloudWatchSetting{
		MetricName: "VolumeTotalReadTime", Statistics: "Average",
		CalcFunc: func(val float64, period float64) float64 { return val * 1000 },
	},
	"latency_%s_write": cloudWatchSetting{
		MetricName: "VolumeTotalWriteTime", Statistics: "Average",
		CalcFunc: func(val float64, period float64) float64 { return val * 1000 },
	},
	"queue_length_%s": cloudWatchSetting{
		MetricName: "VolumeQueueLength", Statistics: "Average",
		CalcFunc: func(val float64, period float64) float64 { return val },
	},
	"idle_time_%s": cloudWatchSetting{
		MetricName: "VolumeIdleTime", Statistics: "Sum",
		CalcFunc: func(val float64, period float64) float64 { return val / period * 100 },
	},
	"throughput_delivered_%s": cloudWatchSetting{
		MetricName: "VolumeThroughputPercentage", Statistics: "Average",
		CalcFunc: func(val float64, period float64) float64 { return val },
	},
	"consumed_ops_%s": cloudWatchSetting{
		MetricName: "VolumeConsumedReadWriteOps", Statistics: "Sum",
		CalcFunc: func(val float64, period float64) float64 { return val },
	},
}

var graphdef = map[string](mp.Graphs){
	"ec2.ebs.bandwidth": mp.Graphs{
		Label: "EBS Bandwidth",
		Unit:  "bytes/sec",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "bw_%s_read", Label: "%s Read", Diff: false},
			mp.Metrics{Name: "bw_%s_write", Label: "%s Write", Diff: false},
		},
	},
	"ec2.ebs.throughput": mp.Graphs{
		Label: "EBS Throughput (op/s)",
		Unit:  "iops",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "throughput_%s_read", Label: "%s Read", Diff: false},
			mp.Metrics{Name: "throughput_%s_write", Label: "%s Write", Diff: false},
		},
	},
	"ec2.ebs.size_per_op": mp.Graphs{
		Label: "EBS Avg Op Size (Bytes/op)",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "size_per_op_%s_read", Label: "%s Read", Diff: false},
			mp.Metrics{Name: "size_per_op_%s_write", Label: "%s Write", Diff: false},
		},
	},
	"ec2.ebs.latency": mp.Graphs{
		Label: "EBS Avg Latency (ms/op)",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "latency_%s_read", Label: "%s Read", Diff: false},
			mp.Metrics{Name: "latency_%s_write", Label: "%s Write", Diff: false},
		},
	},
	"ec2.ebs.queue_length": mp.Graphs{
		Label: "EBS Avg Queue Length (ops)",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "queue_length_%s", Label: "%s", Diff: false},
		},
	},
	"ec2.ebs.idle_time": mp.Graphs{
		Label: "EBS Time Spent Idle (%)",
		Unit:  "percentage",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "idle_time_%s", Label: "%s", Diff: false},
		},
	},
	"ec2.ebs.throughput_delivered": mp.Graphs{
		Label: "EBS Throughput of Provisioned IOPS (%)",
		Unit:  "percentage",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "throughput_delivered_%s", Label: "%s", Diff: false},
		},
	},
	"ec2.ebs.consumed_ops": mp.Graphs{
		Label: "EBS Consumed Ops (Provisioned IOPS)",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "consumed_ops_%s", Label: "%s", Diff: false},
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
	Volumes         *[]*ec2.Volume
}

func (p *EBSPlugin) prepare() error {
	if p.AccessKeyID != "" && p.SecretAccessKey != "" {
		p.Credentials = credentials.NewStaticCredentials(p.AccessKeyID, p.SecretAccessKey, "")
	}

	p.EC2 = ec2.New(&aws.Config{Credentials: p.Credentials, Region: &p.Region})
	resp, err := p.EC2.DescribeVolumes(&ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
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

	p.Volumes = &resp.Volumes
	if len(*p.Volumes) == 0 {
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
			&cloudwatch.Dimension{
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
func (p EBSPlugin) FetchMetrics() (map[string]float64, error) {
	stat := make(map[string]float64)
	p.CloudWatch = cloudwatch.New(&aws.Config{Credentials: p.Credentials, Region: &p.Region})

	for _, vol := range *p.Volumes {
		graphs := graphsToProcess(vol.VolumeType)
		for _, graph := range *graphs {
			for _, met := range graphdef[graph].Metrics {
				cwdef := cloudwatchdef[met.Name]
				val, period, err := p.getLastPoint(vol, cwdef.MetricName, cwdef.Statistics)
				skey := fmt.Sprintf(met.Name, *vol.VolumeId)

				if err != nil {
					retErr := errors.New(*vol.VolumeId + " " + err.Error() + ":" + cwdef.MetricName)
					if err.Error() == "fetched no datapoints" {
						getStderrLogger().Println(retErr)
					} else {
						return nil, retErr
					}
				} else {
					stat[skey] = cwdef.CalcFunc(val, float64(period))
				}
			}
		}
	}

	return stat, nil
}

// GraphDefinition for plugin
func (p EBSPlugin) GraphDefinition() map[string](mp.Graphs) {
	graphMetrics := map[string]([]mp.Metrics){}
	for _, vol := range *p.Volumes {
		graphs := graphsToProcess(vol.VolumeType)

		for _, graph := range *graphs {
			if _, ok := graphMetrics[graph]; !ok {
				graphMetrics[graph] = []mp.Metrics{}
			}

			for _, metric := range graphdef[graph].Metrics {
				m := mp.Metrics{
					Name:  fmt.Sprintf(metric.Name, *vol.VolumeId),
					Label: fmt.Sprintf(metric.Label, *vol.VolumeId+":"+*vol.Attachments[0].Device),
					Diff:  metric.Diff,
				}
				graphMetrics[graph] = append(graphMetrics[graph], m)
			}
		}
	}

	for k := range graphdef {
		graphdef[k] = mp.Graphs{
			Label:   graphdef[k].Label,
			Unit:    graphdef[k].Unit,
			Metrics: graphMetrics[k],
		}
	}

	return graphdef
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func graphsToProcess(volumeType *string) *[]string {
	if stringInSlice(*volumeType, volumeTypesHavingExtraMetrics) {
		var graphsWithExtra = append(defaultGraphs, extraGraphs...)
		return &graphsWithExtra
	}
	return &defaultGraphs
}

func getStderrLogger() *log.Logger {
	if stderrLogger == nil {
		stderrLogger = log.New(os.Stderr, "", log.LstdFlags)
	}
	return stderrLogger
}

func main() {
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
	ec2MC := ec2metadata.New(&ec2metadata.Config{})
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
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = "/tmp/mackerel-plugin-ebs"
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
