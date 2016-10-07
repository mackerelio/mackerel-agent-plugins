package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

// RDSPlugin mackerel plugin for amazon RDS
type RDSPlugin struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Identifier      string
	Engine          string
	Prefix          string
	LabelPrefix     string
}

func getLastPoint(cloudWatch *cloudwatch.CloudWatch, dimension *cloudwatch.Dimension, metricName string) (float64, error) {
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

// FetchMetrics interface for mackerel-plugin
func (p RDSPlugin) FetchMetrics() (map[string]float64, error) {
	auth, err := aws.GetAuth(p.AccessKeyID, p.SecretAccessKey, "", time.Now())
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

	for _, met := range p.rdsMetrics() {
		v, err := getLastPoint(cloudWatch, perInstance, met)
		if err == nil {
			stat[met] = v
		} else {
			log.Printf("%s: %s", met, err)
		}
	}

	return stat, nil
}

func (p RDSPlugin) baseGraphDefs() map[string](mp.Graphs) {
	return map[string](mp.Graphs){
		p.Prefix + ".BinLogDiskUsage": mp.Graphs{
			Label: p.LabelPrefix + " BinLogDiskUsage",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "BinLogDiskUsage", Label: "Usage"},
			},
		},
		p.Prefix + ".DiskQueueDepth": mp.Graphs{
			Label: p.LabelPrefix + " BinLogDiskUsage",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "DiskQueueDepth", Label: "Depth"},
			},
		},
		p.Prefix + ".CPUUtilization": mp.Graphs{
			Label: p.LabelPrefix + " CPU Utilization",
			Unit:  "percentage",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "CPUUtilization", Label: "CPUUtilization"},
			},
		},
		p.Prefix + ".DatabaseConnections": mp.Graphs{
			Label: p.LabelPrefix + " Database Connections",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "DatabaseConnections", Label: "DatabaseConnections"},
			},
		},
		p.Prefix + ".FreeableMemory": mp.Graphs{
			Label: p.LabelPrefix + " Freeable Memory",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "FreeableMemory", Label: "FreeableMemory"},
			},
		},
		p.Prefix + ".FreeStorageSpace": mp.Graphs{
			Label: p.LabelPrefix + " Free Storage Space",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "FreeStorageSpace", Label: "FreeStorageSpace"},
			},
		},
		p.Prefix + ".ReplicaLag": mp.Graphs{
			Label: p.LabelPrefix + " Replica Lag",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ReplicaLag", Label: "ReplicaLag"},
			},
		},
		p.Prefix + ".SwapUsage": mp.Graphs{
			Label: p.LabelPrefix + " Swap Usage",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "SwapUsage", Label: "SwapUsage"},
			},
		},
		p.Prefix + ".IOPS": mp.Graphs{
			Label: p.LabelPrefix + " IOPS",
			Unit:  "iops",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ReadIOPS", Label: "Read"},
				mp.Metrics{Name: "WriteIOPS", Label: "Write"},
			},
		},
		p.Prefix + ".Latency": mp.Graphs{
			Label: p.LabelPrefix + " Latency in second",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ReadLatency", Label: "Read"},
				mp.Metrics{Name: "WriteLatency", Label: "Write"},
			},
		},
		p.Prefix + ".Throughput": mp.Graphs{
			Label: p.LabelPrefix + " Throughput",
			Unit:  "bytes/sec",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ReadThroughput", Label: "Read"},
				mp.Metrics{Name: "WriteThroughput", Label: "Write"},
			},
		},
		p.Prefix + ".NetworkThroughput": mp.Graphs{
			Label: p.LabelPrefix + " Network Throughput",
			Unit:  "bytes/sec",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "NetworkTransmitThroughput", Label: "Transmit"},
				mp.Metrics{Name: "NetworkReceiveThroughput", Label: "Receive"},
			},
		},
	}
}

// GraphDefinition interface for mackerel plugin
func (p RDSPlugin) GraphDefinition() map[string](mp.Graphs) {
	graphdef := p.baseGraphDefs()
	switch p.Engine {
	case "mysql", "mariadb":
		graphdef = mergeGraphDefs(graphdef, p.mySQLGraphDefinition())
	case "aurora":
		graphdef = mergeGraphDefs(graphdef, p.auroraGraphDefinition())
	case "postgresql":
		graphdef = mergeGraphDefs(graphdef, p.postgreSQLGraphDefinition())
	}
	return graphdef
}

func main() {
	optRegion := flag.String("region", "", "AWS Region")
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optIdentifier := flag.String("identifier", "", "DB Instance Identifier")
	optEngine := flag.String("engine", "", "RDS Engine")
	optPrefix := flag.String("metric-key-prefix", "rds", "Metric key prefix")
	optLabelPrefix := flag.String("metric-label-prefix", "", "Metric Label prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	rds := RDSPlugin{
		Prefix: *optPrefix,
	}
	if *optLabelPrefix == "" {
		if *optPrefix == "rds" {
			rds.LabelPrefix = "RDS"
		} else {
			rds.LabelPrefix = strings.Title(*optPrefix)
		}
	} else {
		rds.LabelPrefix = *optLabelPrefix
	}

	if *optRegion == "" {
		rds.Region = aws.InstanceRegion()
	} else {
		rds.Region = *optRegion
	}

	rds.Identifier = *optIdentifier
	rds.AccessKeyID = *optAccessKeyID
	rds.SecretAccessKey = *optSecretAccessKey
	rds.Engine = *optEngine

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
