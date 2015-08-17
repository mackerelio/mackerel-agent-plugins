package main

import (
	"encoding/json"
	"flag"
	"fmt"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"mongodb.background_flushing": mp.Graphs{
		Label: "MongoDB Command",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "duration_ms", Label: "Duration in ms", Diff: true},
		},
	},
	"mongodb.connections": mp.Graphs{
		Label: "MongoDB Connections",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "connections_current", Label: "current"},
		},
	},
	"mongodb.index_counters.btree": mp.Graphs{
		Label: "MongoDB Index Counters Btree",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "btree_hits", Label: "hits", Diff: true},
		},
	},
	"mongodb.index_counters.btree.miss_ratio": mp.Graphs{
		Label: "MongoDB Index Counters Btree Miss-ratio",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "miss_ratio", Label: "Miss ratio"},
		},
	},
	"mongodb.opcounters": mp.Graphs{
		Label: "MongoDB opcounters",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "opcounters_insert", Label: "Insert", Diff: true},
			mp.Metrics{Name: "opcounters_query", Label: "Query", Diff: true},
			mp.Metrics{Name: "opcounters_update", Label: "Update", Diff: true},
			mp.Metrics{Name: "opcounters_delete", Label: "Delete", Diff: true},
			mp.Metrics{Name: "opcounters_getmore", Label: "Getmore", Diff: true},
			mp.Metrics{Name: "opcounters_command", Label: "Command", Diff: true},
		},
	},
}

var graphdef30 map[string](mp.Graphs) = map[string](mp.Graphs){
	"mongodb.background_flushing": mp.Graphs{
		Label: "MongoDB Command",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "duration_ms", Label: "Duration in ms", Diff: true},
		},
	},
	"mongodb.connections": mp.Graphs{
		Label: "MongoDB Connections",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "connections_current", Label: "current"},
		},
	},
	"mongodb.index_counters.btree.miss_ratio": mp.Graphs{
		Label: "MongoDB Index Counters Btree Miss-ratio",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "miss_ratio", Label: "Miss ratio"},
		},
	},
	"mongodb.opcounters": mp.Graphs{
		Label: "MongoDB opcounters",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "opcounters_insert", Label: "Insert", Diff: true},
			mp.Metrics{Name: "opcounters_query", Label: "Query", Diff: true},
			mp.Metrics{Name: "opcounters_update", Label: "Update", Diff: true},
			mp.Metrics{Name: "opcounters_delete", Label: "Delete", Diff: true},
			mp.Metrics{Name: "opcounters_getmore", Label: "Getmore", Diff: true},
			mp.Metrics{Name: "opcounters_command", Label: "Command", Diff: true},
		},
	},
}

var metricPlace22 map[string][]string = map[string][]string{
	"duration_ms":         []string{"backgroundFlushing", "total_ms"},
	"connections_current": []string{"connections", "current"},
	"btree_hits":          []string{"indexCounters", "btree", "hits"},
	"opcounters_insert":   []string{"opcounters", "insert"},
	"opcounters_query":    []string{"opcounters", "query"},
	"opcounters_update":   []string{"opcounters", "update"},
	"opcounters_delete":   []string{"opcounters", "delete"},
	"opcounters_getmore":  []string{"opcounters", "getmore"},
	"opcounters_command":  []string{"opcounters", "command"},
}

var metricPlace24 map[string][]string = map[string][]string{
	"duration_ms":         []string{"backgroundFlushing", "total_ms"},
	"connections_current": []string{"connections", "current"},
	"btree_hits":          []string{"indexCounters", "hits"},
	"opcounters_insert":   []string{"opcounters", "insert"},
	"opcounters_query":    []string{"opcounters", "query"},
	"opcounters_update":   []string{"opcounters", "update"},
	"opcounters_delete":   []string{"opcounters", "delete"},
	"opcounters_getmore":  []string{"opcounters", "getmore"},
	"opcounters_command":  []string{"opcounters", "command"},
}

// indexCounters is removed from mongodb 3.0.
// ref. http://stackoverflow.com/questions/29428793/where-is-the-indexcounter-in-db-serverstatus-on-mongodb-3-0
var metricPlace30 map[string][]string = map[string][]string{
	"duration_ms":         []string{"backgroundFlushing", "total_ms"},
	"connections_current": []string{"connections", "current"},
	"opcounters_insert":   []string{"opcounters", "insert"},
	"opcounters_query":    []string{"opcounters", "query"},
	"opcounters_update":   []string{"opcounters", "update"},
	"opcounters_delete":   []string{"opcounters", "delete"},
	"opcounters_getmore":  []string{"opcounters", "getmore"},
	"opcounters_command":  []string{"opcounters", "command"},
}

func GetFloatValue(s map[string]interface{}, keys []string) (float64, error) {
	var val float64
	sm := s
	var err error
	for i, k := range keys {
		if i+1 < len(keys) {
			switch sm[k].(type) {
			case bson.M:
				sm = sm[k].(bson.M)
			default:
				return 0, fmt.Errorf("Cannot handle as a hash for %s", k)
			}
		} else {
			val, err = strconv.ParseFloat(fmt.Sprint(sm[k]), 64)
			if err != nil {
				return 0, err
			}
		}
	}

	return val, nil
}

type MongoDBPlugin struct {
	Url     string
	Verbose bool
}

func (m MongoDBPlugin) FetchStatus() (bson.M, error) {
	session, err := mgo.Dial(m.Url)
	if err != nil {
		return nil, err
	}

	defer session.Close()
	serverStatus := bson.M{}
	if err := session.Run("serverStatus", &serverStatus); err != nil {
		return nil, err
	}
	if m.Verbose {
		str, err := json.Marshal(serverStatus)
		if err != nil {
			fmt.Println(fmt.Errorf("Marshaling error: %s", err.Error()))
		}
		fmt.Println(string(str))
	}
	return serverStatus, nil
}

func (m MongoDBPlugin) FetchMetrics() (map[string]float64, error) {
	serverStatus, err := m.FetchStatus()
	if err != nil {
		return nil, err
	}
	return m.ParseStatus(serverStatus)
}

func (m MongoDBPlugin) getVersion(serverStatus bson.M) string {
	if reflect.TypeOf(serverStatus["version"]).String() == "string" {
		version := serverStatus["version"].(string)
		return version
	}
	return ""
}

func (m MongoDBPlugin) ParseStatus(serverStatus bson.M) (map[string]float64, error) {
	stat := make(map[string]float64)
	metricPlace := &metricPlace22
	version := m.getVersion(serverStatus)
	if strings.HasPrefix(version, "2.4") {
		metricPlace = &metricPlace24
	} else if strings.HasPrefix(version, "2.6") {
		metricPlace = &metricPlace24
	} else if strings.HasPrefix(version, "3.0") {
		metricPlace = &metricPlace30
	}

	for k, v := range *metricPlace {
		val, err := GetFloatValue(serverStatus, v)
		if err != nil {
			return nil, err
		}

		stat[k] = val
	}

	return stat, nil
}

func (m MongoDBPlugin) GraphDefinition() map[string](mp.Graphs) {
	serverStatus, err := m.FetchStatus()
	if err != nil {
		return graphdef
	}
	version := m.getVersion(serverStatus)
	if strings.HasPrefix(version, "3.0") {
		return graphdef30
	}
	return graphdef
}

func main() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "27017", "Port")
	optUser := flag.String("username", "", "Username")
	optPass := flag.String("password", "", "Password")
	optVerbose := flag.Bool("v", false, "Verbose mode")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var mongodb MongoDBPlugin
	mongodb.Verbose = *optVerbose
	if *optUser == "" && *optPass == "" {
		mongodb.Url = fmt.Sprintf("mongodb://%s:%s", *optHost, *optPort)
	} else {
		mongodb.Url = fmt.Sprintf("mongodb://%s:%s@%s:%s", *optUser, *optPass, *optHost, *optPort)
	}

	helper := mp.NewMackerelPlugin(mongodb)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-mongodb-%s-%s", *optHost, *optPort)
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
