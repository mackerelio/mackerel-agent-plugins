package mpmongodb

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/hashicorp/go-version"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/golib/logging"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var logger = logging.GetLogger("metrics.plugin.mongodb")

func (m MongoDBPlugin) defaultGraphdef() map[string]mp.Graphs {
	labelPrefix := m.LabelPrefix()

	return map[string]mp.Graphs{
		"background_flushing": {
			Label: labelPrefix + " Command",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "duration_ms", Label: "Duration in ms", Diff: true, Type: "uint64"},
			},
		},
		"connections": {
			Label: labelPrefix + " Connections",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "connections_current", Label: "current"},
			},
		},
		"index_counters.btree": {
			Label: labelPrefix + " Index Counters Btree",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "btree_hits", Label: "hits", Diff: true, Type: "uint64"},
			},
		},
		"opcounters": {
			Label: labelPrefix + " opcounters",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "opcounters_insert", Label: "Insert", Diff: true, Type: "uint64"},
				{Name: "opcounters_query", Label: "Query", Diff: true, Type: "uint64"},
				{Name: "opcounters_update", Label: "Update", Diff: true, Type: "uint64"},
				{Name: "opcounters_delete", Label: "Delete", Diff: true, Type: "uint64"},
				{Name: "opcounters_getmore", Label: "Getmore", Diff: true, Type: "uint64"},
				{Name: "opcounters_command", Label: "Command", Diff: true, Type: "uint64"},
			},
		},
	}
}

func (m MongoDBPlugin) graphdef30() map[string]mp.Graphs {
	labelPrefix := m.LabelPrefix()

	return map[string]mp.Graphs{
		"background_flushing": {
			Label: labelPrefix + " Command",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "duration_ms", Label: "Duration in ms", Diff: true, Type: "uint64"},
			},
		},
		"connections": {
			Label: labelPrefix + " Connections",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "connections_current", Label: "current"},
			},
		},
		"opcounters": {
			Label: labelPrefix + " opcounters",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "opcounters_insert", Label: "Insert", Diff: true, Type: "uint64"},
				{Name: "opcounters_query", Label: "Query", Diff: true, Type: "uint64"},
				{Name: "opcounters_update", Label: "Update", Diff: true, Type: "uint64"},
				{Name: "opcounters_delete", Label: "Delete", Diff: true, Type: "uint64"},
				{Name: "opcounters_getmore", Label: "Getmore", Diff: true, Type: "uint64"},
				{Name: "opcounters_command", Label: "Command", Diff: true, Type: "uint64"},
			},
		},
	}
}

// Adapt to version 3.2 or higher.
// Check in version 3.6.
func (m MongoDBPlugin) graphdef32() map[string]mp.Graphs {
	labelPrefix := m.LabelPrefix()

	return map[string]mp.Graphs{
		"connections": {
			Label: labelPrefix + " Connections",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "connections_current", Label: "current"},
			},
		},
		"opcounters": {
			Label: labelPrefix + " opcounters",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "opcounters_insert", Label: "Insert", Diff: true, Type: "uint64"},
				{Name: "opcounters_query", Label: "Query", Diff: true, Type: "uint64"},
				{Name: "opcounters_update", Label: "Update", Diff: true, Type: "uint64"},
				{Name: "opcounters_delete", Label: "Delete", Diff: true, Type: "uint64"},
				{Name: "opcounters_getmore", Label: "Getmore", Diff: true, Type: "uint64"},
				{Name: "opcounters_command", Label: "Command", Diff: true, Type: "uint64"},
			},
		},
	}
}

var metricPlace22 = map[string][]string{
	"duration_ms":         {"backgroundFlushing", "total_ms"},
	"connections_current": {"connections", "current"},
	"btree_hits":          {"indexCounters", "btree", "hits"},
	"opcounters_insert":   {"opcounters", "insert"},
	"opcounters_query":    {"opcounters", "query"},
	"opcounters_update":   {"opcounters", "update"},
	"opcounters_delete":   {"opcounters", "delete"},
	"opcounters_getmore":  {"opcounters", "getmore"},
	"opcounters_command":  {"opcounters", "command"},
}

var metricPlace24 = map[string][]string{
	"duration_ms":         {"backgroundFlushing", "total_ms"},
	"connections_current": {"connections", "current"},
	"btree_hits":          {"indexCounters", "hits"},
	"opcounters_insert":   {"opcounters", "insert"},
	"opcounters_query":    {"opcounters", "query"},
	"opcounters_update":   {"opcounters", "update"},
	"opcounters_delete":   {"opcounters", "delete"},
	"opcounters_getmore":  {"opcounters", "getmore"},
	"opcounters_command":  {"opcounters", "command"},
}

// indexCounters is removed from mongodb 3.0.
// ref. http://stackoverflow.com/questions/29428793/where-is-the-indexcounter-in-db-serverstatus-on-mongodb-3-0
var metricPlace30 = map[string][]string{
	"duration_ms":         {"backgroundFlushing", "total_ms"},
	"connections_current": {"connections", "current"},
	"opcounters_insert":   {"opcounters", "insert"},
	"opcounters_query":    {"opcounters", "query"},
	"opcounters_update":   {"opcounters", "update"},
	"opcounters_delete":   {"opcounters", "delete"},
	"opcounters_getmore":  {"opcounters", "getmore"},
	"opcounters_command":  {"opcounters", "command"},
}

// backgroundFlushing information only appears for instances that use the MMAPv1 storage engine.
// and the MMAPv1 is no longer the default storage engine in MongoDB 3.2
// ref. https://docs.mongodb.org/manual/reference/command/serverStatus/#server-status-backgroundflushing

//Adapt to version 3.2 or higher.
//Check in version 3.6.

var metricPlace32 = map[string][]string{
	"connections_current": {"connections", "current"},
	"opcounters_insert":   {"opcounters", "insert"},
	"opcounters_query":    {"opcounters", "query"},
	"opcounters_update":   {"opcounters", "update"},
	"opcounters_delete":   {"opcounters", "delete"},
	"opcounters_getmore":  {"opcounters", "getmore"},
	"opcounters_command":  {"opcounters", "command"},
}

func getFloatValue(s map[string]interface{}, keys []string) (float64, error) {
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

// MongoDBPlugin mackerel plugin for mongo
type MongoDBPlugin struct {
	URL       string
	Username  string
	Password  string
	Source    string
	KeyPrefix string
	Verbose   bool
}

func (m MongoDBPlugin) fetchStatus() (bson.M, error) {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{m.URL},
		Username: m.Username,
		Password: m.Password,
		Source:   m.Source,
		Direct:   true,
		Timeout:  10 * time.Second,
	}
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		return nil, err
	}

	defer session.Close()
	session.SetMode(mgo.Eventual, true)
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

// FetchMetrics interface for mackerelplugin
func (m MongoDBPlugin) FetchMetrics() (map[string]interface{}, error) {
	serverStatus, err := m.fetchStatus()
	if err != nil {
		return nil, err
	}
	return m.parseStatus(serverStatus)
}

func (m MongoDBPlugin) getVersion(serverStatus bson.M) string {
	if reflect.TypeOf(serverStatus["version"]).String() == "string" {
		version := serverStatus["version"].(string)
		return version
	}
	return ""
}

func (m MongoDBPlugin) parseStatus(serverStatus bson.M) (map[string]interface{}, error) {
	stat := make(map[string]interface{})
	metricPlace := &metricPlace22

	cv, err := version.NewVersion(m.getVersion(serverStatus))
	if err != nil {
		return stat, err
	}

	if v, _ := version.NewVersion("3.2"); cv.Equal(v) || cv.GreaterThan(v) {
		metricPlace = &metricPlace32
	} else if v, _ := version.NewVersion("3.0"); cv.Equal(v) || cv.GreaterThan(v) {
		metricPlace = &metricPlace30
	} else if v, _ := version.NewVersion("2.4"); cv.Equal(v) || cv.GreaterThan(v) {
		metricPlace = &metricPlace24
	}

	for k, v := range *metricPlace {
		val, err := getFloatValue(serverStatus, v)
		if err != nil {
			logger.Warningf("Cannot fetch metric %s: %s", v, err)
		}

		stat[k] = val
	}

	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (m MongoDBPlugin) GraphDefinition() map[string]mp.Graphs {
	serverStatus, err := m.fetchStatus()
	if err != nil {
		return m.defaultGraphdef()
	}

	cv, err := version.NewVersion(m.getVersion(serverStatus))
	if err != nil {
		return m.defaultGraphdef()
	}
	if v, _ := version.NewVersion("3.2"); cv.Equal(v) || cv.GreaterThan(v) {
		return m.graphdef32()
	} else if v, _ := version.NewVersion("3.0"); cv.Equal(v) || cv.GreaterThan(v) {
		return m.graphdef30()
	}

	return m.defaultGraphdef()
}

const defaultPrefix = "mongodb"

// MetricKeyPrefix returns the metrics key prefix
func (m MongoDBPlugin) MetricKeyPrefix() string {
	if m.KeyPrefix == "" {
		m.KeyPrefix = defaultPrefix
	}
	return m.KeyPrefix
}

func (m MongoDBPlugin) LabelPrefix() string {
	return cases.Title(language.Und, cases.NoLower).String(strings.Replace(m.MetricKeyPrefix(), defaultPrefix, "MongoDB", -1))
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "27017", "Port")
	optUser := flag.String("username", "", "Username")
	optPass := flag.String("password", os.Getenv("MONGODB_PASSWORD"), "Password")
	optSource := flag.String("source", "", "authenticationDatabase")
	optVerbose := flag.Bool("v", false, "Verbose mode")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optKeyPrefix := flag.String("metric-key-prefix", "", "Metric key prefix")
	flag.Parse()

	var mongodb MongoDBPlugin
	mongodb.Verbose = *optVerbose
	mongodb.URL = net.JoinHostPort(*optHost, *optPort)
	mongodb.Username = *optUser
	mongodb.Password = *optPass
	mongodb.Source = *optSource
	mongodb.KeyPrefix = *optKeyPrefix

	helper := mp.NewMackerelPlugin(mongodb)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.SetTempfileByBasename(fmt.Sprintf("mackerel-plugin-mongodb-%s-%s", *optHost, *optPort))
	}

	helper.Run()
}
