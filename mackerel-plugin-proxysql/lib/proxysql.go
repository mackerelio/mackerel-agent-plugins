package mpproxysql

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"database/sql"

	mp "github.com/mackerelio/go-mackerel-plugin"
	// Go driver for database/sql package
	_ "github.com/ziutek/mymysql/godrv"
)

func (m *ProxySQLPlugin) proxysqlGraphDef() map[string]mp.Graphs {
	labelPrefix := strings.Title(strings.Replace(m.MetricKeyPrefix(), "proxysql", "ProxySQL", -1))
	return map[string]mp.Graphs{
		"uptime": {
			Label: labelPrefix + " Uptime",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "proxysql_uptime", Label: "Seconds"},
			},
		},
		"traffic": {
			Label: labelPrefix + " Traffic",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "queries_backends_bytes_recv", Label: "Queries backends bytes received", Diff: true},
				{Name: "queries_backends_bytes_sent", Label: "Queries backends bytes sent", Diff: true},
				{Name: "queries_frontends_bytes_recv", Label: "Queries frontends bytes received", Diff: true},
				{Name: "queries_frontends_bytes_sent", Label: "Queries frontends bytes sent", Diff: true},
			},
		},
		"queries.": {
			Label: labelPrefix + " Queries",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "active_transactions", Label: "Active Transactions"},
				{Name: "questions", Label: "Questions", Diff: true},
				{Name: "slow_queries", Label: "Slow queries", Diff: true},
			},
		},
		"commands": {
			Label: labelPrefix + " Commands",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "com_autocommit", Label: "Command Autocommit", Diff: true},
				{Name: "com_autocommit_filtered", Label: "Command Autocommit Filtered", Diff: true},
				{Name: "com_commit", Label: "Command Commit", Diff: true},
				{Name: "com_commit_filtered", Label: "Command Commit Filtered", Diff: true},
				{Name: "com_rollback", Label: "Command Rollback", Diff: true},
				{Name: "com_rollback_filtered", Label: "Command Rollback Filtered", Diff: true},
				{Name: "com_frontend_stmt_prepare", Label: "Command Frontend Statement Prepare", Diff: true},
				{Name: "com_frontend_stmt_execute", Label: "Command Frontent Statement Execute", Diff: true},
				{Name: "com_frontend_stmt_close", Label: "Command Frontend Statement Close", Diff: true},
				{Name: "com_backend_stmt_prepare", Label: "Command Backend Statement Prepare", Diff: true},
				{Name: "com_backend_stmt_execute", Label: "Command Backend Statement Execute", Diff: true},
				{Name: "com_backend_stmt_close", Label: "Command Backend Statement Clode", Diff: true},
			},
		},
		"connections": {
			Label: labelPrefix + " Connections",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "client_connections_aborted", Label: "Client Connections Aborted", Diff: true},
				{Name: "client_connections_connected", Label: "Client Connections Connected", Diff: true},
				{Name: "client_connections_created", Label: "Client Connections Created", Diff: true},
				{Name: "client_connections_non_idle", Label: "Client Connections Non Idle", Diff: true},
				{Name: "server_connections_aborted", Label: "Server Connections Aborted", Diff: true},
				{Name: "server_connections_connected", Label: "Server Connections Connected", Diff: true},
				{Name: "server_connections_created", Label: "Server Connections Created", Diff: true},
				{Name: "server_connections_delayed", Label: "Server Connections Delayed", Diff: true},
			},
		},
		"workers": {
			Label: labelPrefix + " Workers",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "mysql_thread_workers", Label: "Thread Workers", Diff: true},
				{Name: "mysql_monitor_workers", Label: "Monitors Workers", Diff: true},
			},
		},
		"memories": {
			Label: labelPrefix + " Memory Bytes",
			Unit:  mp.UnitBytes,
			Metrics: []mp.Metrics{
				{Name: "sqlite3_memory_bytes", Label: "SQLite3 Memory Bytes", Stacked: true},
				{Name: "connpool_memory_bytes", Label: "Connection Pool Memory Bytes", Stacked: true},
				{Name: "query_cache_memory_bytes", Label: "Query Cache Memory Bytes", Stacked: true},
			},
		},
		"statement": {
			Label: labelPrefix + " Statement",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "stmt_client_active_total", Label: "Statement Client Active Total"},
				{Name: "stmt_client_active_unique", Label: "Statement Client Active Unique"},
				{Name: "stmt_server_active_total", Label: "Statement Server Active Total"},
				{Name: "stmt_server_active_unique", Label: "Statement Server Active Unique"},
				{Name: "stmt_cached", Label: "Statement Cached"},
			},
		},
		"querycache.count": {
			Label: labelPrefix + " Query Cache Counts",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "query_cache_count_get", Label: "Query Cache Count Get"},
				{Name: "query_cache_count_get_ok", Label: "Query Cache Get OK", Stacked: true},
				{Name: "query_cache_count_set", Label: "Query Cache Set", Stacked: true},
			},
		},
		"querycache.bytes": {
			Label: labelPrefix + " Query Cache Bytes In/Out",
			Unit:  mp.UnitBytes,
			Metrics: []mp.Metrics{
				{Name: "query_cache_bytes_in", Label: "Query Cache Bytes In", Diff: true},
				{Name: "query_cache_bytes_out", Label: "Query Cache Bytes Out", Diff: true},
			},
		},
		"querycache.size": {
			Label: labelPrefix + " Query Size",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "query_cache_purged", Label: "Query Cache Purged", Diff: true},
				{Name: "query_cache_entries", Label: "Query Cache Entries", Diff: true},
			},
		},
	}
}

// ProxySQLPlugin mackerel plugin for ProxySQL
type ProxySQLPlugin struct {
	URI      string
	Tempfile string
	prefix   string
}

// MetricKeyPrefix interface for PluginWithPrefix
func (m *ProxySQLPlugin) MetricKeyPrefix() string {
	if m.prefix == "" {
		m.prefix = "proxysql"
	}
	return m.prefix
}

func (m *ProxySQLPlugin) parseStatsMySQLGlobal(rows *sql.Rows, stat map[string]float64) error {
	for rows.Next() {
		var name, value string
		err := rows.Scan(&name, &value)
		if err != nil {
			return fmt.Errorf("parseStatsMySQLGlobal (stats_mysql_global): %s", err)
		}

		stat[strings.ToLower(name)], err = strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("%s", err)
		}
	}
	return rows.Err()
}

func (m *ProxySQLPlugin) fetchStatsMySQLGlobal(db *sql.DB, stat map[string]float64) error {
	rows, err := db.Query("SELECT Variable_Name, Variable_Value FROM stats_mysql_global")
	if err != nil {
		return fmt.Errorf("FetchMetrics (stats_mysql_global): %s", err)
	}

	return m.parseStatsMySQLGlobal(rows, stat)
}

// FetchMetrics interface for mackerelplugin
func (m *ProxySQLPlugin) FetchMetrics() (map[string]float64, error) {
	db, err := sql.Open("mymysql", m.URI)
	if err != nil {
		return nil, fmt.Errorf("FetchMetrics (DB Connect): %s", err)
	}
	defer db.Close()

	stat := make(map[string]float64)

	err = m.fetchStatsMySQLGlobal(db, stat)
	if err != nil {
		return nil, err
	}

	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (m *ProxySQLPlugin) GraphDefinition() map[string]mp.Graphs {
	return m.proxysqlGraphDef()
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "0.0.0.0", "Hostname")
	optPort := flag.String("port", "6032", "Port")
	optSocket := flag.String("socket", "", "Path to unix socket")
	optUser := flag.String("username", "radminuser", "Username")
	optPass := flag.String("password", os.Getenv("PROXYSQL_PASSWORD"), "Password")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optMetricKeyPrefix := flag.String("metric-key-prefix", "proxysql", "metric key prefix")
	flag.Parse()

	var proxysql ProxySQLPlugin

	if *optSocket != "" {
		proxysql.URI = fmt.Sprintf("unix:%s*stats/%s/%s", *optSocket, *optUser, *optPass)
	} else {
		proxysql.URI = fmt.Sprintf("tcp:%s:%s*stats/%s/%s", *optHost, *optPort, *optUser, *optPass)
	}

	proxysql.prefix = *optMetricKeyPrefix
	helper := mp.NewMackerelPlugin(&proxysql)
	helper.Tempfile = *optTempfile

	helper.Run()
}
