package main

import (
	mp "github.com/mackerelio/go-mackerel-plugin"
)

func (p RDSPlugin) rdsMetrics() (metrics []string) {
	for _, v := range p.GraphDefinition() {
		for _, vv := range v.Metrics {
			metrics = append(metrics, vv.Name)
		}
	}
	return
}

func mergeGraphDefs(a, b map[string](mp.Graphs)) map[string](mp.Graphs) {
	for k, v := range b {
		a[k] = v
	}
	return a
}

func (p RDSPlugin) mySQLGraphDefinition() map[string](mp.Graphs) {
	return map[string](mp.Graphs){
		p.Prefix + ".CPUCreditBalance": mp.Graphs{
			Label: p.LabelPrefix + " CPU CreditBalance",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "CPUCreditBalance", Label: "CPUCreditBalance"},
			},
		},
		p.Prefix + ".CPUCreditUsage": mp.Graphs{
			Label: p.LabelPrefix + " CPU CreditUsage",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "CPUCreditUsage", Label: "CPUCreditUsage"},
			},
		},
	}
}

func (p RDSPlugin) auroraGraphDefinition() map[string](mp.Graphs) {
	return map[string](mp.Graphs){
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
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "AuroraBinlogReplicaLag", Label: "AuroraBinlogReplicaLag"},
			},
		},
		p.Prefix + ".AuroraReplicaLag": mp.Graphs{
			Label: p.LabelPrefix + " Aurora ReplicaLag",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "AuroraReplicaLagMaximum", Label: "ReplicaLagMaximum"},
				mp.Metrics{Name: "AuroraBinlogReplicaLag", Label: "ReplicaLagMinimum"},
			},
		},
	}
}

func (p RDSPlugin) postgreSQLGraphDefinition() map[string](mp.Graphs) {
	return map[string](mp.Graphs){
		p.Prefix + ".CPUCreditBalance": mp.Graphs{
			Label: p.LabelPrefix + " CPU CreditBalance",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "CPUCreditBalance", Label: "CPUCreditBalance"},
			},
		},
		p.Prefix + ".CPUCreditUsage": mp.Graphs{
			Label: p.LabelPrefix + " CPU CreditUsage",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "CPUCreditUsage", Label: "CPUCreditUsage"},
			},
		},
		p.Prefix + ".MaximumUsedTransactionIDs": mp.Graphs{
			Label: p.LabelPrefix + " Maximum Used Transaction IDs",
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "MaximumUsedTransactionIDs", Label: "MaximumUsedTransactionIDs"},
			},
		},
		p.Prefix + ".OldestReplicationSlotLag": mp.Graphs{
			Label: p.LabelPrefix + " Oldest Replication Slot Lag",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "OldestReplicationSlotLag", Label: "OldestReplicationSlotLag"},
			},
		},
		p.Prefix + ".ReplicationSlotDiskUsage": mp.Graphs{
			Label: p.LabelPrefix + " Replication Slot Disk Usage",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "ReplicationSlotDiskUsage", Label: "ReplicationSlotDiskUsage"},
			},
		},
		p.Prefix + ".TransactionLogsDiskUsage": mp.Graphs{
			Label: p.LabelPrefix + " Transaction Logs Disk Usage",
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "TransactionLogsDiskUsage", Label: "TransactionLogsDiskUsage"},
			},
		},
		p.Prefix + ".TransactionLogsGeneration": mp.Graphs{
			Label: p.LabelPrefix + " Transaction Logs Generation",
			Unit:  "bytes/sec",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "TransactionLogsGeneration", Label: "TransactionLogsGeneration"},
			},
		},
	}
}
