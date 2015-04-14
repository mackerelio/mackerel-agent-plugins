package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"mysql.cmd": mp.Graphs{
		Label: "MySQL Command",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Com_insert", Label: "Insert", Diff: true, Stacked: true},
			mp.Metrics{Name: "Com_select", Label: "Select", Diff: true, Stacked: true},
			mp.Metrics{Name: "Com_update", Label: "Update", Diff: true, Stacked: true},
			mp.Metrics{Name: "Com_update_multi", Label: "Update Multi", Diff: true, Stacked: true},
			mp.Metrics{Name: "Com_delete", Label: "Delete", Diff: true, Stacked: true},
			mp.Metrics{Name: "Com_delete_multi", Label: "Delete Multi", Diff: true, Stacked: true},
			mp.Metrics{Name: "Com_replace", Label: "Replace", Diff: true, Stacked: true},
			mp.Metrics{Name: "Com_set_option", Label: "Set Option", Diff: true, Stacked: true},
			mp.Metrics{Name: "Qcache_hits", Label: "Query Cache Hits", Diff: true, Stacked: false},
			mp.Metrics{Name: "Questions", Label: "Questions", Diff: true, Stacked: false},
		},
	},
	"mysql.join": mp.Graphs{
		Label: "MySQL Join/Scan",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Select_full_join", Label: "Select Full JOIN", Diff: true, Stacked: false},
			mp.Metrics{Name: "Select_full_range_join", Label: "Select Full Range JOIN", Diff: true, Stacked: false},
			mp.Metrics{Name: "Select_scan", Label: "Select SCAN", Diff: true, Stacked: false},
			mp.Metrics{Name: "Sort_scan", Label: "Sort SCAN", Diff: true, Stacked: false},
		},
	},
	"mysql.threads": mp.Graphs{
		Label: "MySQL Threads",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Max_used_connections", Label: "Max used connections", Diff: false, Stacked: false},
			mp.Metrics{Name: "Threads_connected", Label: "Connected", Diff: false, Stacked: false},
			mp.Metrics{Name: "Threads_running", Label: "Running", Diff: false, Stacked: false},
			mp.Metrics{Name: "Threads_cached", Label: "Cached", Diff: false, Stacked: false},
		},
	},
	"mysql.connections": mp.Graphs{
		Label: "MySQL Connections",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Connections", Label: "Connections", Diff: true, Stacked: false},
			mp.Metrics{Name: "Thread_created", Label: "Created Threads", Diff: true, Stacked: false},
			mp.Metrics{Name: "Aborted_clients", Label: "Aborted Clients", Diff: true, Stacked: false},
			mp.Metrics{Name: "Aborted_connects", Label: "Aborted Connects", Diff: true, Stacked: false},
		},
	},
	"mysql.seconds_behind_master": mp.Graphs{
		Label: "MySQL Slave status",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Seconds_Behind_Master", Label: "Seconds Behind Master", Diff: false, Stacked: false},
		},
	},
	"mysql.table_locks": mp.Graphs{
		Label: "MySQL Table Locks/Slow Queries",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Table_locks_immediate", Label: "Table Locks Immediate", Diff: true, Stacked: false},
			mp.Metrics{Name: "Table_locks_waited", Label: "Table Locks Waited", Diff: true, Stacked: false},
			mp.Metrics{Name: "Slow_queries", Label: "Slow Queries", Diff: true, Stacked: false},
		},
	},
	"mysql.traffic": mp.Graphs{
		Label: "MySQL Traffic",
		Unit:  "bytes/sec",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Bytes_sent", Label: "Sent Bytes", Diff: true, Stacked: false},
			mp.Metrics{Name: "Bytes_received", Label: "Received Bytes", Diff: true, Stacked: false},
		},
	},
}

type MySQLPlugin struct {
	Target        string
	Tempfile      string
	Username      string
	Password      string
	DisableInnoDB bool
}

func (m MySQLPlugin) FetchMetrics() (map[string]float64, error) {
	db := mysql.New("tcp", "", m.Target, m.Username, m.Password, "")
	err := db.Connect()
	if err != nil {
		log.Fatalln("FetchMetrics (DB Connect): ", err)
		return nil, err
	}
	defer db.Close()

	stat := make(map[string]float64)

	rows, _, err := db.Query("show /*!50002 global */ status")
	if err != nil {
		log.Fatalln("FetchMetrics (Status): ", err)
		return nil, err
	}
	for _, row := range rows {
		Variable_name := string(row[0].([]byte))
		if err != nil {
			log.Println("FetchMetrics (Status Fetch): ", err)
		}
		//fmt.Println(Variable_name, Value)
		stat[Variable_name], _ = _atof(string(row[1].([]byte)))
	}

	if m.DisableInnoDB != true {
		row, _, err := db.QueryFirst("SHOW /*!50000 ENGINE*/ INNODB STATUS")
		if err != nil {
			log.Fatalln("FetchMetrics (InnoDB Status): ", err)
			return nil, err
		}
		err = parseInnodbStatus(string(row[2].([]byte)), &stat)

		rows, _, err = db.Query("SHOW VARIABLES")
		if err != nil {
			log.Fatalln("FetchMetrics (Variables): ", err)
			return nil, err
		}
		for _, row := range rows {
			Variable_name := string(row[0].([]byte))
			if err != nil {
				log.Println("FetchMetrics (Fetch Variables): ", err)
			}
			//fmt.Println(Variable_name, Value)
			stat[Variable_name], _ = _atof(string(row[1].([]byte)))
		}
	}

	rows, res, err := db.Query("show slave status")
	if err != nil {
		log.Fatalln("FetchMetrics (Slave Status): ", err)
		return nil, err
	}
	for _, row := range rows {
		idx := res.Map("Seconds_Behind_Master")
		Value := row.Int(idx)
		stat["Seconds_Behind_Master"] = float64(Value)
	}
	return stat, err
}

func (m MySQLPlugin) GraphDefinition() map[string](mp.Graphs) {
	if m.DisableInnoDB != true {
		setInnoDBMetrics()
	}

	return graphdef
}

func setInnoDBMetrics() {
	graphdef["mysql.innodb_rows"] = mp.Graphs{
		Label: "mysql.innodb Rows",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Innodb_rows_read", Label: "Read", Diff: true, Stacked: false},
			mp.Metrics{Name: "Innodb_rows_inserted", Label: "Inserted", Diff: true, Stacked: false},
			mp.Metrics{Name: "Innodb_rows_updated", Label: "Updated", Diff: true, Stacked: false},
			mp.Metrics{Name: "Innodb_rows_deleted", Label: "Deleted", Diff: true, Stacked: false},
		},
	}
	graphdef["mysql.innodb_row_lock_time"] = mp.Graphs{
		Label: "mysql.innodb Row Lock Time",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Innodb_row_lock_time", Label: "Lock Time", Diff: true, Stacked: false},
		},
	}
	graphdef["mysql.innodb_row_lock_waits"] = mp.Graphs{
		Label: "mysql.innodb Row Lock Waits",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Innodb_row_lock_waits", Label: "Lock Waits", Diff: true, Stacked: false},
		},
	}
	graphdef["mysql.innodb_adaptive_hash_index"] = mp.Graphs{
		Label: "mysql.innodb Adaptive Hash Index",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "hash_index_cells_total", Label: "Hash Index Cells Total", Diff: false, Stacked: false},
			mp.Metrics{Name: "hash_index_cells_used", Label: "Hash Index Cells Used", Diff: false, Stacked: false},
		},
	}
	graphdef["mysql.innodb_buffer_pool_read"] = mp.Graphs{
		Label: "mysql.innodb Buffer Pool Read (/sec)",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "read_ahead", Label: "Pages Read Ahead", Diff: false, Stacked: false},
			mp.Metrics{Name: "read_evicted", Label: "Evicted Without Access", Diff: false, Stacked: false},
			mp.Metrics{Name: "read_random_ahead", Label: "Random Read Ahead", Diff: false, Stacked: false},
		},
	}
	graphdef["mysql.innodb_buffer_pool_activity"] = mp.Graphs{
		Label: "mysql.innodb Buffer Pool Activity (Pages)",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "pages_created", Label: "Created", Diff: true, Stacked: false},
			mp.Metrics{Name: "pages_read", Label: "Read", Diff: true, Stacked: false},
			mp.Metrics{Name: "pages_written", Label: "Written", Diff: true, Stacked: false},
		},
	}
	graphdef["mysql.innodb_buffer_pool_efficiency"] = mp.Graphs{
		Label: "mysql.innodb Buffer Pool Efficiency",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Innodb_buffer_pool_reads", Label: "Reads", Diff: true, Stacked: false},
			mp.Metrics{Name: "Innodb_buffer_pool_read_requests", Label: "Read Requests", Diff: true, Stacked: false},
		},
	}
	graphdef["mysql.innodb_buffer_pool"] = mp.Graphs{
		Label: "mysql.innodb Buffer Pool (Pages)",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "pool_size", Label: "Pool Size", Diff: false, Stacked: false},
			mp.Metrics{Name: "database_pages", Label: "Used", Diff: false, Stacked: true},
			mp.Metrics{Name: "free_pages", Label: "Free", Diff: false, Stacked: true},
			mp.Metrics{Name: "modified_pages", Label: "Modified", Diff: false, Stacked: false},
		},
	}
	graphdef["mysql.innodb_checkpoint_age"] = mp.Graphs{
		Label: "mysql.innodb Checkpoint Age",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "uncheckpointed_bytes", Label: "Uncheckpointed", Diff: false, Stacked: false},
		},
	}
	graphdef["mysql.innodb_current_lock_waits"] = mp.Graphs{
		Label: "mysql.innodb Current Lock Waits (secs)",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "innodb_lock_wait_secs", Label: "Innodb Lock Wait", Diff: false, Stacked: false},
		},
	}
	graphdef["mysql.innodb_io"] = mp.Graphs{
		Label: "mysql.innodb I/O",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "file_reads", Label: "File Reads", Diff: true, Stacked: false},
			mp.Metrics{Name: "file_writes", Label: "File Writes", Diff: true, Stacked: false},
			mp.Metrics{Name: "file_fsyncs", Label: "File fsyncs", Diff: true, Stacked: false},
			mp.Metrics{Name: "log_writes", Label: "Log Writes", Diff: true, Stacked: false},
		},
	}
	graphdef["mysql.innodb_io_pending"] = mp.Graphs{
		Label: "mysql.innodb I/O Pending",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "pending_normal_aio_reads", Label: "Normal AIO Reads", Diff: false, Stacked: false},
			mp.Metrics{Name: "pending_normal_aio_writes", Label: "Normal AIO Writes", Diff: false, Stacked: false},
			mp.Metrics{Name: "pending_ibuf_aio_reads", Label: "InnoDB Buffer AIO Reads", Diff: false, Stacked: false},
			mp.Metrics{Name: "pending_aio_log_ios", Label: "AIO Log IOs", Diff: false, Stacked: false},
			mp.Metrics{Name: "pending_aio_sync_ios", Label: "AIO Sync IOs", Diff: false, Stacked: false},
			mp.Metrics{Name: "pending_log_flushes", Label: "Log Flushes", Diff: false, Stacked: false},
			mp.Metrics{Name: "pending_buf_pool_flushes", Label: "Buffer Pool Flushes", Diff: false, Stacked: false},
			mp.Metrics{Name: "pending_log_writes", Label: "Log Writes", Diff: false, Stacked: false},
			mp.Metrics{Name: "pending_chkp_writes", Label: "Checkpoint Writes", Diff: false, Stacked: false},
		},
	}
	graphdef["mysql.innodb_insert_buffer"] = mp.Graphs{
		Label: "mysql.innodb Insert Buffer",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "ibuf_inserts", Label: "Inserts", Diff: true, Stacked: false},
			mp.Metrics{Name: "ibuf_merges", Label: "Merges", Diff: true, Stacked: false},
			mp.Metrics{Name: "ibuf_merged", Label: "Merged", Diff: true, Stacked: false},
		},
	}
	graphdef["mysql.innodb_insert_buffer_usage"] = mp.Graphs{
		Label: "mysql.innodb Insert Buffer Usage (Cells)",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "ibuf_cell_count", Label: "Cell Count", Diff: false, Stacked: false},
			mp.Metrics{Name: "ibuf_used_cells", Label: "Used", Diff: false, Stacked: true},
			mp.Metrics{Name: "ibuf_free_cells", Label: "Free", Diff: false, Stacked: true},
		},
	}
	graphdef["mysql.innodb_lock_structures"] = mp.Graphs{
		Label: "mysql.innodb Lock Structures",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "innodb_lock_structs", Label: "Structures", Diff: false, Stacked: false},
		},
	}
	graphdef["mysql.innodb_log"] = mp.Graphs{
		Label: "mysql.innodb Log",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "log_bytes_written", Label: "Written", Diff: true, Stacked: false},
			mp.Metrics{Name: "log_bytes_flushed", Label: "Flushed", Diff: true, Stacked: false},
			mp.Metrics{Name: "unflushed_log", Label: "Unflushed", Diff: false, Stacked: false},
			mp.Metrics{Name: "innodb_log_buffer_size", Label: "Buffer Size", Diff: false, Stacked: true},
		},
	}
	graphdef["mysql.innodb_memory_allocation"] = mp.Graphs{
		Label: "mysql.innodb Memory Allocation",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "additional_pool_alloc", Label: "Additional Pool Allocated", Diff: false, Stacked: false},
			mp.Metrics{Name: "total_mem_alloc", Label: "Total Memory Allocated", Diff: false, Stacked: false},
		},
	}
	graphdef["mysql.innodb_semaphores"] = mp.Graphs{
		Label: "mysql.innodb Semaphores",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "spin_waits", Label: "Spin Waits", Diff: true, Stacked: false},
			mp.Metrics{Name: "spin_rounds", Label: "Spin Rounds", Diff: true, Stacked: false},
			mp.Metrics{Name: "os_waits", Label: "OS Waits", Diff: true, Stacked: false},
		},
	}
	graphdef["mysql.innodb_tables_in_use"] = mp.Graphs{
		Label: "mysql.innodb Tables In Use",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "innodb_locked_tables", Label: "Table in Use", Diff: false, Stacked: false},
			mp.Metrics{Name: "innodb_tables_in_use", Label: "Locked Tables", Diff: false, Stacked: false},
		},
	}
	graphdef["mysql.innodb_transactions_active_locked"] = mp.Graphs{
		Label: "mysql.innodb Transactions Active/Locked",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "current_transactions", Label: "Current", Diff: false, Stacked: false},
			mp.Metrics{Name: "active_transactions", Label: "Active", Diff: false, Stacked: false},
			mp.Metrics{Name: "locked_transactions", Label: "Locked", Diff: false, Stacked: false},
			mp.Metrics{Name: "read_views", Label: "Read Views", Diff: false, Stacked: false},
		},
	}
	graphdef["mysql.innodb_transactions"] = mp.Graphs{
		Label: "mysql.innodb Transactions",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "history_list", Label: "History List", Diff: false, Stacked: false},
			mp.Metrics{Name: "innodb_transactions", Label: "InnoDB Transactions", Diff: true, Stacked: false},
		},
	}
}

func parseInnodbStatus(str string, p *map[string]float64) error {

	is_transaction := false
	prev_line := ""

	for _, line := range strings.Split(str, "\n") {
		record := strings.Fields(line)

		// Innodb Semaphores
		if strings.Index(line, "Mutex spin waits") == 0 {
			_increase_map(p, "spin_waits", record[3])
			_increase_map(p, "spin_rounds", record[5])
			_increase_map(p, "os_waits", record[8])
			continue
		}
		if strings.Index(line, "RW-shared spins") == 0 && strings.Index(line, ";") > 0 {
			// 5.5, 5.6
			_increase_map(p, "spin_waits", record[2])
			_increase_map(p, "os_waits", record[5])
			_increase_map(p, "spin_waits", record[8])
			_increase_map(p, "os_waits", record[11])
			continue
		}
		if strings.Index(line, "RW-shared spins") == 0 && strings.Index(line, "; RW-excl spins") < 0 {
			// 5.1
			_increase_map(p, "spin_waits", record[2])
			_increase_map(p, "os_waits", record[7])
			continue
		}
		if strings.Index(line, "RW-excl spins") == 0 {
			// 5.5, 5.6
			_increase_map(p, "spin_waits", record[2])
			_increase_map(p, "os_waits", record[7])
			continue
		}
		if strings.Index(line, "seconds the semaphore:") > 0 {
			_increase_map(p, "innodb_sem_waits", "1")
			wait, _ := _atof(record[9])
			wait = wait * 1000
			_increase_map(p, "innodb_sem_wait_time_ms", fmt.Sprintf("%.f", wait))
			continue
		}

		// Innodb Transactions
		if strings.Index(line, "Trx id counter") == 0 {
			lo_val := ""
			if len(record) >= 5 {
				lo_val = record[4]
			}
			val := _make_bigint(record[3], lo_val)
			_increase_map(p, "innodb_transactions", fmt.Sprintf("%d", val))
			is_transaction = true
			continue
		}
		if strings.Index(line, "Purge done for trx") == 0 {
			if record[7] == "undo" {
				record[7] = ""
			}
			val := _make_bigint(record[6], record[7])
			trx := (*p)["innodb_transactions"] - float64(val)
			_increase_map(p, "unpurged_txns", fmt.Sprintf("%.f", trx))
			continue
		}
		if strings.Index(line, "History list length") == 0 {
			_increase_map(p, "history_list", record[3])
			continue
		}
		if is_transaction && strings.Index(line, "---TRANSACTION") == 0 {
			_increase_map(p, "current_transactions", "1")
			if strings.Index(line, "ACTIVE") > 0 {
				_increase_map(p, "active_transactions", "1")
			}
			continue
		}
		if is_transaction && strings.Index(line, "------- TRX HAS BEEN") == 0 {
			_increase_map(p, "innodb_lock_wait_secs", "1")
			continue
		}
		if strings.Index(line, "read views open inside InnoDB") > 0 {
			(*p)["read_views"], _ = _atof(record[0])
			continue
		}
		if strings.Index(line, "------- TRX HAS BEEN") == 0 {
			_increase_map(p, "innodb_tables_in_use", record[4])
			_increase_map(p, "innodb_locked_tables", record[6])
			continue
		}
		if is_transaction && strings.Index(line, "lock struct(s)") == 0 {
			if strings.Index(line, "LOCK WAIT") > 0 {
				_increase_map(p, "innodb_lock_structs", record[2])
				_increase_map(p, "locked_transactions", "1")
			} else {
				_increase_map(p, "innodb_lock_structs", record[0])
			}
			continue
		}

		// File I/O
		if strings.Index(line, " OS file reads, ") > 0 {
			(*p)["file_reads"], _ = _atof(record[0])
			(*p)["file_writes"], _ = _atof(record[4])
			(*p)["file_fsyncs"], _ = _atof(record[8])
			continue
		}
		if strings.Index(line, "Pending normal aio reads:") == 0 {
			(*p)["pending_normal_aio_reads"], _ = _atof(record[4])
			(*p)["pending_normal_aio_writes"], _ = _atof(record[7])
			continue
		}
		if strings.Index(line, "ibuf aio reads") == 0 {
			(*p)["pending_ibuf_aio_reads"], _ = _atof(record[3])
			(*p)["pending_aio_log_ios"], _ = _atof(record[6])
			(*p)["pending_aio_sync_ios"], _ = _atof(record[9])
			continue
		}
		if strings.Index(line, "Pending flushes (fsync)") == 0 {
			(*p)["pending_log_flushes"], _ = _atof(record[4])
			(*p)["pending_buf_pool_flushes"], _ = _atof(record[7])
			continue
		}

		// Insert Buffer and Adaptive Hash Index
		if strings.Index(line, "Ibuf for space 0: size ") == 0 {
			(*p)["ibuf_used_cells"], _ = _atof(record[5])
			(*p)["ibuf_free_cells"], _ = _atof(record[9])
			(*p)["ibuf_cell_count"], _ = _atof(record[12])
			continue
		}
		if strings.Index(line, "Ibuf: size ") == 0 {
			(*p)["ibuf_used_cells"], _ = _atof(record[2])
			(*p)["ibuf_free_cells"], _ = _atof(record[6])
			(*p)["ibuf_cell_count"], _ = _atof(record[9])
			if strings.Index(line, "merges") > 0 {
				(*p)["ibuf_merges"], _ = _atof(record[10])
			}
			continue
		}
		if strings.Index(line, ", delete mark ") > 0 && strings.Index(prev_line, "merged operations:") == 0 {
			(*p)["ibuf_inserts"], _ = _atof(record[1])
			v1, _ := _atof(record[1])
			v2, _ := _atof(record[4])
			v3, _ := _atof(record[6])
			(*p)["ibuf_merged"] = v1 + v2 + v3
			continue
		}
		if strings.Index(line, " merged recs, ") > 0 {
			(*p)["ibuf_inserts"], _ = _atof(record[0])
			(*p)["ibuf_merged"], _ = _atof(record[2])
			(*p)["ibuf_merges"], _ = _atof(record[5])
			continue
		}
		if strings.Index(line, "Hash table size ") == 0 {
			(*p)["hash_index_cells_total"], _ = _atof(record[3])
			if strings.Index(line, "used cells") > 0 {
				(*p)["hash_index_cells_used"], _ = _atof(record[6])
			} else {
				(*p)["hash_index_cells_used"] = 0
			}
			continue
		}

		// Log
		if strings.Index(line, " log i/o's done, ") > 0 {
			(*p)["log_writes"], _ = _atof(record[0])
			continue
		}
		if strings.Index(line, " pending log writes, ") > 0 {
			(*p)["pending_log_writes"], _ = _atof(record[0])
			(*p)["pending_chkp_writes"], _ = _atof(record[4])
			continue
		}
		if strings.Index(line, "Log sequence number") == 0 {
			val, _ := _atof(record[3])
			if len(record) >= 5 {
				val = float64(_make_bigint(record[3], record[4]))
			}
			(*p)["log_bytes_written"] = val
			continue
		}
		if strings.Index(line, "Log flushed up to") == 0 {
			val, _ := _atof(record[4])
			if len(record) >= 6 {
				val = float64(_make_bigint(record[4], record[5]))
			}
			(*p)["log_bytes_flushed"] = val
			continue
		}
		if strings.Index(line, "Last checkpoint at") == 0 {
			val, _ := _atof(record[3])
			if len(record) >= 5 {
				val = float64(_make_bigint(record[3], record[4]))
			}
			(*p)["last_checkpoint"] = val
			continue
		}

		// Buffer Pool and Memory
		if strings.Index(line, "Total memory allocated") == 0 && strings.Index(line, "in additional pool allocated") > 0 {
			(*p)["total_mem_alloc"], _ = _atof(record[3])
			(*p)["additional_pool_alloc"], _ = _atof(record[8])
			continue
		}
		if strings.Index(line, "Adaptive hash index ") == 0 {
			(*p)["adaptive_hash_memory"], _ = _atof(record[3])
			continue
		}
		if strings.Index(line, "Page hash           ") == 0 {
			(*p)["page_hash_memory"], _ = _atof(record[2])
			continue
		}
		if strings.Index(line, "Dictionary cache    ") == 0 {
			(*p)["dictionary_cache_memory"], _ = _atof(record[2])
			continue
		}
		if strings.Index(line, "File system         ") == 0 {
			(*p)["file_system_memory"], _ = _atof(record[2])
			continue
		}
		if strings.Index(line, "Lock system         ") == 0 {
			(*p)["lock_system_memory"], _ = _atof(record[2])
			continue
		}
		if strings.Index(line, "Recovery system     ") == 0 {
			(*p)["recovery_system_memory"], _ = _atof(record[2])
			continue
		}
		if strings.Index(line, "Threads             ") == 0 {
			(*p)["thread_hash_memory"], _ = _atof(record[1])
			continue
		}
		if strings.Index(line, "innodb_io_pattern   ") == 0 {
			(*p)["innodb_io_pattern_memory"], _ = _atof(record[1])
			continue
		}
		if strings.Index(line, "Buffer pool size ") == 0 {
			(*p)["pool_size"], _ = _atof(record[3])
			continue
		}
		if strings.Index(line, "Free buffers") == 0 {
			(*p)["free_pages"], _ = _atof(record[2])
			continue
		}
		if strings.Index(line, "Database pages") == 0 {
			(*p)["database_pages"], _ = _atof(record[2])
			continue
		}
		if strings.Index(line, "Modified db pages") == 0 {
			(*p)["modified_pages"], _ = _atof(record[3])
			continue
		}
		if strings.Index(line, "Pages read ahead") == 0 {
			(*p)["read_ahead"], _ = _atof(record[3])
			(*p)["read_evicted"], _ = _atof(record[7])
			(*p)["read_random_ahead"], _ = _atof(record[11])
			continue
		}
		if strings.Index(line, "Pages read") == 0 {
			(*p)["pages_read"], _ = _atof(record[2])
			(*p)["pages_created"], _ = _atof(record[4])
			(*p)["pages_written"], _ = _atof(record[6])
			continue
		}

		// Row Operations
		if strings.Index(line, "Number of rows inserted") == 0 {
			(*p)["rows_inserted"], _ = _atof(record[4])
			(*p)["rows_updated"], _ = _atof(record[6])
			(*p)["rows_deleted"], _ = _atof(record[8])
			(*p)["rows_read"], _ = _atof(record[10])
			continue
		}
		if strings.Index(line, " queries inside InnoDB, ") == 0 {
			(*p)["queries_inside"], _ = _atof(record[0])
			(*p)["queries_queued"], _ = _atof(record[4])
			continue
		}

		// for next loop
		prev_line = line
	}

	// finalize
	(*p)["queries_queued"] = (*p)["log_bytes_written"] - (*p)["log_bytes_flushed"]
	(*p)["uncheckpointed_bytes"] = (*p)["log_bytes_written"] - (*p)["last_checkpoint"]

	return nil
}

// atof
func _atof(str string) (float64, error) {
	str = strings.Replace(str, ",", "", -1)
	str = strings.Replace(str, ";", "", -1)
	str = strings.Replace(str, "/s", "", -1)
	str = strings.Trim(str, " ")
	return strconv.ParseFloat(str, 64)
}

func _increase_map(p *map[string]float64, key string, src string) {
	val, err := _atof(src)
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

func _increase(src *float64, data float64) {
	*src = *src + data
}

func _make_bigint(hi string, lo string) int64 {
	if lo == "" {
		val, _ := strconv.ParseInt(hi, 16, 64)
		return val
	}

	var hi_val int64 = 0
	var lo_val int64 = 0
	if hi != "" {
		hi_val, _ = strconv.ParseInt(hi, 10, 64)
	}
	if lo != "" {
		lo_val, _ = strconv.ParseInt(lo, 10, 64)
	}

	val := hi_val * lo_val

	return val
}

func main() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "3306", "Port")
	optUser := flag.String("username", "root", "Username")
	optPass := flag.String("password", "", "Password")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optInnoDB := flag.Bool("disable_innodb", false, "Disable InnoDB metrics")
	flag.Parse()

	var mysql MySQLPlugin

	mysql.Target = fmt.Sprintf("%s:%s", *optHost, *optPort)
	mysql.Username = *optUser
	mysql.Password = *optPass
	mysql.DisableInnoDB = *optInnoDB
	helper := mp.NewMackerelPlugin(mysql)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-mysql-%s-%s", *optHost, *optPort)
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
