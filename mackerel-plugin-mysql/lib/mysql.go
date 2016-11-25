package mpmysql

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/ziutek/mymysql/mysql"
	// MySQL Driver
	_ "github.com/ziutek/mymysql/native"
)

var (
	processState map[string]bool
)

func init() {
	processState = make(map[string]bool, 0)
	processState["State_closing_tables"] = true
	processState["State_copying_to_tmp_table"] = true
	processState["State_end"] = true
	processState["State_freeing_items"] = true
	processState["State_init"] = true
	processState["State_locked"] = true
	processState["State_login"] = true
	processState["State_preparing"] = true
	processState["State_reading_from_net"] = true
	processState["State_sending_data"] = true
	processState["State_sorting_result"] = true
	processState["State_statistics"] = true
	processState["State_updating"] = true
	processState["State_writing_to_net"] = true
	processState["State_none"] = true
	processState["State_other"] = true
}

func (m MySQLPlugin) defaultGraphdef() map[string]mp.Graphs {
	labelPrefix := strings.Title(strings.Replace(m.MetricKeyPrefix(), "mysql", "MySQL", -1))

	return map[string]mp.Graphs{
		"cmd": {
			Label: labelPrefix + " Command",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "Com_insert", Label: "Insert", Diff: true, Stacked: true, Type: "uint64"},
				{Name: "Com_select", Label: "Select", Diff: true, Stacked: true, Type: "uint64"},
				{Name: "Com_update", Label: "Update", Diff: true, Stacked: true, Type: "uint64"},
				{Name: "Com_update_multi", Label: "Update Multi", Diff: true, Stacked: true, Type: "uint64"},
				{Name: "Com_delete", Label: "Delete", Diff: true, Stacked: true, Type: "uint64"},
				{Name: "Com_delete_multi", Label: "Delete Multi", Diff: true, Stacked: true, Type: "uint64"},
				{Name: "Com_replace", Label: "Replace", Diff: true, Stacked: true, Type: "uint64"},
				{Name: "Com_set_option", Label: "Set Option", Diff: true, Stacked: true, Type: "uint64"},
				{Name: "Qcache_hits", Label: "Query Cache Hits", Diff: true, Stacked: false, Type: "uint64"},
				{Name: "Questions", Label: "Questions", Diff: true, Stacked: false, Type: "uint64"},
			},
		},
		"join": {
			Label: labelPrefix + " Join/Scan",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "Select_full_join", Label: "Select Full JOIN", Diff: true, Stacked: false},
				{Name: "Select_full_range_join", Label: "Select Full Range JOIN", Diff: true, Stacked: false},
				{Name: "Select_scan", Label: "Select SCAN", Diff: true, Stacked: false},
				{Name: "Sort_scan", Label: "Sort SCAN", Diff: true, Stacked: false},
			},
		},
		"threads": {
			Label: labelPrefix + " Threads",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "Max_used_connections", Label: "Max used connections", Diff: false, Stacked: false},
				{Name: "Threads_connected", Label: "Connected", Diff: false, Stacked: false},
				{Name: "Threads_running", Label: "Running", Diff: false, Stacked: false},
				{Name: "Threads_cached", Label: "Cached", Diff: false, Stacked: false},
			},
		},
		"connections": {
			Label: labelPrefix + " Connections",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "Connections", Label: "Connections", Diff: true, Stacked: false},
				{Name: "Threads_created", Label: "Created Threads", Diff: true, Stacked: false},
				{Name: "Aborted_clients", Label: "Aborted Clients", Diff: true, Stacked: false},
				{Name: "Aborted_connects", Label: "Aborted Connects", Diff: true, Stacked: false},
			},
		},
		"seconds_behind_master": {
			Label: labelPrefix + " Slave status",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "Seconds_Behind_Master", Label: "Seconds Behind Master", Diff: false, Stacked: false},
			},
		},
		"table_locks": {
			Label: labelPrefix + " Table Locks/Slow Queries",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "Table_locks_immediate", Label: "Table Locks Immediate", Diff: true, Stacked: false},
				{Name: "Table_locks_waited", Label: "Table Locks Waited", Diff: true, Stacked: false},
				{Name: "Slow_queries", Label: "Slow Queries", Diff: true, Stacked: false},
			},
		},
		"traffic": {
			Label: labelPrefix + " Traffic",
			Unit:  "bytes/sec",
			Metrics: []mp.Metrics{
				{Name: "Bytes_sent", Label: "Sent Bytes", Diff: true, Stacked: false},
				{Name: "Bytes_received", Label: "Received Bytes", Diff: true, Stacked: false},
			},
		},
		"capacity": {
			Label: labelPrefix + " Capacity",
			Unit:  "percentage",
			Metrics: []mp.Metrics{
				{Name: "PercentageOfConnections", Label: "Percentage Of Connections", Diff: false, Stacked: false},
				{Name: "PercentageOfBufferPool", Label: "Percentage Of Buffer Pool", Diff: false, Stacked: false},
			},
		},
	}
}

// MySQLPlugin mackerel plugin for MySQL
type MySQLPlugin struct {
	Target         string
	Tempfile       string
	prefix         string
	Username       string
	Password       string
	DisableInnoDB  bool
	isUnixSocket   bool
	EnableExtended bool
}

// MetricKeyPrefix retruns the metrics key prefix
func (m MySQLPlugin) MetricKeyPrefix() string {
	if m.prefix == "" {
		m.prefix = "mysql"
	}
	return m.prefix
}

func (m MySQLPlugin) fetchShowStatus(db mysql.Conn, stat map[string]float64) error {
	rows, _, err := db.Query("show /*!50002 global */ status")
	if err != nil {
		log.Fatalln("FetchMetrics (Status): ", err)
		return err
	}

	for _, row := range rows {
		if len(row) > 1 {
			variableName := string(row[0].([]byte))
			if err != nil {
				log.Fatalln("FetchMetrics (Status Fetch): ", err)
				return err
			}
			stat[variableName], _ = atof(string(row[1].([]byte)))
		} else {
			log.Fatalln("FetchMetrics (Status): row length is too small: ", len(row))
		}
	}
	if m.EnableExtended {
		err = fetchShowStatusBackwardCompatibile(stat)
		if err != nil {
			log.Fatalln("FetchExtendedMetrics (Status Fetch): ", err)
		}
	}
	return nil
}

func fetchShowStatusBackwardCompatibile(stat map[string]float64) error {
	// https://github.com/percona/percona-monitoring-plugins/blob/v1.1.6/cacti/scripts/ss_get_mysql_stats.php#L585
	if val, found := stat["table_open_cache"]; found {
		stat["table_cache"] = val
	}
	return nil
}

func (m MySQLPlugin) fetchShowInnodbStatus(db mysql.Conn, stat map[string]float64) error {
	row, _, err := db.QueryFirst("SHOW /*!50000 ENGINE*/ INNODB STATUS")
	if err != nil {
		log.Fatalln("FetchMetrics (InnoDB Status): ", err)
		return err
	}

	if len(row) > 0 {
		parseInnodbStatus(string(row[len(row)-1].([]byte)), &stat)
	} else {
		log.Fatalln("FetchMetrics (InnoDB Status): row length is too small: ", len(row))
	}
	return nil
}

func (m MySQLPlugin) fetchShowVariables(db mysql.Conn, stat map[string]float64) error {
	rows, _, err := db.Query("SHOW VARIABLES")
	if err != nil {
		log.Fatalln("FetchMetrics (Variables): ", err)
	}

	for _, row := range rows {
		if len(row) > 1 {
			variableName := string(row[0].([]byte))
			if err != nil {
				log.Println("FetchMetrics (Fetch Variables): ", err)
			}
			stat[variableName], _ = atof(string(row[1].([]byte)))
		} else {
			log.Fatalln("FetchMetrics (Variables): row length is too small: ", len(row))
		}
	}
	return nil
}

func (m MySQLPlugin) fetchShowSlaveStatus(db mysql.Conn, stat map[string]float64) error {
	rows, res, err := db.Query("show slave status")
	if err != nil {
		log.Fatalln("FetchMetrics (Slave Status): ", err)
		return err
	}

	for _, row := range rows {
		idx := res.Map("Seconds_Behind_Master")
		switch row[idx].(type) {
		case nil:
			// nop
		default:
			Value := row.Int(idx)
			stat["Seconds_Behind_Master"] = float64(Value)
		}
	}
	return nil
}

func (m MySQLPlugin) fetchProcesslist(db mysql.Conn, stat map[string]float64) error {
	rows, _, err := db.Query("SHOW PROCESSLIST")
	if err != nil {
		log.Fatalln("FetchMetrics (Processlist): ", err)
		return err
	}

	for k := range processState {
		stat[k] = 0
	}

	for _, row := range rows {
		if len(row) > 1 {
			var state string
			if row[6] == nil {
				state = "NULL"
			} else {
				state = string(row[6].([]byte))
			}
			parseProcesslist(state, &stat)
		} else {
			log.Fatalln("FetchMetrics (Processlist): row length is too small: ", len(row))
		}
	}

	return nil
}

func (m MySQLPlugin) calculateCapacity(stat map[string]float64) {
	stat["PercentageOfConnections"] = 100.0 * stat["Threads_connected"] / stat["max_connections"]
	stat["PercentageOfBufferPool"] = 100.0 * stat["database_pages"] / stat["pool_size"]
}

// FetchMetrics interface for mackerelplugin
func (m MySQLPlugin) FetchMetrics() (map[string]interface{}, error) {
	proto := "tcp"
	if m.isUnixSocket {
		proto = "unix"
	}
	db := mysql.New(proto, "", m.Target, m.Username, m.Password, "")
	err := db.Connect()
	if err != nil {
		log.Fatalln("FetchMetrics (DB Connect): ", err)
		return nil, err
	}
	defer db.Close()

	stat := make(map[string]float64)
	m.fetchShowStatus(db, stat)

	if m.DisableInnoDB != true {
		m.fetchShowInnodbStatus(db, stat)
		m.fetchShowVariables(db, stat)
	}

	m.fetchShowSlaveStatus(db, stat)

	if m.EnableExtended {
		m.fetchProcesslist(db, stat)
	}

	m.calculateCapacity(stat)

	statRet := make(map[string]interface{})
	for key, value := range stat {
		statRet[key] = value
	}

	return statRet, err
}

// GraphDefinition interface for mackerelplugin
func (m MySQLPlugin) GraphDefinition() map[string]mp.Graphs {
	graphdef := m.defaultGraphdef()
	if !m.DisableInnoDB {
		graphdef = m.addGraphdefWithInnoDBMetrics(graphdef)
	}
	if m.EnableExtended {
		graphdef = m.addExtendedGraphdef(graphdef)
	}
	return graphdef
}

func (m MySQLPlugin) addGraphdefWithInnoDBMetrics(graphdef map[string]mp.Graphs) map[string]mp.Graphs {
	prefix := m.MetricKeyPrefix()
	graphdef["innodb_rows"] = mp.Graphs{
		Label: prefix + ".innodb Rows",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "Innodb_rows_read", Label: "Read", Diff: true, Stacked: false},
			{Name: "Innodb_rows_inserted", Label: "Inserted", Diff: true, Stacked: false},
			{Name: "Innodb_rows_updated", Label: "Updated", Diff: true, Stacked: false},
			{Name: "Innodb_rows_deleted", Label: "Deleted", Diff: true, Stacked: false},
		},
	}
	graphdef["innodb_row_lock_time"] = mp.Graphs{
		Label: prefix + ".innodb Row Lock Time",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "Innodb_row_lock_time", Label: "Lock Time", Diff: true, Stacked: false},
		},
	}
	graphdef["innodb_row_lock_waits"] = mp.Graphs{
		Label: prefix + ".innodb Row Lock Waits",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "Innodb_row_lock_waits", Label: "Lock Waits", Diff: true, Stacked: false},
		},
	}
	graphdef["innodb_adaptive_hash_index"] = mp.Graphs{
		Label: prefix + ".innodb Adaptive Hash Index",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "hash_index_cells_total", Label: "Hash Index Cells Total", Diff: false, Stacked: false},
			{Name: "hash_index_cells_used", Label: "Hash Index Cells Used", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_buffer_pool_read"] = mp.Graphs{
		Label: prefix + ".innodb Buffer Pool Read (/sec)",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "read_ahead", Label: "Pages Read Ahead", Diff: false, Stacked: false},
			{Name: "read_evicted", Label: "Evicted Without Access", Diff: false, Stacked: false},
			{Name: "read_random_ahead", Label: "Random Read Ahead", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_buffer_pool_activity"] = mp.Graphs{
		Label: prefix + ".innodb Buffer Pool Activity (Pages)",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "pages_created", Label: "Created", Diff: true, Stacked: false},
			{Name: "pages_read", Label: "Read", Diff: true, Stacked: false},
			{Name: "pages_written", Label: "Written", Diff: true, Stacked: false},
		},
	}
	graphdef["innodb_buffer_pool_efficiency"] = mp.Graphs{
		Label: prefix + ".innodb Buffer Pool Efficiency",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "Innodb_buffer_pool_reads", Label: "Reads", Diff: true, Stacked: false},
			{Name: "Innodb_buffer_pool_read_requests", Label: "Read Requests", Diff: true, Stacked: false},
		},
	}
	graphdef["innodb_buffer_pool"] = mp.Graphs{
		Label: prefix + ".innodb Buffer Pool (Pages)",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "pool_size", Label: "Pool Size", Diff: false, Stacked: false},
			{Name: "database_pages", Label: "Used", Diff: false, Stacked: true},
			{Name: "free_pages", Label: "Free", Diff: false, Stacked: true},
			{Name: "modified_pages", Label: "Modified", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_checkpoint_age"] = mp.Graphs{
		Label: prefix + ".innodb Checkpoint Age",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "uncheckpointed_bytes", Label: "Uncheckpointed", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_current_lock_waits"] = mp.Graphs{
		Label: prefix + ".innodb Current Lock Waits (secs)",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "innodb_lock_wait_secs", Label: "Innodb Lock Wait", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_io"] = mp.Graphs{
		Label: prefix + ".innodb I/O",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "file_reads", Label: "File Reads", Diff: true, Stacked: false},
			{Name: "file_writes", Label: "File Writes", Diff: true, Stacked: false},
			{Name: "file_fsyncs", Label: "File fsyncs", Diff: true, Stacked: false},
			{Name: "log_writes", Label: "Log Writes", Diff: true, Stacked: false},
		},
	}
	graphdef["innodb_io_pending"] = mp.Graphs{
		Label: prefix + ".innodb I/O Pending",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "pending_normal_aio_reads", Label: "Normal AIO Reads", Diff: false, Stacked: false},
			{Name: "pending_normal_aio_writes", Label: "Normal AIO Writes", Diff: false, Stacked: false},
			{Name: "pending_ibuf_aio_reads", Label: "InnoDB Buffer AIO Reads", Diff: false, Stacked: false},
			{Name: "pending_aio_log_ios", Label: "AIO Log IOs", Diff: false, Stacked: false},
			{Name: "pending_aio_sync_ios", Label: "AIO Sync IOs", Diff: false, Stacked: false},
			{Name: "pending_log_flushes", Label: "Log Flushes", Diff: false, Stacked: false},
			{Name: "pending_buf_pool_flushes", Label: "Buffer Pool Flushes", Diff: false, Stacked: false},
			{Name: "pending_log_writes", Label: "Log Writes", Diff: false, Stacked: false},
			{Name: "pending_chkp_writes", Label: "Checkpoint Writes", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_insert_buffer"] = mp.Graphs{
		Label: prefix + ".innodb Insert Buffer",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "ibuf_inserts", Label: "Inserts", Diff: true, Stacked: false},
			{Name: "ibuf_merges", Label: "Merges", Diff: true, Stacked: false},
			{Name: "ibuf_merged", Label: "Merged", Diff: true, Stacked: false},
		},
	}
	graphdef["innodb_insert_buffer_usage"] = mp.Graphs{
		Label: prefix + ".innodb Insert Buffer Usage (Cells)",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "ibuf_cell_count", Label: "Cell Count", Diff: false, Stacked: false},
			{Name: "ibuf_used_cells", Label: "Used", Diff: false, Stacked: true},
			{Name: "ibuf_free_cells", Label: "Free", Diff: false, Stacked: true},
		},
	}
	graphdef["innodb_lock_structures"] = mp.Graphs{
		Label: prefix + ".innodb Lock Structures",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "innodb_lock_structs", Label: "Structures", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_log"] = mp.Graphs{
		Label: prefix + ".innodb Log",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "log_bytes_written", Label: "Written", Diff: true, Stacked: false},
			{Name: "log_bytes_flushed", Label: "Flushed", Diff: true, Stacked: false},
			{Name: "unflushed_log", Label: "Unflushed", Diff: false, Stacked: false},
			{Name: "innodb_log_buffer_size", Label: "Buffer Size", Diff: false, Stacked: true},
		},
	}
	graphdef["innodb_memory_allocation"] = mp.Graphs{
		Label: prefix + ".innodb Memory Allocation",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "additional_pool_alloc", Label: "Additional Pool Allocated", Diff: false, Stacked: false},
			{Name: "total_mem_alloc", Label: "Total Memory Allocated", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_semaphores"] = mp.Graphs{
		Label: prefix + ".innodb Semaphores",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "spin_waits", Label: "Spin Waits", Diff: true, Stacked: false},
			{Name: "spin_rounds", Label: "Spin Rounds", Diff: true, Stacked: false},
			{Name: "os_waits", Label: "OS Waits", Diff: true, Stacked: false},
		},
	}
	graphdef["innodb_tables_in_use"] = mp.Graphs{
		Label: prefix + ".innodb Tables In Use",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "innodb_locked_tables", Label: "Table in Use", Diff: false, Stacked: false},
			{Name: "innodb_tables_in_use", Label: "Locked Tables", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_transactions_active_locked"] = mp.Graphs{
		Label: prefix + ".innodb Transactions Active/Locked",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "current_transactions", Label: "Current", Diff: false, Stacked: false},
			{Name: "active_transactions", Label: "Active", Diff: false, Stacked: false},
			{Name: "locked_transactions", Label: "Locked", Diff: false, Stacked: false},
			{Name: "read_views", Label: "Read Views", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_transactions"] = mp.Graphs{
		Label: prefix + ".innodb Transactions",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "history_list", Label: "History List", Diff: false, Stacked: false},
			{Name: "innodb_transactions", Label: "InnoDB Transactions", Diff: true, Stacked: false},
		},
	}
	return graphdef
}

func (m MySQLPlugin) addExtendedGraphdef(graphdef map[string]mp.Graphs) map[string]mp.Graphs {
	//TODO
	prefix := m.MetricKeyPrefix()
	graphdef["query_cache"] = mp.Graphs{
		Label: prefix + ".query Cache",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "Qcache_queries_in_cache", Label: "Qcache Queries In Cache", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "Qcache_hits", Label: "Qcache Hits", Diff: true, Stacked: false, Type: "uint64"},
			{Name: "Qcache_inserts", Label: "Qcache Inserts", Diff: true, Stacked: false, Type: "uint64"},
			{Name: "Qcache_not_cached", Label: "Qcache Not Cached", Diff: true, Stacked: false, Type: "uint64"},
			{Name: "Qcache_lowmem_prunes", Label: "Qcache Lowmem Prunes", Diff: true, Stacked: false, Type: "uint64"},
		},
	}
	graphdef["query_cache_memory"] = mp.Graphs{
		Label: prefix + ".query Cache Memory",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "query_cache_size", Label: "Query Cache Size", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "Qcache_free_memory", Label: "Qcache Free Memory", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "Qcache_total_blocks", Label: "Qcache Total Blocks", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "Qcache_free_blocks", Label: "Qcache Free Blocks", Diff: false, Stacked: false, Type: "uint64"},
		},
	}
	graphdef["temporary_objects"] = mp.Graphs{
		Label: prefix + ".temporary Objects",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "Created_tmp_tables", Label: "Created Tmp Tables", Diff: true, Stacked: false, Type: "uint64"},
			{Name: "Created_tmp_disk_tables", Label: "Created Tmp Disk Tables", Diff: true, Stacked: false, Type: "uint64"},
			{Name: "Created_tmp_files", Label: "Created Tmp Files", Diff: true, Stacked: false, Type: "uint64"},
		},
	}
	graphdef["files_and_tables"] = mp.Graphs{
		Label: prefix + ".files and Tables",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "table_cache", Label: "Table Cache", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "Open_tables", Label: "Open Tables", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "Open_files", Label: "Open Files", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "Opened_tables", Label: "Opened Tables", Diff: false, Stacked: false, Type: "uint64"},
		},
	}
	graphdef["processlist"] = mp.Graphs{
		Label: prefix + ".processlist",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "State_closing_tables", Label: "State Closing Tables", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "State_copying_to_tmp_table", Label: "State Copying To Tmp Table", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "State_end", Label: "State End", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "State_freeing_items", Label: "State Freeing Items", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "State_init", Label: "State Init", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "State_locked", Label: "State Locked", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "State_login", Label: "State Login", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "State_preparing", Label: "State Preparing", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "State_reading_from_net", Label: "State Reading From Net", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "State_sending_data", Label: "State Sending Data", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "State_sorting_result", Label: "State Sorting Result", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "State_statistics", Label: "State Statistics", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "State_updating", Label: "State Updating", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "State_writing_to_net", Label: "State Writing To Net", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "State_none", Label: "State None", Diff: false, Stacked: false, Type: "uint64"},
			{Name: "State_other", Label: "State Other", Diff: false, Stacked: false, Type: "uint64"},
		},
	}
	graphdef["sorts"] = mp.Graphs{
		Label: prefix + ".sorts",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "Sort_rows", Label: "Sort Rows", Diff: true, Stacked: false, Type: "uint64"},
			{Name: "Sort_range", Label: "Sort Range", Diff: true, Stacked: false, Type: "uint64"},
			{Name: "Sort_merge_passes", Label: "Sort Merge Passes", Diff: true, Stacked: false, Type: "uint64"},
			{Name: "Sort_scan", Label: "Sort Scan", Diff: true, Stacked: false, Type: "uint64"},
		},
	}

	return graphdef
}

func setIfEmpty(p *map[string]float64, key string, val float64) {
	_, ok := (*p)[key]
	if !ok {
		(*p)[key] = val
	}
}

func parseInnodbStatus(str string, p *map[string]float64) {
	isTransaction := false
	prevLine := ""

	for _, line := range strings.Split(str, "\n") {
		record := strings.Fields(line)

		// Innodb Semaphores
		if strings.Index(line, "Mutex spin waits") == 0 {
			increaseMap(p, "spin_waits", record[3])
			increaseMap(p, "spin_rounds", record[5])
			increaseMap(p, "os_waits", record[8])
			continue
		}
		if strings.Index(line, "RW-shared spins") == 0 && strings.Index(line, ";") > 0 {
			// 5.5, 5.6
			increaseMap(p, "spin_waits", record[2])
			increaseMap(p, "os_waits", record[5])
			increaseMap(p, "spin_waits", record[8])
			increaseMap(p, "os_waits", record[11])
			continue
		}
		if strings.Index(line, "RW-shared spins") == 0 && strings.Index(line, "; RW-excl spins") < 0 {
			// 5.1
			increaseMap(p, "spin_waits", record[2])
			increaseMap(p, "os_waits", record[7])
			continue
		}
		if strings.Index(line, "RW-excl spins") == 0 {
			// 5.5, 5.6
			increaseMap(p, "spin_waits", record[2])
			increaseMap(p, "os_waits", record[7])
			continue
		}
		if strings.Index(line, "seconds the semaphore:") > 0 {
			increaseMap(p, "innodb_sem_waits", "1")
			wait, _ := atof(record[9])
			wait = wait * 1000
			increaseMap(p, "innodb_sem_wait_time_ms", fmt.Sprintf("%.f", wait))
			continue
		}

		// Innodb Transactions
		if strings.Index(line, "Trx id counter") == 0 {
			loVal := ""
			if len(record) >= 5 {
				loVal = record[4]
			}
			val := makeBigint(record[3], loVal)
			increaseMap(p, "innodb_transactions", fmt.Sprintf("%d", val))
			isTransaction = true
			continue
		}
		if strings.Index(line, "Purge done for trx") == 0 {
			if record[7] == "undo" {
				record[7] = ""
			}
			val := makeBigint(record[6], record[7])
			trx := (*p)["innodb_transactions"] - float64(val)
			increaseMap(p, "unpurged_txns", fmt.Sprintf("%.f", trx))
			continue
		}
		if strings.Index(line, "History list length") == 0 {
			increaseMap(p, "history_list", record[3])
			continue
		}
		if isTransaction && strings.Index(line, "---TRANSACTION") == 0 {
			increaseMap(p, "current_transactions", "1")
			if strings.Index(line, "ACTIVE") > 0 {
				increaseMap(p, "active_transactions", "1")
			}
			continue
		}
		if isTransaction && strings.Index(line, "------- TRX HAS BEEN") == 0 {
			increaseMap(p, "innodb_lock_wait_secs", "1")
			continue
		}
		if strings.Index(line, "read views open inside InnoDB") > 0 {
			(*p)["read_views"], _ = atof(record[0])
			continue
		}
		if strings.Index(line, "------- TRX HAS BEEN") == 0 {
			increaseMap(p, "innodb_tables_in_use", record[4])
			increaseMap(p, "innodb_locked_tables", record[6])
			continue
		}
		if isTransaction && strings.Index(line, "lock struct(s)") == 0 {
			if strings.Index(line, "LOCK WAIT") > 0 {
				increaseMap(p, "innodb_lock_structs", record[2])
				increaseMap(p, "locked_transactions", "1")
			} else {
				increaseMap(p, "innodb_lock_structs", record[0])
			}
			continue
		}

		// File I/O
		if strings.Index(line, " OS file reads, ") > 0 {
			(*p)["file_reads"], _ = atof(record[0])
			(*p)["file_writes"], _ = atof(record[4])
			(*p)["file_fsyncs"], _ = atof(record[8])
			continue
		}
		if strings.Index(line, "Pending normal aio reads:") == 0 {
			(*p)["pending_normal_aio_reads"], _ = atof(record[4])
			(*p)["pending_normal_aio_writes"], _ = atof(record[7])
			continue
		}
		if strings.Index(line, "ibuf aio reads") == 0 {
			(*p)["pending_ibuf_aio_reads"], _ = atof(record[3])
			(*p)["pending_aio_log_ios"], _ = atof(record[6])
			(*p)["pending_aio_sync_ios"], _ = atof(record[9])
			continue
		}
		if strings.Index(line, "Pending flushes (fsync)") == 0 {
			(*p)["pending_log_flushes"], _ = atof(record[4])
			(*p)["pending_buf_pool_flushes"], _ = atof(record[7])
			continue
		}

		// Insert Buffer and Adaptive Hash Index
		if strings.Index(line, "Ibuf for space 0: size ") == 0 {
			(*p)["ibuf_used_cells"], _ = atof(record[5])
			(*p)["ibuf_free_cells"], _ = atof(record[9])
			(*p)["ibuf_cell_count"], _ = atof(record[12])
			continue
		}
		if strings.Index(line, "Ibuf: size ") == 0 {
			(*p)["ibuf_used_cells"], _ = atof(record[2])
			(*p)["ibuf_free_cells"], _ = atof(record[6])
			(*p)["ibuf_cell_count"], _ = atof(record[9])
			if strings.Index(line, "merges") > 0 {
				(*p)["ibuf_merges"], _ = atof(record[10])
			}
			continue
		}
		if strings.Index(line, ", delete mark ") > 0 && strings.Index(prevLine, "merged operations:") == 0 {
			(*p)["ibuf_inserts"], _ = atof(record[1])
			v1, _ := atof(record[1])
			v2, _ := atof(record[4])
			v3, _ := atof(record[6])
			(*p)["ibuf_merged"] = v1 + v2 + v3
			continue
		}
		if strings.Index(line, " merged recs, ") > 0 {
			(*p)["ibuf_inserts"], _ = atof(record[0])
			(*p)["ibuf_merged"], _ = atof(record[2])
			(*p)["ibuf_merges"], _ = atof(record[5])
			continue
		}
		if strings.Index(line, "Hash table size ") == 0 {
			(*p)["hash_index_cells_total"], _ = atof(record[3])
			if strings.Index(line, "used cells") > 0 {
				(*p)["hash_index_cells_used"], _ = atof(record[6])
			} else {
				(*p)["hash_index_cells_used"] = 0
			}
			continue
		}

		// Log
		if strings.Index(line, " log i/o's done, ") > 0 {
			(*p)["log_writes"], _ = atof(record[0])
			continue
		}
		if strings.Index(line, " pending log writes, ") > 0 {
			(*p)["pending_log_writes"], _ = atof(record[0])
			(*p)["pending_chkp_writes"], _ = atof(record[4])
			continue
		}
		if strings.Index(line, "Log sequence number") == 0 {
			val, _ := atof(record[3])
			if len(record) >= 5 {
				val = float64(makeBigint(record[3], record[4]))
			}
			(*p)["log_bytes_written"] = val
			continue
		}
		if strings.Index(line, "Log flushed up to") == 0 {
			val, _ := atof(record[4])
			if len(record) >= 6 {
				val = float64(makeBigint(record[4], record[5]))
			}
			(*p)["log_bytes_flushed"] = val
			continue
		}
		if strings.Index(line, "Last checkpoint at") == 0 {
			val, _ := atof(record[3])
			if len(record) >= 5 {
				val = float64(makeBigint(record[3], record[4]))
			}
			(*p)["last_checkpoint"] = val
			continue
		}

		// Buffer Pool and Memory
		// 5.6 or before
		if strings.Index(line, "Total memory allocated") == 0 && strings.Index(line, "in additional pool allocated") > 0 {
			(*p)["total_mem_alloc"], _ = atof(record[3])
			(*p)["additional_pool_alloc"], _ = atof(record[8])
			continue
		}
		// 5.7
		if strings.Index(line, "Total large memory allocated") == 0 {
			(*p)["total_mem_alloc"], _ = atof(record[4])
			continue
		}

		if strings.Index(line, "Adaptive hash index ") == 0 {
			v, _ := atof(record[3])
			setIfEmpty(p, "adaptive_hash_memory", v)
			continue
		}
		if strings.Index(line, "Page hash           ") == 0 {
			v, _ := atof(record[2])
			setIfEmpty(p, "page_hash_memory", v)
			continue
		}
		if strings.Index(line, "Dictionary cache    ") == 0 {
			v, _ := atof(record[2])
			setIfEmpty(p, "dictionary_cache_memory", v)
			continue
		}
		if strings.Index(line, "File system         ") == 0 {
			v, _ := atof(record[2])
			setIfEmpty(p, "file_system_memory", v)
			continue
		}
		if strings.Index(line, "Lock system         ") == 0 {
			v, _ := atof(record[2])
			setIfEmpty(p, "lock_system_memory", v)
			continue
		}
		if strings.Index(line, "Recovery system     ") == 0 {
			v, _ := atof(record[2])
			setIfEmpty(p, "recovery_system_memory", v)
			continue
		}
		if strings.Index(line, "Threads             ") == 0 {
			v, _ := atof(record[1])
			setIfEmpty(p, "thread_hash_memory", v)
			continue
		}
		if strings.Index(line, "innodb_io_pattern   ") == 0 {
			v, _ := atof(record[1])
			setIfEmpty(p, "innodb_io_pattern_memory", v)
			continue
		}
		if strings.Index(line, "Buffer pool size ") == 0 {
			v, _ := atof(record[3])
			setIfEmpty(p, "pool_size", v)
			continue
		}
		if strings.Index(line, "Free buffers") == 0 {
			v, _ := atof(record[2])
			setIfEmpty(p, "free_pages", v)
			continue
		}
		if strings.Index(line, "Database pages") == 0 {
			v, _ := atof(record[2])
			setIfEmpty(p, "database_pages", v)
			continue
		}
		if strings.Index(line, "Modified db pages") == 0 {
			v, _ := atof(record[3])
			setIfEmpty(p, "modified_pages", v)
			continue
		}
		if strings.Index(line, "Pages read ahead") == 0 {
			v, _ := atof(record[3])
			setIfEmpty(p, "read_ahead", v)
			v, _ = atof(record[7])
			setIfEmpty(p, "read_evicted", v)
			v, _ = atof(record[11])
			setIfEmpty(p, "read_random_ahead", v)
			continue
		}
		if strings.Index(line, "Pages read") == 0 {
			v, _ := atof(record[2])
			setIfEmpty(p, "pages_read", v)
			v, _ = atof(record[4])
			setIfEmpty(p, "pages_created", v)
			v, _ = atof(record[6])
			setIfEmpty(p, "pages_written", v)
			continue
		}

		// Row Operations
		if strings.Index(line, "Number of rows inserted") == 0 {
			(*p)["rows_inserted"], _ = atof(record[4])
			(*p)["rows_updated"], _ = atof(record[6])
			(*p)["rows_deleted"], _ = atof(record[8])
			(*p)["rows_read"], _ = atof(record[10])
			continue
		}
		if strings.Index(line, " queries inside InnoDB, ") == 0 {
			(*p)["queries_inside"], _ = atof(record[0])
			(*p)["queries_queued"], _ = atof(record[4])
			continue
		}

		// for next loop
		prevLine = line
	}

	// finalize
	(*p)["queries_queued"] = (*p)["log_bytes_written"] - (*p)["log_bytes_flushed"]
	(*p)["uncheckpointed_bytes"] = (*p)["log_bytes_written"] - (*p)["last_checkpoint"]
}

func parseProcesslist(state string, p *map[string]float64) {

	if state == "" {
		state = "none"
	} else if state == "Table lock" {
		state = "Locked"
	} else if strings.HasPrefix(state, "Waiting for ") && strings.HasSuffix(state, "lock") {
		state = "Locked"
	}
	state = strings.Replace(strings.ToLower(state), " ", "_", -1)

	state = strings.Join([]string{"State_", state}, "")
	if _, found := processState[state]; !found {
		state = "State_other"
	}

	increaseMap(p, state, "1")
}

func atof(str string) (float64, error) {
	str = strings.Replace(str, ",", "", -1)
	str = strings.Replace(str, ";", "", -1)
	str = strings.Replace(str, "/s", "", -1)
	str = strings.Trim(str, " ")
	return strconv.ParseFloat(str, 64)
}

func increaseMap(p *map[string]float64, key string, src string) {
	val, err := atof(src)
	if err != nil {
		val = 0
	}
	_, exists := (*p)[key]
	if !exists {
		(*p)[key] = val
		return
	}
	(*p)[key] = (*p)[key] + val
}

func makeBigint(hi string, lo string) int64 {
	if lo == "" {
		val, _ := strconv.ParseInt(hi, 16, 64)
		return val
	}

	var hiVal int64
	var loVal int64
	if hi != "" {
		hiVal, _ = strconv.ParseInt(hi, 10, 64)
	}
	if lo != "" {
		loVal, _ = strconv.ParseInt(lo, 10, 64)
	}

	val := hiVal * loVal

	return val
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "3306", "Port")
	optSocket := flag.String("socket", "", "Port")
	optUser := flag.String("username", "root", "Username")
	optPass := flag.String("password", "", "Password")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optInnoDB := flag.Bool("disable_innodb", false, "Disable InnoDB metrics")
	optMetricKeyPrefix := flag.String("metric-key-prefix", "mysql", "metric key prefix")
	optEnableExtended := flag.Bool("enable_extended", false, "Enable Extended metrics")
	flag.Parse()

	var mysql MySQLPlugin

	if *optSocket != "" {
		mysql.Target = *optSocket
		mysql.isUnixSocket = true
	} else {
		mysql.Target = fmt.Sprintf("%s:%s", *optHost, *optPort)
	}
	mysql.Username = *optUser
	mysql.Password = *optPass
	mysql.DisableInnoDB = *optInnoDB
	mysql.prefix = *optMetricKeyPrefix
	mysql.EnableExtended = *optEnableExtended
	helper := mp.NewMackerelPlugin(mysql)
	helper.Tempfile = *optTempfile
	helper.Run()
}
