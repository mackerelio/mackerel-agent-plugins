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
	"mysql.innodb_rows": mp.Graphs{
		Label: "MySQL InnoDB Rows",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Innodb_rows_read", Label: "Read", Diff: true, Stacked: false},
			mp.Metrics{Name: "Innodb_rows_inserted", Label: "Inserted", Diff: true, Stacked: false},
			mp.Metrics{Name: "Innodb_rows_updated", Label: "Updated", Diff: true, Stacked: false},
			mp.Metrics{Name: "Innodb_rows_deleted", Label: "Deleted", Diff: true, Stacked: false},
		},
	},
	"mysql.innodb_row_lock_time": mp.Graphs{
		Label: "MySQL InnoDB Row Lock Time",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Innodb_row_lock_time", Label: "Lock Time", Diff: true, Stacked: false},
		},
	},
	"mysql.innodb_row_lock_waits": mp.Graphs{
		Label: "MySQL InnoDB Row Lock Waits",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Innodb_row_lock_waits", Label: "Lock Waits", Diff: true, Stacked: false},
		},
	},
}

type MySQLPlugin struct {
	Target   string
	Tempfile string
	Username string
	Password string
}

func (m MySQLPlugin) FetchMetrics() (map[string]float64, error) {
	db := mysql.New("tcp", "", m.Target, m.Username, m.Password, "")
	err := db.Connect()
	if err != nil {
		log.Fatalln("FetchMetrics: ", err)
		return nil, err
	}
	defer db.Close()

	stat := make(map[string]float64)

	rows, _, err := db.Query("show /*!50002 global */ status")
	if err != nil {
		log.Fatalln("FetchMetrics: ", err)
		return nil, err
	}
	for _, row := range rows {
		Variable_name := string(row[0].([]byte))
		if err != nil {
			log.Println("FetchMetrics: ", err)
		}
		//fmt.Println(Variable_name, Value)
		stat[Variable_name], _ = _atof(string(row[1].([]byte)))
	}

	row, _, err := db.QueryFirst("SHOW /*!50000 ENGINE*/ INNODB STATUS")
	if err != nil {
		log.Fatalln("FetchMetrics: ", err)
		return nil, err
	}
	err = parseInnodbStatus(string(row[2].([]byte)), &stat)

	rows, res, err := db.Query("show slave status")
	if err != nil {
		log.Fatalln("FetchMetrics: ", err)
		return nil, err
	}
	for _, row = range rows {
		idx := res.Map("Seconds_Behind_Master")
		stat["Seconds_Behind_Master"], _ = _atof(string(idx))
	}
	return stat, err
}

func (m MySQLPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func parseInnodbStatus(str string, p *map[string]float64) error {

	is_transaction := false

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
		if strings.Index(line, "Pending normal aio reads:") > 0 {
			(*p)["pending_normal_aio_reads"], _ = _atof(record[4])
			(*p)["pending_normal_aio_writes"], _ = _atof(record[7])
			continue
		}
		if strings.Index(line, "ibuf aio reads") > 0 {
			(*p)["pending_ibuf_aio_reads"], _ = _atof(record[3])
			(*p)["pending_aio_log_ios"], _ = _atof(record[6])
			(*p)["pending_aio_sync_ios"], _ = _atof(record[9])
			continue
		}
		if strings.Index(line, "Pending flushes (fsync)") > 0 {
			(*p)["pending_log_flushes"], _ = _atof(record[4])
			(*p)["pending_buf_pool_flushes"], _ = _atof(record[7])
			continue
		}
	}

	return nil
}

// atof
func _atof(str string) (float64, error) {
	str = strings.Replace(str, ",", "", -1)
	str = strings.Replace(str, ";", "", -1)
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
	flag.Parse()

	var mysql MySQLPlugin

	mysql.Target = fmt.Sprintf("%s:%s", *optHost, *optPort)
	mysql.Username = *optUser
	mysql.Password = *optPass
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
