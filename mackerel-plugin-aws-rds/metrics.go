package main

import (
	mp "github.com/mackerelio/go-mackerel-plugin"
)

// MetricsdefMySQL lastpoint metrics name to get in Cloudwatch
var MetricsdefMySQL = []string{
	"CPUUtilization", "BinLogDiskUsage", "DatabaseConnections",
	"FreeableMemory", "FreeStorageSpace", "SwapUsage",
	"NetworkTransmitThroughput", "NetworkReceiveThroughput", "DiskQueueDepth",
	"ReadIOPS", "WriteIOPS",
	"ReadLatency", "WriteLatency",
	"ReadThroughput", "WriteThroughput",
}

// MetricsdefAurora lastpoint metrics name to get in Cloudwatch
var MetricsdefAurora = []string{
	"CPUUtilization", "DatabaseConnections", "FreeableMemory", "FreeLocalStorage",
	"NetworkReceiveThroughput", "NetworkThroughput", "NetworkTransmitThroughput",
	"BinLogDiskUsage", "Deadlocks", "ActiveTransactions", "BlockedTransactions",
	"EngineUptime", "Queries", "LoginFailures",
	"ResultSetCacheHitRatio", "BufferCacheHitRatio",
	"AuroraBinlogReplicaLag", "AuroraReplicaLagMaximum", "AuroraReplicaLagMinimum",
	"CommitLatency", "DDLLatency", "DMLLatency", "DeleteLatency",
	"InsertLatency", "SelectLatency", "UpdateLatency",
	"CommitThroughput", "DDLThroughput", "DMLThroughput", "DeleteThroughput",
	"InsertThroughput", "SelectThroughput", "UpdateThroughput",
}

// MetricsdefMariaDB lastpoint metrics name to get in Cloudwatch
var MetricsdefMariaDB = []string{
	"CPUUtilization", "BinLogDiskUsage",
	"DatabaseConnections", "FreeableMemory", "FreeStorageSpace", "SwapUsage",
	"NetworkTransmitThroughput", "NetworkReceiveThroughput", "DiskQueueDepth",
	"ReadIOPS", "WriteIOPS",
	"ReadLatency", "WriteLatency",
	"ReadThroughput", "WriteThroughput",
}

// MetricsdefPostgreSQL lastpoint metrics name to get in Cloudwatch
var MetricsdefPostgreSQL = []string{
	"CPUUtilization", "DatabaseConnections",
	"FreeableMemory", "FreeStorageSpace", "MaximumUsedTransactionIDs",
	"NetworkTransmitThroughput", "NetworkReceiveThroughput", "DiskQueueDepth",
	"SwapUsage", "OldestReplicationSlotLag", "ReplicationSlotDiskUsage",
	"TransactionLogsDiskUsage", "TransactionLogsGeneration",
	"ReadIOPS", "WriteIOPS",
	"ReadLatency", "WriteLatency",
	"ReadThroughput", "WriteThroughput",
}

// MySQLGraphDefinition graph definition
func (p RDSPlugin) MySQLGraphDefinition() map[string](mp.Graphs) {
	graphdef := map[string](mp.Graphs){
		p.LabelPrefix + ".CPUUtilization": mp.Graphs{
			Label: p.LabelPrefix + " CPU Utilization",
			Unit:  "percentage",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "CPUUtilization", Label: "CPUUtilization"},
			},
		},
		p.LabelPrefix + ".CPUCreditBalance": mp.Graphs{
			Label: p.LabelPrefix + " CPU CreditBalance",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "CPUCreditBalance", Label: "CPUCreditBalance"},
			},
		},
		p.LabelPrefix + ".CPUCreditUsage": mp.Graphs{
			Label: p.LabelPrefix + " CPU CreditUsage",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "CPUCreditUsage", Label: "CPUCreditUsage"},
			},
		},
		p.LabelPrefix + ".BinLogDiskUsage": mp.Graphs{
			Label: p.LabelPrefix + " BinLog Disk Usage",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "BinLogDiskUsage", Label: "BinLogDiskUsage"},
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
		p.LabelPrefix + ".SwapUsage": mp.Graphs{
			Label: p.LabelPrefix + " Swap Usage",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "SwapUsage", Label: "SwapUsage"},
			},
		},
		p.LabelPrefix + ".NetworkThroughput": mp.Graphs{
			Label: p.LabelPrefix + " Network Throughput",
			Unit:  "bytes/sec",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "NetworkTransmitThroughput", Label: "Transmit"},
				mp.Metrics{Name: "NetworkReceiveThroughput", Label: "Receive"},
			},
		},
		p.LabelPrefix + ".DiskQueueDepth": mp.Graphs{
			Label: p.LabelPrefix + " Disk Queue Depth",
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "DiskQueueDepth", Label: "DiskQueueDepth"},
			},
		},
		p.LabelPrefix + ".ReplicaLag": mp.Graphs{
			Label: p.LabelPrefix + " Replica Lag",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ReplicaLag", Label: "ReplicaLag"},
			},
		},
		p.LabelPrefix + ".IOPS": mp.Graphs{
			Label: p.LabelPrefix + " IOPS",
			Unit:  "iops",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ReadIOPS", Label: "Read"},
				mp.Metrics{Name: "WriteIOPS", Label: "Write"},
			},
		},
		p.LabelPrefix + ".Latency": mp.Graphs{
			Label: p.LabelPrefix + " Latency",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ReadLatency", Label: "Read"},
				mp.Metrics{Name: "WriteLatency", Label: "Write"},
			},
		},
		p.LabelPrefix + ".Throughput": mp.Graphs{
			Label: p.LabelPrefix + " Throughput",
			Unit:  "bytes/sec",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ReadThroughput", Label: "Read"},
				mp.Metrics{Name: "WriteThroughput", Label: "Write"},
			},
		},
	}
	return graphdef
}

// AuroraGraphDefinition graph definition
func (p RDSPlugin) AuroraGraphDefinition() map[string](mp.Graphs) {
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
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "AuroraBinlogReplicaLag", Label: "AuroraBinlogReplicaLag"},
			},
		},
		p.LabelPrefix + ".AuroraReplicaLag": mp.Graphs{
			Label: p.LabelPrefix + " Aurora ReplicaLag",
			Unit:  "float",
			Metrics: [](mp.Metrics){
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

// MariaDBGraphDefinition graph definition
func (p RDSPlugin) MariaDBGraphDefinition() map[string](mp.Graphs) {
	graphdef := map[string](mp.Graphs){
		p.LabelPrefix + ".CPUUtilization": mp.Graphs{
			Label: p.LabelPrefix + " CPU Utilization",
			Unit:  "percentage",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "CPUUtilization", Label: "CPUUtilization"},
			},
		},
		p.LabelPrefix + ".CPUCreditBalance": mp.Graphs{
			Label: p.LabelPrefix + " CPU CreditBalance",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "CPUCreditBalance", Label: "CPUCreditBalance"},
			},
		},
		p.LabelPrefix + ".CPUCreditUsage": mp.Graphs{
			Label: p.LabelPrefix + " CPU CreditUsage",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "CPUCreditUsage", Label: "CPUCreditUsage"},
			},
		},
		p.LabelPrefix + ".BinLogDiskUsage": mp.Graphs{
			Label: p.LabelPrefix + " BinLog Disk Usage",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "BinLogDiskUsage", Label: "BinLogDiskUsage"},
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
		p.LabelPrefix + ".SwapUsage": mp.Graphs{
			Label: p.LabelPrefix + " Swap Usage",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "SwapUsage", Label: "SwapUsage"},
			},
		},
		p.LabelPrefix + ".NetworkThroughput": mp.Graphs{
			Label: p.LabelPrefix + " Network Throughput",
			Unit:  "bytes/sec",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "NetworkTransmitThroughput", Label: "Transmit"},
				mp.Metrics{Name: "NetworkReceiveThroughput", Label: "Receive"},
			},
		},
		p.LabelPrefix + ".DiskQueueDepth": mp.Graphs{
			Label: p.LabelPrefix + " Disk Queue Depth",
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "DiskQueueDepth", Label: "DiskQueueDepth"},
			},
		},
		p.LabelPrefix + ".IOPS": mp.Graphs{
			Label: p.LabelPrefix + " IOPS",
			Unit:  "iops",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ReadIOPS", Label: "Read"},
				mp.Metrics{Name: "WriteIOPS", Label: "Write"},
			},
		},
		p.LabelPrefix + ".Latency": mp.Graphs{
			Label: p.LabelPrefix + " Latency",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ReadLatency", Label: "Read"},
				mp.Metrics{Name: "WriteLatency", Label: "Write"},
			},
		},
		p.LabelPrefix + ".Throughput": mp.Graphs{
			Label: p.LabelPrefix + " Throughput",
			Unit:  "bytes/sec",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ReadThroughput", Label: "Read"},
				mp.Metrics{Name: "WriteThroughput", Label: "Write"},
			},
		},
	}
	return graphdef
}

// PostgreSQLGraphDefinition graph definition
func (p RDSPlugin) PostgreSQLGraphDefinition() map[string](mp.Graphs) {
	graphdef := map[string](mp.Graphs){
		p.LabelPrefix + ".CPUUtilization": mp.Graphs{
			Label: p.LabelPrefix + " CPU Utilization",
			Unit:  "percentage",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "CPUUtilization", Label: "CPUUtilization"},
			},
		},
		p.LabelPrefix + ".CPUCreditBalance": mp.Graphs{
			Label: p.LabelPrefix + " CPU CreditBalance",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "CPUCreditBalance", Label: "CPUCreditBalance"},
			},
		},
		p.LabelPrefix + ".CPUCreditUsage": mp.Graphs{
			Label: p.LabelPrefix + " CPU CreditUsage",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "CPUCreditUsage", Label: "CPUCreditUsage"},
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
		p.LabelPrefix + ".MaximumUsedTransactionIDs": mp.Graphs{
			Label: p.LabelPrefix + " Maximum Used Transaction IDs",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "MaximumUsedTransactionIDs", Label: "MaximumUsedTransactionIDs"},
			},
		},
		p.LabelPrefix + ".NetworkThroughput": mp.Graphs{
			Label: p.LabelPrefix + " Network Throughput",
			Unit:  "bytes/sec",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "NetworkTransmitThroughput", Label: "Transmit"},
				mp.Metrics{Name: "NetworkReceiveThroughput", Label: "Receive"},
			},
		},
		p.LabelPrefix + ".DiskQueueDepth": mp.Graphs{
			Label: p.LabelPrefix + " Disk Queue Depth",
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "DiskQueueDepth", Label: "DiskQueueDepth"},
			},
		},
		p.LabelPrefix + ".SwapUsage": mp.Graphs{
			Label: p.LabelPrefix + " Swap Usage",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "SwapUsage", Label: "SwapUsage"},
			},
		},
		p.LabelPrefix + ".OldestReplicationSlotLag": mp.Graphs{
			Label: p.LabelPrefix + " Oldest Replication Slot Lag",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "OldestReplicationSlotLag", Label: "OldestReplicationSlotLag"},
			},
		},
		p.LabelPrefix + ".ReplicationSlotDiskUsage": mp.Graphs{
			Label: p.LabelPrefix + " Replication Slot Disk Usage",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ReplicationSlotDiskUsage", Label: "ReplicationSlotDiskUsage"},
			},
		},
		p.LabelPrefix + ".TransactionLogsDiskUsage": mp.Graphs{
			Label: p.LabelPrefix + " Transaction Logs Disk Usage",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "TransactionLogsDiskUsage", Label: "TransactionLogsDiskUsage"},
			},
		},
		p.LabelPrefix + ".TransactionLogsGeneration": mp.Graphs{
			Label: p.LabelPrefix + " Transaction Logs Generation",
			Unit:  "bytes/sec",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "TransactionLogsGeneration", Label: "TransactionLogsGeneration"},
			},
		},
		p.LabelPrefix + ".IOPS": mp.Graphs{
			Label: p.LabelPrefix + " IOPS",
			Unit:  "iops",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ReadIOPS", Label: "Read"},
				mp.Metrics{Name: "WriteIOPS", Label: "Write"},
			},
		},
		p.LabelPrefix + ".Latency": mp.Graphs{
			Label: p.LabelPrefix + " Latency",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ReadLatency", Label: "Read"},
				mp.Metrics{Name: "WriteLatency", Label: "Write"},
			},
		},
		p.LabelPrefix + ".Throughput": mp.Graphs{
			Label: p.LabelPrefix + " Throughput",
			Unit:  "bytes/sec",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ReadThroughput", Label: "Read"},
				mp.Metrics{Name: "WriteThroughput", Label: "Write"},
			},
		},
	}
	return graphdef
}
