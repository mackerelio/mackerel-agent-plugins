package mpaws-rds

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

func mergeGraphDefs(a, b map[string]mp.Graphs) map[string]mp.Graphs {
	for k, v := range b {
		a[k] = v
	}
	return a
}

func (p RDSPlugin) mySQLGraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		p.Prefix + ".CPUCreditBalance": {
			Label: p.LabelPrefix + " CPU CreditBalance",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "CPUCreditBalance", Label: "CPUCreditBalance"},
			},
		},
		p.Prefix + ".CPUCreditUsage": {
			Label: p.LabelPrefix + " CPU CreditUsage",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "CPUCreditUsage", Label: "CPUCreditUsage"},
			},
		},
	}
}

func (p RDSPlugin) auroraGraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		p.Prefix + ".Deadlocks": {
			Label: p.LabelPrefix + " Dead Locks",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "Deadlocks", Label: "Deadlocks"},
			},
		},
		p.Prefix + ".Transaction": {
			Label: p.LabelPrefix + " Transaction",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "ActiveTransactions", Label: "Active"},
				{Name: "BlockedTransactions", Label: "Blocked"},
			},
		},
		p.Prefix + ".EngineUptime": {
			Label: p.LabelPrefix + " Engine Uptime",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "EngineUptime", Label: "EngineUptime"},
			},
		},
		p.Prefix + ".Queries": {
			Label: p.LabelPrefix + " Queries",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "Queries", Label: "Queries"},
			},
		},
		p.Prefix + ".LoginFailures": {
			Label: p.LabelPrefix + " Login Failures",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "LoginFailures", Label: "LoginFailures"},
			},
		},
		p.Prefix + ".CacheHitRatio": {
			Label: p.LabelPrefix + " Cache Hit Ratio",
			Unit:  "percentage",
			Metrics: []mp.Metrics{
				{Name: "ResultSetCacheHitRatio", Label: "ResultSet"},
				{Name: "BufferCacheHitRatio", Label: "Buffer"},
			},
		},
		p.Prefix + ".AuroraBinlogReplicaLag": {
			Label: p.LabelPrefix + " Aurora Binlog ReplicaLag",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "AuroraBinlogReplicaLag", Label: "AuroraBinlogReplicaLag"},
			},
		},
		p.Prefix + ".AuroraReplicaLag": {
			Label: p.LabelPrefix + " Aurora ReplicaLag",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "AuroraReplicaLagMaximum", Label: "ReplicaLagMaximum"},
				{Name: "AuroraBinlogReplicaLag", Label: "ReplicaLagMinimum"},
			},
		},
	}
}

func (p RDSPlugin) postgreSQLGraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		p.Prefix + ".CPUCreditBalance": {
			Label: p.LabelPrefix + " CPU CreditBalance",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "CPUCreditBalance", Label: "CPUCreditBalance"},
			},
		},
		p.Prefix + ".CPUCreditUsage": {
			Label: p.LabelPrefix + " CPU CreditUsage",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "CPUCreditUsage", Label: "CPUCreditUsage"},
			},
		},
		p.Prefix + ".MaximumUsedTransactionIDs": {
			Label: p.LabelPrefix + " Maximum Used Transaction IDs",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "MaximumUsedTransactionIDs", Label: "MaximumUsedTransactionIDs"},
			},
		},
		p.Prefix + ".OldestReplicationSlotLag": {
			Label: p.LabelPrefix + " Oldest Replication Slot Lag",
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "OldestReplicationSlotLag", Label: "OldestReplicationSlotLag"},
			},
		},
		p.Prefix + ".ReplicationSlotDiskUsage": {
			Label: p.LabelPrefix + " Replication Slot Disk Usage",
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "ReplicationSlotDiskUsage", Label: "ReplicationSlotDiskUsage"},
			},
		},
		p.Prefix + ".TransactionLogsDiskUsage": {
			Label: p.LabelPrefix + " Transaction Logs Disk Usage",
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "TransactionLogsDiskUsage", Label: "TransactionLogsDiskUsage"},
			},
		},
		p.Prefix + ".TransactionLogsGeneration": {
			Label: p.LabelPrefix + " Transaction Logs Generation",
			Unit:  "bytes/sec",
			Metrics: []mp.Metrics{
				{Name: "TransactionLogsGeneration", Label: "TransactionLogsGeneration"},
			},
		},
	}
}
