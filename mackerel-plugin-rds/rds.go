package main

import (
	"errors"
	"flag"
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"log"
	"os"
	"time"
)

var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"rds.CPUUtilization": mp.Graphs{
		Label: "RDS CPU Utilization",
		Unit:  "percentage",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "CPUUtilization", Label: "CPUUtilization"},
		},
	},
	"rds.DatabaseConnections": mp.Graphs{
		Label: "RDS Database Connections",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "DatabaseConnections", Label: "DatabaseConnections"},
		},
	},
	"rds.FreeableMemory": mp.Graphs{
		Label: "RDS Freeable Memory",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "FreeableMemory", Label: "FreeableMemory"},
		},
	},
	"rds.FreeStorageSpace": mp.Graphs{
		Label: "RDS Free Storage Space",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "FreeStorageSpace", Label: "FreeStorageSpace"},
		},
	},
	"rds.ReplicaLag": mp.Graphs{
		Label: "RDS Replica Lag",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "ReplicaLag", Label: "ReplicaLag"},
		},
	},
	"rds.SwapUsage": mp.Graphs{
		Label: "RDS Swap Usage",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "SwapUsage", Label: "SwapUsage"},
		},
	},
	"rds.IOPS": mp.Graphs{
		Label: "RDS IOPS",
		Unit:  "iops",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "ReadIOPS", Label: "Read"},
			mp.Metrics{Name: "WriteIOPS", Label: "Write"},
		},
	},
	"rds.Latency": mp.Graphs{
		Label: "RDS Latency in second",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "ReadLatency", Label: "Read"},
			mp.Metrics{Name: "WriteLatency", Label: "Write"},
		},
	},
}

type RDSPlugin struct {
	Region          string
	AccessKeyId     string
	SecretAccessKey string
	Identifier      string
}

func GetLastPoint(cloudWatch *cloudwatch.CloudWatch, dimension *cloudwatch.Dimension, metricName string) (float64, error) {
	now := time.Now()

	response, err := cloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsRequest{
		Dimensions: []cloudwatch.Dimension{*dimension},
		StartTime:  now.Add(time.Duration(180) * time.Second * -1), // 3 min (to fetch at least 1 data-point)
		EndTime:    now,
		MetricName: metricName,
		Period:     60,
		Statistics: []string{"Average"},
		Namespace:  "AWS/RDS",
	})
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

func (p RDSPlugin) FetchMetrics() (map[string]float64, error) {
	auth, err := aws.GetAuth(p.AccessKeyId, p.SecretAccessKey, "", time.Now())
	if err != nil {
		return nil, err
	}

	cloudWatch, err := cloudwatch.NewCloudWatch(auth, aws.Regions[p.Region].CloudWatchServicepoint)
	if err != nil {
		return nil, err
	}

	stat := make(map[string]float64)

	perInstance := &cloudwatch.Dimension{
		Name:  "DBInstanceIdentifier",
		Value: p.Identifier,
	}

	for _, met := range [...]string{
		"BinLogDiskUsage", "CPUUtilization", "DatabaseConnections", "DiskQueueDepth", "FreeableMemory",
		"FreeStorageSpace", "ReplicaLag", "SwapUsage", "ReadIOPS", "WriteIOPS", "ReadLatency",
		"WriteLatency",
	} {
		v, err := GetLastPoint(cloudWatch, perInstance, met)
		if err == nil {
			stat[met] = v
		} else {
			log.Printf("%s: %s", met, err)
		}
	}

	return stat, nil
}

func (p RDSPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optRegion := flag.String("region", "", "AWS Region")
	optAccessKeyId := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optIdentifier := flag.String("identifier", "", "DB Instance Identifier")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var rds RDSPlugin

	if *optRegion == "" {
		rds.Region = aws.InstanceRegion()
	} else {
		rds.Region = *optRegion
	}

	rds.Identifier = *optIdentifier
	rds.AccessKeyId = *optAccessKeyId
	rds.SecretAccessKey = *optSecretAccessKey

	helper := mp.NewMackerelPlugin(rds)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = "/tmp/mackerel-plugin-rds"
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
