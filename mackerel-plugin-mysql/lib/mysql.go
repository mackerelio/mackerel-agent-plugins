package mpmysql

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/ziutek/mymysql/mysql"

	// MySQL Driver
	"github.com/ziutek/mymysql/native"
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

func (m *MySQLPlugin) defaultGraphdef() map[string]mp.Graphs {
	labelPrefix := strings.Title(strings.Replace(m.MetricKeyPrefix(), "mysql", "MySQL", -1))

	capacityMetrics := []mp.Metrics{
		{Name: "PercentageOfConnections", Label: "Percentage Of Connections", Diff: false, Stacked: false},
	}
	if m.DisableInnoDB != true {
		capacityMetrics = append(capacityMetrics, mp.Metrics{
			Name: "PercentageOfBufferPool", Label: "Percentage Of Buffer Pool", Diff: false, Stacked: false,
		})
	}

	return map[string]mp.Graphs{
		"cmd": {
			Label: labelPrefix + " Command",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "Com_insert", Label: "Insert", Diff: true, Stacked: true},
				{Name: "Com_insert_select", Label: "Insert Select", Diff: true, Stacked: true},
				{Name: "Com_select", Label: "Select", Diff: true, Stacked: true},
				{Name: "Com_update", Label: "Update", Diff: true, Stacked: true},
				{Name: "Com_update_multi", Label: "Update Multi", Diff: true, Stacked: true},
				{Name: "Com_delete", Label: "Delete", Diff: true, Stacked: true},
				{Name: "Com_delete_multi", Label: "Delete Multi", Diff: true, Stacked: true},
				{Name: "Com_replace", Label: "Replace", Diff: true, Stacked: true},
				{Name: "Com_replace_select", Label: "Replace Select", Diff: true, Stacked: true},
				{Name: "Com_load", Label: "Load", Diff: true, Stacked: true},
				{Name: "Com_set_option", Label: "Set Option", Diff: true, Stacked: true},
				// Duplicate of query_cache.Qcache_hits but remains for compatibility reason
				{Name: "Qcache_hits", Label: "Query Cache Hits", Diff: true, Stacked: false},
				{Name: "Questions", Label: "Questions", Diff: true, Stacked: false},
			},
		},
		"join": {
			Label: labelPrefix + " Join/Scan",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "Select_full_join", Label: "Select Full JOIN", Diff: true, Stacked: false},
				{Name: "Select_full_range_join", Label: "Select Full Range JOIN", Diff: true, Stacked: false},
				{Name: "Select_range", Label: "Select Range", Diff: true, Stacked: false},
				{Name: "Select_range_check", Label: "Select Range Check", Diff: true, Stacked: false},
				{Name: "Select_scan", Label: "Select SCAN", Diff: true, Stacked: false},
			},
		},
		"threads": {
			Label: labelPrefix + " Threads",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "thread_cache_size", Label: "Cache Size", Diff: false, Stacked: false},
				// Duplicate of connections.Threads_connected but remains for compatibility reason
				{Name: "Threads_connected", Label: "Connected", Diff: false, Stacked: false},
				{Name: "Threads_running", Label: "Running", Diff: false, Stacked: false},
				{Name: "Threads_created", Label: "Created", Diff: true, Stacked: false},
				{Name: "Threads_cached", Label: "Cached", Diff: false, Stacked: false},
			},
		},
		"connections": {
			Label: labelPrefix + " Connections",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "max_connections", Label: "Max Connections", Diff: false, Stacked: false},
				{Name: "Max_used_connections", Label: "Max Used Connections", Diff: false, Stacked: false},
				{Name: "Connections", Label: "Connections", Diff: true, Stacked: false},
				{Name: "Threads_connected", Label: "Threads Connected", Diff: false, Stacked: false},
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
				{Name: "Bytes_sent", Label: "Sent Bytes", Diff: true, Stacked: false, Scale: (1.0 / 60.0)},
				{Name: "Bytes_received", Label: "Received Bytes", Diff: true, Stacked: false, Scale: (1.0 / 60.0)},
			},
		},
		"capacity": {
			Label:   labelPrefix + " Capacity",
			Unit:    "percentage",
			Metrics: capacityMetrics,
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
	isAuroraReader bool
	Debug          bool
}

// MetricKeyPrefix returns the metrics key prefix
func (m *MySQLPlugin) MetricKeyPrefix() string {
	if m.prefix == "" {
		m.prefix = "mysql"
	}
	return m.prefix
}

func (m *MySQLPlugin) fetchVersion(db mysql.Conn) (version [3]int, err error) {
	rows, _, err := db.Query("SHOW VARIABLES WHERE VARIABLE_NAME = 'VERSION'")
	if err != nil {
		return
	}
	for _, row := range rows {
		if len(row) > 1 {
			versionString := string(row[1].([]byte))
			if i := strings.IndexRune(versionString, '-'); i >= 0 {
				// Trim -log or -debug, -MariaDB-...
				versionString = versionString[:i]
			}
			xs := strings.Split(versionString, ".")
			if len(xs) >= 2 {
				version[0], _ = strconv.Atoi(xs[0])
				version[1], _ = strconv.Atoi(xs[1])
				if len(xs) >= 3 {
					version[2], _ = strconv.Atoi(xs[2])
				}
			}
			break
		}
	}
	if version[0] == 0 {
		err = errors.New("failed to get mysql version")
		return
	}
	return
}

func (m *MySQLPlugin) fetchShowStatus(db mysql.Conn, stat map[string]float64) error {
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
			v, err := atof(string(row[1].([]byte)))
			if err != nil {
				continue
			}
			stat[variableName] = v
		} else {
			log.Fatalln("FetchMetrics (Status): row length is too small: ", len(row))
		}
	}
	return nil
}

func (m *MySQLPlugin) fetchShowInnodbStatus(db mysql.Conn, stat map[string]float64) error {
	row, _, err := db.QueryFirst("SHOW /*!50000 ENGINE*/ INNODB STATUS")
	if err != nil {
		log.Println("FetchMetrics (InnoDB Status): ", err)
		log.Fatalln("Hint: If you don't use InnoDB and see InnoDB Status error, you should set -disable_innodb")
	}

	var trxIDHexFormat bool
	v, err := m.fetchVersion(db)
	if err != nil {
		log.Println(err)
	}
	// Transaction IDs are printed in hex format in version < 5.6.4.
	//   Ref: https://github.com/mysql/mysql-server/commit/3420dc52b68c9afcee0a19ba7c19a73c2fbb2913
	//        https://github.com/mysql/mysql-server/blob/mysql-5.6.3/storage/innobase/include/trx0types.h#L32
	//        https://github.com/mysql/mysql-server/blob/mysql-5.6.4/storage/innobase/include/trx0types.h#L32
	// MariaDB 10.x is recognized as newer than 5.6.4, which should be correct.
	//   Ref: https://github.com/MariaDB/server/blob/mariadb-10.0.0/storage/innobase/include/trx0types.h#L32
	if v[0] < 5 || v[0] == 5 && v[1] < 6 || v[0] == 5 && v[1] == 6 && v[2] < 4 {
		trxIDHexFormat = true
	}

	if len(row) > 0 {
		parseInnodbStatus(string(row[len(row)-1].([]byte)), trxIDHexFormat, stat)
	} else {
		return fmt.Errorf("row length is too small: %d", len(row))
	}
	return nil
}

func (m *MySQLPlugin) fetchShowVariables(db mysql.Conn, stat map[string]float64) error {
	rows, _, err := db.Query("SHOW VARIABLES")
	if err != nil {
		log.Fatalln("FetchMetrics (Variables): ", err)
	}

	rawStat := make(map[string]string)
	for _, row := range rows {
		if len(row) > 1 {
			variableName := string(row[0].([]byte))
			if err != nil {
				log.Println("FetchMetrics (Fetch Variables): ", err)
			}
			value := string(row[1].([]byte))
			rawStat[variableName] = value
			v, err := atof(value)
			if err != nil {
				continue
			}
			stat[variableName] = v
		} else {
			log.Fatalln("FetchMetrics (Variables): row length is too small: ", len(row))
		}
	}

	m.isAuroraReader = rawStat["aurora_version"] != "" && rawStat["innodb_read_only"] == "ON"

	if m.EnableExtended {
		err = fetchShowVariablesBackwardCompatibile(stat)
		if err != nil {
			log.Fatalln("FetchExtendedMetrics (Fetch Variables): ", err)
		}
		if _, found := stat["key_cache_block_size"]; found {
			if _, found = stat["Key_blocks_unused"]; found {
				stat["key_buf_bytes_used"] = stat["key_buffer_size"] - stat["Key_blocks_unused"]*stat["key_cache_block_size"]
			}
			if _, found = stat["Key_blocks_not_flushed"]; found {
				stat["key_buf_bytes_unflushed"] = stat["Key_blocks_not_flushed"] * stat["key_cache_block_size"]
			}
		}
	}
	return nil
}

func fetchShowVariablesBackwardCompatibile(stat map[string]float64) error {
	// https://github.com/percona/percona-monitoring-plugins/blob/v1.1.6/cacti/scripts/ss_get_mysql_stats.php#L585
	if val, found := stat["table_open_cache"]; found {
		stat["table_cache"] = val
	}
	return nil
}

func (m *MySQLPlugin) fetchShowSlaveStatus(db mysql.Conn, stat map[string]float64) error {
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

func (m *MySQLPlugin) fetchProcesslist(db mysql.Conn, stat map[string]float64) error {
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
			parseProcesslist(state, stat)
		} else {
			log.Fatalln("FetchMetrics (Processlist): row length is too small: ", len(row))
		}
	}

	return nil
}

// for backward compatibility, convert stat names
func (m *MySQLPlugin) convertInnodbStats(stat map[string]float64) {
	for dst, src := range map[string]string{
		"database_pages":    "Innodb_buffer_pool_pages_data",
		"free_pages":        "Innodb_buffer_pool_pages_free",
		"file_fsyncs":       "Innodb_data_fsyncs",
		"file_reads":        "Innodb_data_reads",
		"file_writes":       "Innodb_data_writes",
		"modified_pages":    "Innodb_buffer_pool_pages_dirty",
		"pool_size":         "Innodb_buffer_pool_pages_total",
		"pages_created":     "Innodb_pages_created",
		"pages_read":        "Innodb_pages_read",
		"pages_written":     "Innodb_pages_written",
		"read_ahead":        "Innodb_buffer_pool_read_ahead",
		"read_evicted":      "Innodb_buffer_pool_read_ahead_evicted",
		"read_random_ahead": "Innodb_buffer_pool_read_ahead_rnd",
	} {
		if v, found := stat[src]; found {
			setIfEmpty(stat, dst, v)
		}
	}
}

func (m *MySQLPlugin) calculateCapacity(stat map[string]float64) {
	stat["PercentageOfConnections"] = 100.0 * stat["Threads_connected"] / stat["max_connections"]
	if m.DisableInnoDB != true {
		stat["PercentageOfBufferPool"] = 100.0 * stat["database_pages"] / stat["pool_size"]
	}
}

// FetchMetrics interface for mackerelplugin
func (m *MySQLPlugin) FetchMetrics() (map[string]float64, error) {
	proto := "tcp"
	if m.isUnixSocket {
		proto = "unix"
	}
	db := mysql.New(proto, "", m.Target, m.Username, m.Password, "")
	switch c := db.(type) {
	case *native.Conn:
		c.Debug = m.Debug
	}
	err := db.Connect()
	if err != nil {
		log.Fatalln("FetchMetrics (DB Connect): ", err)
		return nil, err
	}
	defer db.Close()

	stat := make(map[string]float64)
	m.fetchShowStatus(db, stat)

	m.fetchShowVariables(db, stat)

	if m.DisableInnoDB != true {
		m.convertInnodbStats(stat)
		if !m.isAuroraReader {
			err := m.fetchShowInnodbStatus(db, stat)
			if err != nil {
				log.Println("FetchMetrics (InnoDB Status): ", err)
				m.DisableInnoDB = true
			}
		}
	}

	m.fetchShowSlaveStatus(db, stat)

	if m.EnableExtended {
		m.fetchProcesslist(db, stat)
	}

	m.calculateCapacity(stat)

	explicitMetricNames := m.metricNames()
	statRet := make(map[string]float64)
	for key, value := range stat {
		if _, ok := explicitMetricNames[key]; !ok {
			continue
		}
		statRet[key] = value
	}

	return statRet, err
}

func (m *MySQLPlugin) metricNames() map[string]struct{} {
	a := make(map[string]struct{})
	for _, g := range m.GraphDefinition() {
		for _, metric := range g.Metrics {
			a[metric.Name] = struct{}{}
		}
	}
	return a
}

// GraphDefinition interface for mackerelplugin
func (m *MySQLPlugin) GraphDefinition() map[string]mp.Graphs {
	graphdef := m.defaultGraphdef()
	if !m.DisableInnoDB {
		graphdef = m.addGraphdefWithInnoDBMetrics(graphdef)
	}
	if m.EnableExtended {
		graphdef = m.addExtendedGraphdef(graphdef)
	}
	return graphdef
}

func (m *MySQLPlugin) addGraphdefWithInnoDBMetrics(graphdef map[string]mp.Graphs) map[string]mp.Graphs {
	labelPrefix := strings.Title(strings.Replace(m.MetricKeyPrefix(), "mysql", "MySQL", -1))
	graphdef["innodb_rows"] = mp.Graphs{
		Label: labelPrefix + " innodb Rows",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "Innodb_rows_read", Label: "Read", Diff: true, Stacked: false},
			{Name: "Innodb_rows_inserted", Label: "Inserted", Diff: true, Stacked: false},
			{Name: "Innodb_rows_updated", Label: "Updated", Diff: true, Stacked: false},
			{Name: "Innodb_rows_deleted", Label: "Deleted", Diff: true, Stacked: false},
		},
	}
	graphdef["innodb_row_lock_time"] = mp.Graphs{
		Label: labelPrefix + " innodb Row Lock Time",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "Innodb_row_lock_time", Label: "Lock Time", Diff: true, Stacked: false},
		},
	}
	graphdef["innodb_row_lock_waits"] = mp.Graphs{
		Label: labelPrefix + " innodb Row Lock Waits",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "Innodb_row_lock_waits", Label: "Lock Waits", Diff: true, Stacked: false},
		},
	}
	graphdef["innodb_adaptive_hash_index"] = mp.Graphs{
		Label: labelPrefix + " innodb Adaptive Hash Index",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "hash_index_cells_total", Label: "Hash Index Cells Total", Diff: false, Stacked: false},
			{Name: "hash_index_cells_used", Label: "Hash Index Cells Used", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_buffer_pool_read"] = mp.Graphs{
		Label: labelPrefix + " innodb Buffer Pool Read (/sec)",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "read_ahead", Label: "Pages Read Ahead", Diff: true, Stacked: false, Scale: (1.0 / 60.0)},
			{Name: "read_evicted", Label: "Evicted Without Access", Diff: true, Stacked: false, Scale: (1.0 / 60.0)},
			{Name: "read_random_ahead", Label: "Random Read Ahead", Diff: true, Stacked: false, Scale: (1.0 / 60.0)},
		},
	}
	graphdef["innodb_buffer_pool_activity"] = mp.Graphs{
		Label: labelPrefix + " innodb Buffer Pool Activity (Pages)",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "pages_created", Label: "Created", Diff: true, Stacked: false},
			{Name: "pages_read", Label: "Read", Diff: true, Stacked: false},
			{Name: "pages_written", Label: "Written", Diff: true, Stacked: false},
		},
	}
	graphdef["innodb_buffer_pool_efficiency"] = mp.Graphs{
		Label: labelPrefix + " innodb Buffer Pool Efficiency",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "Innodb_buffer_pool_reads", Label: "Reads", Diff: true, Stacked: false},
			{Name: "Innodb_buffer_pool_read_requests", Label: "Read Requests", Diff: true, Stacked: false},
		},
	}
	graphdef["innodb_buffer_pool"] = mp.Graphs{
		Label: labelPrefix + " innodb Buffer Pool (Pages)",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "pool_size", Label: "Pool Size", Diff: false, Stacked: false},
			{Name: "database_pages", Label: "Used", Diff: false, Stacked: true},
			{Name: "free_pages", Label: "Free", Diff: false, Stacked: true},
			{Name: "modified_pages", Label: "Modified", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_checkpoint_age"] = mp.Graphs{
		Label: labelPrefix + " innodb Checkpoint Age",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "uncheckpointed_bytes", Label: "Uncheckpointed", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_current_lock_waits"] = mp.Graphs{
		Label: labelPrefix + " innodb Current Lock Waits (secs)",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "innodb_lock_wait_secs", Label: "Innodb Lock Wait", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_io"] = mp.Graphs{
		Label: labelPrefix + " innodb I/O",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "file_reads", Label: "File Reads", Diff: true, Stacked: false},
			{Name: "file_writes", Label: "File Writes", Diff: true, Stacked: false},
			{Name: "file_fsyncs", Label: "File fsyncs", Diff: true, Stacked: false},
			{Name: "log_writes", Label: "Log Writes", Diff: true, Stacked: false},
		},
	}
	graphdef["innodb_io_pending"] = mp.Graphs{
		Label: labelPrefix + " innodb I/O Pending",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "pending_normal_aio_reads", Label: "Normal AIO Reads", Diff: false, Stacked: false},
			{Name: "pending_normal_aio_writes", Label: "Normal AIO Writes", Diff: false, Stacked: false},
			{Name: "pending_ibuf_aio_reads", Label: "InnoDB Buffer AIO Reads", Diff: false, Stacked: false},
			{Name: "pending_aio_log_ios", Label: "AIO Log IOs", Diff: false, Stacked: false},
			{Name: "pending_aio_sync_ios", Label: "AIO Sync IOs", Diff: false, Stacked: false},
			{Name: "pending_log_flushes", Label: "Log Flushes (fsync)", Diff: false, Stacked: false},
			{Name: "pending_buf_pool_flushes", Label: "Buffer Pool Flushes", Diff: false, Stacked: false},
			{Name: "pending_log_writes", Label: "Log Writes", Diff: false, Stacked: false},
			{Name: "pending_chkp_writes", Label: "Checkpoint Writes", Diff: false, Stacked: false},
			{Name: "log_pending_log_flushes", Label: "Log Flushes (log)", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_insert_buffer"] = mp.Graphs{
		Label: labelPrefix + " innodb Insert Buffer",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "ibuf_inserts", Label: "Inserts", Diff: true, Stacked: false},
			{Name: "ibuf_merges", Label: "Merges", Diff: true, Stacked: false},
			{Name: "ibuf_merged", Label: "Merged", Diff: true, Stacked: false},
		},
	}
	graphdef["innodb_insert_buffer_usage"] = mp.Graphs{
		Label: labelPrefix + " innodb Insert Buffer Usage (Cells)",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "ibuf_cell_count", Label: "Cell Count", Diff: false, Stacked: false},
			{Name: "ibuf_used_cells", Label: "Used", Diff: false, Stacked: true},
			{Name: "ibuf_free_cells", Label: "Free", Diff: false, Stacked: true},
		},
	}
	graphdef["innodb_lock_structures"] = mp.Graphs{
		Label: labelPrefix + " innodb Lock Structures",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "innodb_lock_structs", Label: "Structures", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_log"] = mp.Graphs{
		Label: labelPrefix + " innodb Log",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "log_bytes_written", Label: "Written", Diff: true, Stacked: false},
			{Name: "log_bytes_flushed", Label: "Flushed", Diff: true, Stacked: false},
			{Name: "unflushed_log", Label: "Unflushed", Diff: false, Stacked: false},
			{Name: "innodb_log_buffer_size", Label: "Buffer Size", Diff: false, Stacked: true},
		},
	}
	graphdef["innodb_memory_allocation"] = mp.Graphs{
		Label: labelPrefix + " innodb Memory Allocation",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "additional_pool_alloc", Label: "Additional Pool Allocated", Diff: false, Stacked: false},
			{Name: "total_mem_alloc", Label: "Total Memory Allocated", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_semaphores"] = mp.Graphs{
		Label: labelPrefix + " innodb Semaphores",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "spin_waits", Label: "Spin Waits", Diff: true, Stacked: false},
			{Name: "spin_rounds", Label: "Spin Rounds", Diff: true, Stacked: false},
			{Name: "os_waits", Label: "OS Waits", Diff: true, Stacked: false},
		},
	}
	graphdef["innodb_tables_in_use"] = mp.Graphs{
		Label: labelPrefix + " innodb Tables In Use",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "innodb_locked_tables", Label: "Locked Tables", Diff: false, Stacked: false},
			{Name: "innodb_tables_in_use", Label: "Table in Use", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_transactions_active_locked"] = mp.Graphs{
		Label: labelPrefix + " innodb Transactions Active/Locked",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "current_transactions", Label: "Current", Diff: false, Stacked: false},
			{Name: "active_transactions", Label: "Active", Diff: false, Stacked: false},
			{Name: "locked_transactions", Label: "Locked", Diff: false, Stacked: false},
			{Name: "read_views", Label: "Read Views", Diff: false, Stacked: false},
		},
	}
	graphdef["innodb_transactions"] = mp.Graphs{
		Label: labelPrefix + " innodb Transactions",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "history_list", Label: "History List", Diff: false, Stacked: false},
			{Name: "innodb_transactions", Label: "InnoDB Transactions", Diff: true, Stacked: false},
		},
	}
	return graphdef
}

func (m *MySQLPlugin) addExtendedGraphdef(graphdef map[string]mp.Graphs) map[string]mp.Graphs {
	//TODO
	labelPrefix := strings.Title(strings.Replace(m.MetricKeyPrefix(), "mysql", "MySQL", -1))
	graphdef["query_cache"] = mp.Graphs{
		Label: labelPrefix + " query Cache",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "Qcache_queries_in_cache", Label: "Qcache Queries In Cache", Diff: false, Stacked: false},
			{Name: "Qcache_hits", Label: "Qcache Hits", Diff: true, Stacked: false},
			{Name: "Qcache_inserts", Label: "Qcache Inserts", Diff: true, Stacked: false},
			{Name: "Qcache_not_cached", Label: "Qcache Not Cached", Diff: true, Stacked: false},
			{Name: "Qcache_lowmem_prunes", Label: "Qcache Lowmem Prunes", Diff: true, Stacked: false},
		},
	}
	graphdef["query_cache_memory"] = mp.Graphs{
		Label: labelPrefix + " query Cache Memory",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "query_cache_size", Label: "Query Cache Size", Diff: false, Stacked: false},
			{Name: "Qcache_free_memory", Label: "Qcache Free Memory", Diff: false, Stacked: false},
			{Name: "Qcache_total_blocks", Label: "Qcache Total Blocks", Diff: false, Stacked: false},
			{Name: "Qcache_free_blocks", Label: "Qcache Free Blocks", Diff: false, Stacked: false},
		},
	}
	graphdef["temporary_objects"] = mp.Graphs{
		Label: labelPrefix + " temporary Objects",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "Created_tmp_tables", Label: "Created Tmp Tables", Diff: true, Stacked: false},
			{Name: "Created_tmp_disk_tables", Label: "Created Tmp Disk Tables", Diff: true, Stacked: false},
			{Name: "Created_tmp_files", Label: "Created Tmp Files", Diff: true, Stacked: false},
		},
	}
	graphdef["files_and_tables"] = mp.Graphs{
		Label: labelPrefix + " files and Tables",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "table_cache", Label: "Table Cache", Diff: false, Stacked: false},
			{Name: "Open_tables", Label: "Open Tables", Diff: false, Stacked: false},
			{Name: "Open_files", Label: "Open Files", Diff: false, Stacked: false},
			{Name: "Opened_tables", Label: "Opened Tables", Diff: true, Stacked: false},
		},
	}
	graphdef["processlist"] = mp.Graphs{
		Label: labelPrefix + " processlist",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "State_closing_tables", Label: "State Closing Tables", Diff: false, Stacked: true},
			{Name: "State_copying_to_tmp_table", Label: "State Copying To Tmp Table", Diff: false, Stacked: true},
			{Name: "State_end", Label: "State End", Diff: false, Stacked: true},
			{Name: "State_freeing_items", Label: "State Freeing Items", Diff: false, Stacked: true},
			{Name: "State_init", Label: "State Init", Diff: false, Stacked: true},
			{Name: "State_locked", Label: "State Locked", Diff: false, Stacked: true},
			{Name: "State_login", Label: "State Login", Diff: false, Stacked: true},
			{Name: "State_preparing", Label: "State Preparing", Diff: false, Stacked: true},
			{Name: "State_reading_from_net", Label: "State Reading From Net", Diff: false, Stacked: true},
			{Name: "State_sending_data", Label: "State Sending Data", Diff: false, Stacked: true},
			{Name: "State_sorting_result", Label: "State Sorting Result", Diff: false, Stacked: true},
			{Name: "State_statistics", Label: "State Statistics", Diff: false, Stacked: true},
			{Name: "State_updating", Label: "State Updating", Diff: false, Stacked: true},
			{Name: "State_writing_to_net", Label: "State Writing To Net", Diff: false, Stacked: true},
			{Name: "State_none", Label: "State None", Diff: false, Stacked: true},
			{Name: "State_other", Label: "State Other", Diff: false, Stacked: true},
		},
	}
	graphdef["sorts"] = mp.Graphs{
		Label: labelPrefix + " sorts",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "Sort_rows", Label: "Sort Rows", Diff: true, Stacked: false},
			{Name: "Sort_range", Label: "Sort Range", Diff: true, Stacked: false},
			{Name: "Sort_merge_passes", Label: "Sort Merge Passes", Diff: true, Stacked: false},
			{Name: "Sort_scan", Label: "Sort Scan", Diff: true, Stacked: false},
		},
	}
	graphdef["handlers"] = mp.Graphs{
		Label: labelPrefix + " handlers",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "Handler_write", Label: "Handler Write", Diff: true, Stacked: true},
			{Name: "Handler_update", Label: "Handler Update", Diff: true, Stacked: true},
			{Name: "Handler_delete", Label: "Handler Delete", Diff: true, Stacked: true},
			{Name: "Handler_read_first", Label: "Handler Read First", Diff: true, Stacked: true},
			{Name: "Handler_read_key", Label: "Handler Read Key", Diff: true, Stacked: true},
			{Name: "Handler_read_last", Label: "Handler Read Last", Diff: true, Stacked: true},
			{Name: "Handler_read_next", Label: "Handler Read Next", Diff: true, Stacked: true},
			{Name: "Handler_read_prev", Label: "Handler Read Prev", Diff: true, Stacked: true},
			{Name: "Handler_read_rnd", Label: "Handler Read Rnd", Diff: true, Stacked: true},
			{Name: "Handler_read_rnd_next", Label: "Handler Read Rnd Next", Diff: true, Stacked: true},
		},
	}
	graphdef["transaction_handler"] = mp.Graphs{
		Label: labelPrefix + " transaction handler",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "Handler_commit", Label: "Handler Commit", Diff: true, Stacked: false},
			{Name: "Handler_rollback", Label: "Handler Rollback", Diff: true, Stacked: false},
			{Name: "Handler_savepoint", Label: "Handler Savepoint", Diff: true, Stacked: false},
			{Name: "Handler_savepoint_rollback", Label: "Handler Savepoint Rollback", Diff: true, Stacked: false},
		},
	}
	graphdef["myisam_indexes"] = mp.Graphs{
		Label: labelPrefix + " MyISAM Indexes",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "Key_read_requests", Label: "Key Read Requests", Diff: true, Stacked: false},
			{Name: "Key_reads", Label: "Key Reads", Diff: true, Stacked: false},
			{Name: "Key_write_requests", Label: "Key Write Requests", Diff: true, Stacked: false},
			{Name: "Key_writes", Label: "Key Writes", Diff: true, Stacked: false},
		},
	}
	graphdef["myisam_key_cache"] = mp.Graphs{
		Label: labelPrefix + " MyISAM Key Cache",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "key_buffer_size", Label: "Key Buffer Size", Diff: false, Stacked: false},
			{Name: "key_buf_bytes_used", Label: "Key Buf Bytes Used", Diff: false, Stacked: false},
			{Name: "key_buf_bytes_unflushed", Label: "Key Buf Bytes Unflushed", Diff: false, Stacked: false},
		},
	}

	return graphdef
}

func setIfEmpty(p map[string]float64, key string, val float64) {
	_, ok := p[key]
	if !ok {
		p[key] = val
	}
}

func calculateAio(s string) (int, error) {
	v := strings.TrimSpace(s)
	i := strings.Index(v, "[")
	j := strings.Index(v, "]")

	if i >= 0 && j > i {
		counts := strings.Split(v[i+1:j], ",")
		total := 0
		for _, c := range counts {
			p, err := strconv.Atoi(strings.TrimSpace(c))
			if err != nil {
				return 0, err
			}
			total += p
		}
		return total, nil
	}
	v = strings.TrimSpace(strings.TrimRight(v, ","))
	if v == "" {
		return 0, nil
	}
	return strconv.Atoi(v)
}

func parseInnodbStatus(str string, trxIDHexFormat bool, p map[string]float64) {
	isTransaction := false
	prevLine := ""

	for _, line := range strings.Split(str, "\n") {
		line = strings.TrimLeft(line, " ")
		record := strings.Fields(line)

		// Innodb Semaphores
		if strings.HasPrefix(line, "Mutex spin waits") {
			increaseMap(p, "spin_waits", record[3])
			increaseMap(p, "spin_rounds", record[5])
			increaseMap(p, "os_waits", record[8])
			continue
		}
		if strings.HasPrefix(line, "RW-shared spins") && strings.Contains(line, ";") {
			// 5.5, 5.6
			increaseMap(p, "spin_waits", record[2])
			increaseMap(p, "os_waits", record[5])
			increaseMap(p, "spin_waits", record[8])
			increaseMap(p, "os_waits", record[11])
			continue
		}
		if strings.HasPrefix(line, "RW-shared spins") && !strings.Contains(line, "; RW-excl spins") {
			// 5.1
			increaseMap(p, "spin_waits", record[2])
			increaseMap(p, "os_waits", record[7])
			continue
		}
		if strings.HasPrefix(line, "RW-excl spins") {
			// 5.5, 5.6
			increaseMap(p, "spin_waits", record[2])
			increaseMap(p, "os_waits", record[7])
			continue
		}
		if strings.HasPrefix(line, "RW-sx spins") {
			// 5.7
			increaseMap(p, "spin_waits", record[2])
			increaseMap(p, "os_waits", record[7])
			continue
		}
		if strings.Contains(line, "seconds the semaphore:") {
			increaseMap(p, "innodb_sem_waits", "1")
			wait, err := atof(record[9])
			if err != nil {
				continue
			}
			wait = wait * 1000
			increaseMap(p, "innodb_sem_wait_time_ms", fmt.Sprintf("%.f", wait))
			continue
		}

		// Innodb Transactions
		if strings.HasPrefix(line, "Trx id counter") {
			loVal := ""
			if len(record) >= 5 {
				loVal = record[4]
			}
			val := makeBigint(record[3], loVal, trxIDHexFormat)
			increaseMap(p, "innodb_transactions", fmt.Sprintf("%d", val))
			isTransaction = true
			continue
		}
		if strings.HasPrefix(line, "Purge done for trx") {
			if record[7] == "undo" {
				record[7] = ""
			}
			val := makeBigint(record[6], record[7], trxIDHexFormat)
			trx := p["innodb_transactions"] - float64(val)
			increaseMap(p, "unpurged_txns", fmt.Sprintf("%.f", trx))
			continue
		}
		if strings.HasPrefix(line, "History list length") {
			increaseMap(p, "history_list", record[3])
			continue
		}
		if isTransaction && strings.HasPrefix(line, "---TRANSACTION") {
			increaseMap(p, "current_transactions", "1")
			if strings.Contains(line, "ACTIVE") {
				increaseMap(p, "active_transactions", "1")
			}
			continue
		}
		if isTransaction && strings.HasPrefix(line, "------- TRX HAS BEEN") {
			increaseMap(p, "innodb_lock_wait_secs", record[5])
			continue
		}
		if strings.Contains(line, "read views open inside InnoDB") {
			setMap(p, "read_views", record[0])
			continue
		}
		if isTransaction && strings.HasPrefix(line, "mysql tables in use") {
			increaseMap(p, "innodb_tables_in_use", record[4])
			increaseMap(p, "innodb_locked_tables", record[6])
			continue
		}
		if isTransaction && strings.Index(line, "lock struct(s)") > 0 {
			if strings.HasPrefix(line, "LOCK WAIT") {
				increaseMap(p, "innodb_lock_structs", record[2])
				increaseMap(p, "locked_transactions", "1")
			} else {
				increaseMap(p, "innodb_lock_structs", record[0])
			}
			continue
		}

		// File I/O
		if strings.HasPrefix(line, "Pending normal aio reads:") {
			// There are many 'Pending normal aio' lines.
			//
			//   Pending normal aio reads: 0, aio writes: 0,
			//   Pending normal aio reads: 0 [0, 0, 0, 0] , aio writes: 0 [0, 0, 0, 0] ,
			//   Pending normal aio reads: [0, 0, 0, 0] , aio writes: [0, 0, 0, 0] ,
			//   ...etc

			reads := 0
			writes := 0
			s := strings.TrimPrefix(line, "Pending normal aio reads:")
			writeSep := ", aio writes:"
			var err error
			if strings.Contains(s, writeSep) {
				values := strings.Split(s, writeSep)
				reads, err = calculateAio(values[0])
				if err != nil {
					log.Println(err)
				}
				writes, err = calculateAio(values[1])
				if err != nil {
					log.Println(err)
				}
			} else {
				reads, err = calculateAio(s)
				if err != nil {
					log.Println(err)
				}
			}
			p["pending_normal_aio_reads"] = float64(reads)
			p["pending_normal_aio_writes"] = float64(writes)
			continue
		}
		if strings.HasPrefix(line, "ibuf aio reads") && len(record) >= 10 {
			setMap(p, "pending_ibuf_aio_reads", record[3])
			setMap(p, "pending_aio_log_ios", record[6])
			setMap(p, "pending_aio_sync_ios", record[9])
			continue
		}
		if strings.Index(line, "Pending flushes (fsync)") == 0 {
			setMap(p, "pending_log_flushes", record[4])
			setMap(p, "pending_buf_pool_flushes", record[7])
			continue
		}

		// Insert Buffer and Adaptive Hash Index
		if strings.HasPrefix(line, "Ibuf for space 0: size ") {
			setMap(p, "ibuf_used_cells", record[5])
			setMap(p, "ibuf_free_cells", record[9])
			setMap(p, "ibuf_cell_count", record[12])
			continue
		}
		if strings.HasPrefix(line, "Ibuf: size ") {
			setMap(p, "ibuf_used_cells", record[2])
			setMap(p, "ibuf_free_cells", record[6])
			setMap(p, "ibuf_cell_count", record[9])
			if strings.Contains(line, "merges") {
				setMap(p, "ibuf_merges", record[10])
			}
			continue
		}
		if strings.Contains(line, ", delete mark ") && strings.HasPrefix(prevLine, "merged operations:") {
			setMap(p, "ibuf_inserts", record[1])
			v1, e1 := atof(record[1])
			v2, e2 := atof(record[4])
			v3, e3 := atof(record[6])
			if e1 == nil && e2 == nil && e3 == nil {
				p["ibuf_merged"] = v1 + v2 + v3
			}
			continue
		}
		if strings.Contains(line, " merged recs, ") {
			setMap(p, "ibuf_inserts", record[0])
			setMap(p, "ibuf_merged", record[2])
			setMap(p, "ibuf_merges", record[5])
			continue
		}
		if strings.HasPrefix(line, "Hash table size ") {
			setMap(p, "hash_index_cells_total", record[3])
			if strings.Contains(line, "used cells") {
				setMap(p, "hash_index_cells_used", record[6])
			} else {
				p["hash_index_cells_used"] = 0
			}
			continue
		}

		// Log
		if strings.Contains(line, " log i/o's done, ") {
			setMap(p, "log_writes", record[0])
			continue
		}
		// MySQL < 5.7
		if strings.Contains(line, " pending log writes, ") {
			//      Log sequence number 3091027710
			//      Log flushed up to   3090240098
			//      Pages flushed up to 3074432960
			//      Last checkpoint at  3050856266
			//      0 pending log writes, 0 pending chkp writes
			//      1187 log i/o's done, 14.67 log i/o's/second
			setMap(p, "pending_log_writes", record[0])
			setMap(p, "pending_chkp_writes", record[4])
			continue
		}
		// MySQL >= 5.7 < 8
		if strings.Contains(line, " pending log flushes, ") {
			//  Log sequence number 12665751
			//  Log flushed up to   12665751
			//  Pages flushed up to 12665751
			//  Last checkpoint at  12665742
			//  0 pending log flushes, 0 pending chkp writes
			// 10 log i/o's done, 0.00 log i/o's/second
			setMap(p, "log_pending_log_flushes", record[0])
			setMap(p, "pending_chkp_writes", record[4])
			continue
		}
		// TODO: MySQL 8 does not output pending log writes / flushes, so we perhaps need another way to obtain these metrics.
		//	Log sequence number          28622392
		//	Log buffer assigned up to    28622392
		//	Log buffer completed up to   28622392
		//	Log written up to            28622392
		//	Log flushed up to            28622392
		//	Added dirty pages up to      28622392
		//	Pages flushed up to          28622392
		//	Last checkpoint at           28622392
		//	25 log i/o's done, 0.00 log i/o's/second

		if strings.HasPrefix(line, "Log sequence number") {
			val, err := atof(record[3])
			if err != nil {
				continue
			}
			if len(record) >= 5 {
				val = float64(makeBigint(record[3], record[4], false))
			}
			p["log_bytes_written"] = val
			continue
		}
		if strings.HasPrefix(line, "Log flushed up to") {
			val, err := atof(record[4])
			if err != nil {
				continue
			}
			if len(record) >= 6 {
				val = float64(makeBigint(record[4], record[5], false))
			}
			p["log_bytes_flushed"] = val
			continue
		}
		if strings.HasPrefix(line, "Last checkpoint at") {
			val, err := atof(record[3])
			if err != nil {
				continue
			}
			if len(record) >= 5 {
				val = float64(makeBigint(record[3], record[4], false))
			}
			p["last_checkpoint"] = val
			continue
		}

		// Buffer Pool and Memory
		// 5.6 or before
		if strings.HasPrefix(line, "Total memory allocated") && strings.Contains(line, "in additional pool allocated") {
			setMap(p, "total_mem_alloc", record[3])
			setMap(p, "additional_pool_alloc", record[8])
			continue
		}
		// 5.7
		if strings.HasPrefix(line, "Total large memory allocated") {
			setMap(p, "total_mem_alloc", record[4])
			continue
		}

		if strings.HasPrefix(line, "Adaptive hash index ") {
			setMapIfEmpty(p, "adaptive_hash_memory", record[3])
			continue
		}
		if strings.HasPrefix(line, "Page hash           ") {
			setMapIfEmpty(p, "page_hash_memory", record[2])
			continue
		}
		if strings.HasPrefix(line, "Dictionary cache    ") {
			setMapIfEmpty(p, "dictionary_cache_memory", record[2])
			continue
		}
		if strings.HasPrefix(line, "File system         ") {
			setMapIfEmpty(p, "file_system_memory", record[2])
			continue
		}
		if strings.HasPrefix(line, "Lock system         ") {
			setMapIfEmpty(p, "lock_system_memory", record[2])
			continue
		}
		if strings.HasPrefix(line, "Recovery system     ") {
			setMapIfEmpty(p, "recovery_system_memory", record[2])
			continue
		}
		if strings.HasPrefix(line, "Threads             ") {
			setMapIfEmpty(p, "thread_hash_memory", record[1])
			continue
		}
		if strings.HasPrefix(line, "innodb_io_pattern   ") {
			setMapIfEmpty(p, "innodb_io_pattern_memory", record[1])
			continue
		}

		// for next loop
		prevLine = line
	}

	// finalize
	p["unflushed_log"] = p["log_bytes_written"] - p["log_bytes_flushed"]
	p["uncheckpointed_bytes"] = p["log_bytes_written"] - p["last_checkpoint"]
}

func parseProcesslist(state string, p map[string]float64) {

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

func setMap(p map[string]float64, key, src string) {
	v, err := atof(src)
	if err != nil {
		return
	}
	p[key] = v
}

func setMapIfEmpty(p map[string]float64, key, src string) {
	if _, ok := p[key]; ok {
		return
	}
	setMap(p, key, src)
}

func increaseMap(p map[string]float64, key string, src string) {
	val, err := atof(src)
	if err != nil {
		val = 0
	}
	_, exists := p[key]
	if !exists {
		p[key] = val
		return
	}
	p[key] = p[key] + val
}

func makeBigint(hi string, lo string, hexFormat bool) int64 {
	if lo == "" {
		if hexFormat {
			val, _ := strconv.ParseInt(hi, 16, 64)
			return val
		}
		val, _ := strconv.ParseInt(hi, 10, 64)
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

	val := hiVal<<32 + loVal

	return val
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "3306", "Port")
	optSocket := flag.String("socket", "", "Path to unix socket")
	optUser := flag.String("username", "root", "Username")
	optPass := flag.String("password", os.Getenv("MYSQL_PASSWORD"), "Password")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optInnoDB := flag.Bool("disable_innodb", false, "Disable InnoDB metrics")
	optMetricKeyPrefix := flag.String("metric-key-prefix", "mysql", "metric key prefix")
	optEnableExtended := flag.Bool("enable_extended", false, "Enable Extended metrics")
	optDebug := flag.Bool("debug", false, "Print debugging logs to stderr")
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
	mysql.Debug = *optDebug
	helper := mp.NewMackerelPlugin(&mysql)
	helper.Tempfile = *optTempfile
	helper.Run()
}
