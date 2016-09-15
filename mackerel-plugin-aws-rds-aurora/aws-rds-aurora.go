package main

import (
	"flag"

	"github.com/crowdmob/goamz/aws"
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

// GraphDefinition interface for mackerel plugin
func (p AuroraPlugin) GraphDefinition() map[string](mp.Graphs) {
	graphdef := map[string](mp.Graphs){
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
		p.Prefix + ".NetworkThroughput": mp.Graphs{
			Label: p.LabelPrefix + " Network Throughput",
			Unit:  "bytes/sec",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "NetworkThroughput", Label: "Throughput"},
				mp.Metrics{Name: "NetworkTransmitThroughput", Label: "Transmit"},
				mp.Metrics{Name: "NetworkReceiveThroughput", Label: "Receive"},
			},
		},
		p.Prefix + ".BinLogDiskUsage": mp.Graphs{
			Label: p.LabelPrefix + " BinLog Disk Usage",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "BinLogDiskUsage", Label: "BinLogDiskUsage"},
			},
		},
		p.Prefix + ".Deadlocks": mp.Graphs{
			Label: p.LabelPrefix + " Dead Locks",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "Deadlocks", Label: "Deadlocks"},
			},
		},
		p.Prefix + ".Transaction": mp.Graphs{
			Label: p.LabelPrefix + " Transaction",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ActiveTransactions", Label: "Active"},
				mp.Metrics{Name: "BlockedTransactions", Label: "Blocked"},
			},
		},
		p.Prefix + ".EngineUptime": mp.Graphs{
			Label: p.LabelPrefix + " Engine Uptime",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "EngineUptime", Label: "EngineUptime"},
			},
		},
		p.Prefix + ".Queries": mp.Graphs{
			Label: p.LabelPrefix + " Queries",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "Queries", Label: "Queries"},
			},
		},
		p.Prefix + ".LoginFailures": mp.Graphs{
			Label: p.LabelPrefix + " Login Failures",
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "LoginFailures", Label: "LoginFailures"},
			},
		},
		p.Prefix + ".CacheHitRatio": mp.Graphs{
			Label: p.LabelPrefix + " Cache Hit Ratio",
			Unit:  "percentage",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ResultSetCacheHitRatio", Label: "ResultSet"},
				mp.Metrics{Name: "BufferCacheHitRatio", Label: "Buffer"},
			},
		},
		p.Prefix + ".AuroraBinlogReplicaLag": mp.Graphs{
			Label: p.LabelPrefix + " Aurora Binlog ReplicaLag",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "AuroraBinlogReplicaLag", Label: "AuroraBinlogReplicaLag"},
			},
		},
		p.Prefix + ".AuroraReplicaLag": mp.Graphs{
			Label: p.LabelPrefix + " Aurora ReplicaLag",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "AuroraReplicaLag", Label: "ReplicaLag"},
				mp.Metrics{Name: "AuroraReplicaLagMaximum", Label: "ReplicaLagMaximum"},
				mp.Metrics{Name: "AuroraBinlogReplicaLag", Label: "ReplicaLagMinimum"},
			},
		},
		p.Prefix + ".Latency": mp.Graphs{
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
		p.Prefix + ".Throughput": mp.Graphs{
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
}
