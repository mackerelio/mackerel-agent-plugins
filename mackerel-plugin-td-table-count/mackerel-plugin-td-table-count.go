package main

import (
	"flag"
	"os"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/mackerelio/mackerel-agent/logging"
	td "github.com/mattn/go-treasuredata"
)

var logger = logging.GetLogger("metrics.plugin.td-table")

var tableNames = []string{}

type TDTablePlugin struct {
	ApiKey           string
	Database         string
	IgnoreTableNames []string
	Tempfile         string
}

func (m TDTablePlugin) FetchMetrics() (map[string]float64, error) {
	stat := make(map[string]float64)
	cli := td.NewClient(m.ApiKey)

	tables, err := cli.TableList(m.Database)
	if err != nil {
		return nil, err
	}

	for _, table := range tables {
		ignore := false
		for _, ignoreTableName := range m.IgnoreTableNames {
			if table.Name == ignoreTableName {
				ignore = true
			}
		}

		if ignore == false {
			stat[table.Name] = float64(table.Count)
			tableNames = append(tableNames, table.Name)
		}
	}

	return stat, nil
}

func (m TDTablePlugin) GraphDefinition() map[string](mp.Graphs) {
	metrics := []mp.Metrics{}
	for _, tableName := range tableNames {
		metrics = append(metrics, mp.Metrics{
			Name: tableName,
			Diff: false,
		})
	}

	return map[string](mp.Graphs){
		"td.count": mp.Graphs{
			Label:   "TD Number of rows",
			Unit:    "float",
			Metrics: metrics,
		},
	}
}

func main() {
	optApiKey := flag.String("api-key", "", "API Key")
	optDatabase := flag.String("database", "", "Database name")
	optIgnoreTableNames := flag.String("ignore-table", "", "Ignore Table name (Can be Comma-Separated)")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var plugin TDTablePlugin
	plugin.ApiKey = *optApiKey
	plugin.Database = *optDatabase

	ignoreTableNames := []string{}
	if *optIgnoreTableNames != "" {
		ignoreTableNames = strings.Split(*optIgnoreTableNames, ",")
	}
	plugin.IgnoreTableNames = ignoreTableNames

	helper := mp.NewMackerelPlugin(plugin)

	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = "/tmp/mackerel-plugin-td-table"
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
