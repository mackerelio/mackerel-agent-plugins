package main

import (
	"flag"
	"fmt"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"log"
	"os"
	"strconv"
)

var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"mysql.connections": mp.Graphs{
		Label: "MySQL Connections",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Key: "curr_connections", Label: "Connections", Diff: false},
		},
	},
	"mysql.cmd": mp.Graphs{
		Label: "MySQL Command",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Key: "Com_insert", Label: "Insert", Diff: true},
			mp.Metrics{Key: "Com_select", Label: "Select", Diff: true},
			mp.Metrics{Key: "Com_update", Label: "Update", Diff: true},
			mp.Metrics{Key: "Com_update_multi", Label: "Update Multi", Diff: true},
			mp.Metrics{Key: "Com_delete", Label: "Delete", Diff: true},
			mp.Metrics{Key: "Com_delete_multi", Label: "Delete Multi", Diff: true},
			mp.Metrics{Key: "Com_replace", Label: "Replace", Diff: true},
			mp.Metrics{Key: "Com_set_option", Label: "Set Option", Diff: true},
			mp.Metrics{Key: "Qcache_hits", Label: "Query Cache Hits", Diff: true},
			mp.Metrics{Key: "Questions", Label: "Questions", Diff: true},
		},
	},
	"mysql.join": mp.Graphs{
		Label: "MySQL Join/Scan",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Key: "Select_full_join", Label: "Select Full JOIN", Diff: true},
			mp.Metrics{Key: "Select_full_range_join", Label: "Select Full Range JOIN", Diff: true},
			mp.Metrics{Key: "Select_scan", Label: "Select SCAN", Diff: true},
			mp.Metrics{Key: "Sort_scan", Label: "Sort SCAN", Diff: true},
		},
	},
	"mysql.Threads": mp.Graphs{
		Label: "MySQL Threads",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Key: "Max_used_connections", Label: "Max used connections", Diff: true},
			mp.Metrics{Key: "Threads_connected", Label: "Connected", Diff: true},
			mp.Metrics{Key: "Threads_running", Label: "Running", Diff: true},
			mp.Metrics{Key: "Threads_cached", Label: "Cached", Diff: true},
		},
	},
}

type MySQLPlugin struct {
	Target   string
	Tempfile string
	Username string
	Password string
}

func (m MySQLPlugin) FetchData() (map[string]float64, error) {
	db := mysql.New("tcp", "", m.Target, m.Username, m.Password, "mysql")
	err := db.Connect()
	if err != nil {
		log.Fatalln("FetchData: ", err)
		return nil, err
	}
	defer db.Close()

	stat := make(map[string]float64)

	rows, _, err := db.Query("show /*!50002 global */ status")
	if err != nil {
		log.Fatalln("FetchData2: ", err)
		return nil, err
	}
	for _, row := range rows {
		Variable_name := string(row[0].([]byte))
		Value, err := strconv.Atoi(string(row[1].([]byte)))
		if err != nil {
			log.Println("FetchData2: ", err)
		}
		// fmt.Println(Variable_name, Value)
		stat[Variable_name] = float64(Value)
	}
	return stat, err
}

func (m MySQLPlugin) GetGraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func (m MySQLPlugin) GetTempfilename() string {
	return m.Tempfile
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
	if *optTempfile != "" {
		mysql.Tempfile = *optTempfile
	} else {
		mysql.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-mysql-%s-%s", *optHost, *optPort)
	}

	helper := mp.MackerelPlugin{mysql}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
