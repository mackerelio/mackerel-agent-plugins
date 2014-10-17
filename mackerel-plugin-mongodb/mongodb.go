package main

import (
	"errors"
	"flag"
	"fmt"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"strconv"
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

var metricPlace map[string][]string = map[string][]string{
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
				return 0, errors.New("Cannot handle as a hash")
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
	Url string
}

func (m MongoDBPlugin) FetchMetrics() (map[string]float64, error) {
	session, err := mgo.Dial(m.Url)
	if err != nil {
		return nil, err
	}

	defer session.Close()
	serverStatus := bson.M{}
	if err := session.Run("serverStatus", &serverStatus); err != nil {
		return nil, err
	}

	stat := make(map[string]float64)
	for k, v := range metricPlace {
		val, err := GetFloatValue(serverStatus, v)
		if err != nil {
			return nil, err
		}

		stat[k] = val
	}

	return stat, err
}

func (m MongoDBPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "27017", "Port")
	optUser := flag.String("username", "", "Username")
	optPass := flag.String("password", "", "Password")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var mongodb MongoDBPlugin
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
