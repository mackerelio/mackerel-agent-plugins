package mpproxysql

import (
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGraphDefinition(t *testing.T) {
	var proxysql ProxySQLPlugin
	expect := 7

	graphdef := proxysql.GraphDefinition()
	if len(graphdef) != expect {
		t.Errorf("GetTempfilename: %d should be %d", len(graphdef), expect)
	}
}

func TestParseStatsMysqlGlobal(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Variable_Name", "Variable_Value"}).
		AddRow("ProxySQL_Uptime", "3600").
		AddRow("Active_Transactions", "0").
		AddRow("Client_Connections_aborted", "2").
		AddRow("Client_Connections_connected", "6").
		AddRow("Client_Connections_created", "6").
		AddRow("Server_Connections_aborted", "0").
		AddRow("Server_Connections_connected", "7").
		AddRow("Server_Connections_created", "7").
		AddRow("Server_Connections_delayed", "0").
		AddRow("Client_Connections_non_idle", "0").
		AddRow("Questions", "38").
		AddRow("Slow_queries", "0").
		AddRow("MySQL_Thread_Workers", "4").
		AddRow("MySQL_Monitor_Workers", "8").
		AddRow("Query_Cache_Memory_bytes", "0").
		AddRow("Query_Cache_count_GET", "0").
		AddRow("Query_Cache_count_GET_OK", "0").
		AddRow("Query_Cache_count_SET", "0").
		AddRow("Query_Cache_bytes_IN", "0").
		AddRow("Query_Cache_bytes_OUT", "0").
		AddRow("Query_Cache_Purged", "0").
		AddRow("Query_Cache_Entries", "0")

	query := "SELECT Variable_Name, Variable_Value FROM stats_mysql_global"
	mock.ExpectQuery(query).WillReturnRows(rows)

	result := make(map[string]float64)
	var proxysql ProxySQLPlugin
	err = proxysql.fetchStatsMySQLGlobal(db, result)
	if err != nil {
		t.Errorf("Failed to parse: %s", err)
	}

	expect := map[string]float64{
		"proxysql_uptime":              3600,
		"active_transactions":          0,
		"client_connections_aborted":   2,
		"client_connections_connected": 6,
		"client_connections_created":   6,
		"server_connections_aborted":   0,
		"server_connections_connected": 7,
		"server_connections_created":   7,
		"server_connections_delayed":   0,
		"client_connections_non_idle":  0,
		"questions":                    38,
		"slow_queries":                 0,
		"mysql_thread_workers":         4,
		"mysql_monitor_workers":        8,
		"query_cache_memory_bytes":     0,
		"query_cache_count_get":        0,
		"query_cache_count_get_ok":     0,
		"query_cache_count_set":        0,
		"query_cache_bytes_in":         0,
		"query_cache_bytes_out":        0,
		"query_cache_purged":           0,
		"query_cache_entries":          0,
	}

	for k := range expect {
		if expect[k] != result[k] {
			t.Errorf("%s does not match\nExpect: %v\nResult: %v", k, expect[k], result[k])
		}
	}
}
