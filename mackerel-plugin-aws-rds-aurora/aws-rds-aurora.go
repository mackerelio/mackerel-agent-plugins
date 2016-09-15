package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

// AuroraPlugin mackerel plugin for amazon Aurora
type AuroraPlugin struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Identifier      string
	Prefix          string
	LabelPrefix     string
	Metrics         []string
}

var metricsdefAurora = []string{
	"CPUUtilization", "DatabaseConnections", "FreeableMemory", "FreeLocalStorage",
	"NetworkReceiveThroughput", "NetworkThroughput", "NetworkTransmitThroughput",
	"BinLogDiskUsage", "Deadlocks", "ActiveTransactions", "BlockedTransactions",
	"EngineUptime", "Queries", "LoginFailures",
	"ResultSetCacheHitRatio", "BufferCacheHitRatio",
	"AuroraBinlogReplicaLag", "AuroraReplicaLag", "AuroraReplicaLagMaximum", "AuroraReplicaLagMinimum",
	"CommitLatency", "DDLLatency", "DMLLatency", "DeleteLatency",
	"InsertLatency", "SelectLatency", "UpdateLatency",
	"CommitThroughput", "DDLThroughput", "DMLThroughput", "DeleteThroughput",
	"InsertThroughput", "SelectThroughput", "UpdateThroughput",
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
func (p AuroraPlugin) FetchMetrics() (map[string]float64, error) {
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

	for _, met := range p.Metrics {
		v, err := getLastPoint(cloudWatch, perInstance, met)
		if err == nil {
			stat[met] = v
		} else {
			log.Printf("%s: %s", met, err)
		}
	}

	return stat, nil
}

// GraphDefinition interface for mackerel plugin
func (p AuroraPlugin) GraphDefinition() map[string](mp.Graphs) {
	graphdef := map[string](mp.Graphs){
		p.LabelPrefix + ".CPUUtilization": mp.Graphs{
			Label: p.LabelPrefix + " CPU Utilization",
			Unit:  "percentage",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "CPUUtilization", Label: "CPUUtilization"},
			},
		},
		p.LabelPrefix + ".DatabaseConnections": mp.Graphs{
			Label: p.LabelPrefix + " Database Connections",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "DatabaseConnections", Label: "DatabaseConnections"},
			},
		},
		p.LabelPrefix + ".FreeableMemory": mp.Graphs{
			Label: p.LabelPrefix + " Freeable Memory",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "FreeableMemory", Label: "FreeableMemory"},
			},
		},
		p.LabelPrefix + ".FreeStorageSpace": mp.Graphs{
			Label: p.LabelPrefix + " Free Storage Space",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "FreeStorageSpace", Label: "FreeStorageSpace"},
			},
		},
		p.LabelPrefix + ".NetworkThroughput": mp.Graphs{
			Label: p.LabelPrefix + " Network Throughput",
			Unit:  "bytes/sec",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "NetworkThroughput", Label: "Throughput"},
				mp.Metrics{Name: "NetworkTransmitThroughput", Label: "Transmit"},
				mp.Metrics{Name: "NetworkReceiveThroughput", Label: "Receive"},
			},
		},
		p.LabelPrefix + ".BinLogDiskUsage": mp.Graphs{
			Label: p.LabelPrefix + " BinLog Disk Usage",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "BinLogDiskUsage", Label: "BinLogDiskUsage"},
			},
		},
		p.LabelPrefix + ".Deadlocks": mp.Graphs{
			Label: p.LabelPrefix + " Dead Locks",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "Deadlocks", Label: "Deadlocks"},
			},
		},
		p.LabelPrefix + ".Transaction": mp.Graphs{
			Label: p.LabelPrefix + " Transaction",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ActiveTransactions", Label: "Active"},
				mp.Metrics{Name: "BlockedTransactions", Label: "Blocked"},
			},
		},
		p.LabelPrefix + ".EngineUptime": mp.Graphs{
			Label: p.LabelPrefix + " Engine Uptime",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "EngineUptime", Label: "EngineUptime"},
			},
		},
		p.LabelPrefix + ".Queries": mp.Graphs{
			Label: p.LabelPrefix + " Queries",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "Queries", Label: "Queries"},
			},
		},
		p.LabelPrefix + ".LoginFailures": mp.Graphs{
			Label: p.LabelPrefix + " Login Failures",
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "LoginFailures", Label: "LoginFailures"},
			},
		},
		p.LabelPrefix + ".CacheHitRatio": mp.Graphs{
			Label: p.LabelPrefix + " Cache Hit Ratio",
			Unit:  "percentage",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ResultSetCacheHitRatio", Label: "ResultSet"},
				mp.Metrics{Name: "BufferCacheHitRatio", Label: "Buffer"},
			},
		},
		p.LabelPrefix + ".AuroraBinlogReplicaLag": mp.Graphs{
			Label: p.LabelPrefix + " Aurora Binlog ReplicaLag",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "AuroraBinlogReplicaLag", Label: "AuroraBinlogReplicaLag"},
			},
		},
		p.LabelPrefix + ".AuroraReplicaLag": mp.Graphs{
			Label: p.LabelPrefix + " Aurora ReplicaLag",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "AuroraReplicaLag", Label: "ReplicaLag"},
				mp.Metrics{Name: "AuroraReplicaLagMaximum", Label: "ReplicaLagMaximum"},
				mp.Metrics{Name: "AuroraBinlogReplicaLag", Label: "ReplicaLagMinimum"},
			},
		},
		p.LabelPrefix + ".Latency": mp.Graphs{
			Label: p.LabelPrefix + " Latency",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "SelectLatency", Label: "Select"},
				mp.Metrics{Name: "InsertLatency", Label: "Insert"},
				mp.Metrics{Name: "UpdateLatency", Label: "Update"},
				mp.Metrics{Name: "DeleteLatency", Label: "Delete"},
				mp.Metrics{Name: "CommitLatency", Label: "Commit"},
				mp.Metrics{Name: "DDLLatency", Label: "DDL"},
				mp.Metrics{Name: "DMLLatency", Label: "DML"},
			},
		},
		p.LabelPrefix + ".Throughput": mp.Graphs{
			Label: p.LabelPrefix + " Throughput",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "SelectThroughput", Label: "Select"},
				mp.Metrics{Name: "InsertThroughput", Label: "Insert"},
				mp.Metrics{Name: "UpdateThroughput", Label: "Update"},
				mp.Metrics{Name: "DeleteThroughput", Label: "Delete"},
				mp.Metrics{Name: "CommitThroughput", Label: "Commit"},
				mp.Metrics{Name: "DDLThroughput", Label: "DDL"},
				mp.Metrics{Name: "DMLThroughput", Label: "DML"},
			},
		},
	}

	return graphdef
}

func main() {
	optRegion := flag.String("region", "", "AWS Region")
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optIdentifier := flag.String("identifier", "", "DB Instance Identifier")
	optLabelPrefix := flag.String("metric-label-prefix", "", "Metric Label prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var aurora AuroraPlugin

	if *optLabelPrefix == "" {
		aurora.LabelPrefix = "Aurora"
	} else {
		aurora.LabelPrefix = *optLabelPrefix
	}

	if *optRegion == "" {
		aurora.Region = aws.InstanceRegion()
	} else {
		aurora.Region = *optRegion
	}

	aurora.Identifier = *optIdentifier
	aurora.AccessKeyID = *optAccessKeyID
	aurora.SecretAccessKey = *optSecretAccessKey
	aurora.Metrics = metricsdefAurora

	helper := mp.NewMackerelPlugin(aurora)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-aws-rds-aurora")
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
