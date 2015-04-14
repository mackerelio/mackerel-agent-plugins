package main

import (
	"flag"
	"fmt"
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

func GetTables(m TDTablePlugin) ([]td.Table, error) {
	cli := td.NewClient(m.ApiKey)

	tables, err := cli.TableList(m.Database)
	if err != nil {
		return nil, err
	}

	filteredTables := []td.Table{}

	for _, table := range tables {
		ignore := false
		for _, ignoreTableName := range m.IgnoreTableNames {
			if table.Name == ignoreTableName {
				ignore = true
			}
		}

		if ignore == false {
			filteredTables = append(filteredTables, table)
		}
	}

	return filteredTables, nil
}

func (m TDTablePlugin) FetchMetrics() (map[string]float64, error) {
	stat := make(map[string]float64)

	tables, _ := GetTables(m)
	for _, table := range tables {
		stat[table.Name] = float64(table.Count)
	}

	return stat, nil
}

func (m TDTablePlugin) GraphDefinition() map[string](mp.Graphs) {
	tables, _ := GetTables(m)

	var metrics []mp.Metrics
	for _, table := range tables {
		metrics = append(metrics, mp.Metrics{
			Name:    table.Name,
			Label:   table.Name,
			Diff:    false,
			Stacked: true,
		})
	}

	graph := mp.Graphs{
		Label:   fmt.Sprintf("TD %s Database Number of rows", m.Database),
		Unit:    "integer",
		Metrics: metrics,
	}

	graphdef := map[string](mp.Graphs){
		fmt.Sprintf("td-table.%s", m.Database): graph,
	}

	return graphdef
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
