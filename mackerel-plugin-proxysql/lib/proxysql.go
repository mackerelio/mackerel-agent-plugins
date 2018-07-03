package mpproxysql

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"database/sql"

	mp "github.com/mackerelio/go-mackerel-plugin"
	// Go driver for database/sql package
	_ "github.com/ziutek/mymysql/godrv"
)

func (m *ProxySQLPlugin) proyxsqlGraphDef() map[string]mp.Graphs {
	labelPrefix := strings.Title(strings.Replace(m.MetricKeyPrefix(), "proxysql", "ProxySQL", -1))
	return map[string]mp.Graphs{
		"uptime": {
			Label: labelPrefix + "Uptime",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "ProxySQL_Uptime", Label: "Seconds"},
			},
		},
	}
}

// ProxySQLPlugin mackerel plugin for ProxySQL
type ProxySQLPlugin struct {
	URI      string
	Tempfile string
	prefix   string
}

// MetricKeyPrefix interface for PluginWithPrefix
func (m *ProxySQLPlugin) MetricKeyPrefix() string {
	if m.prefix == "" {
		m.prefix = "proxysql"
	}
	return m.prefix
}

func (m *ProxySQLPlugin) parseStatsMySQLGlobal(rows *sql.Rows, stat map[string]float64) error {
	for rows.Next() {
		var name, value string
		err := rows.Scan(&name, &value)
		if err != nil {
			return fmt.Errorf("parseStatsMySQLGlobal (stats_mysql_global): %s", err)
		}

		stat[name], err = strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("%s", err)
		}
	}
	return rows.Err()
}

func (m *ProxySQLPlugin) fetchStatsMySQLGlobal(db *sql.DB, stat map[string]float64) error {
	rows, err := db.Query("SELECT Variable_Name, Variable_Value FROM stats_mysql_global")
	if err != nil {
		return fmt.Errorf("FetchMetrics (stats_mysql_global): %s", err)
	}

	return m.parseStatsMySQLGlobal(rows, stat)
}

// FetchMetrics interface for mackerelplugin
func (m *ProxySQLPlugin) FetchMetrics() (map[string]float64, error) {
	db, err := sql.Open("mymysql", m.URI)
	if err != nil {
		return nil, fmt.Errorf("FetchMetrics (DB Connect): %s", err)
	}
	defer db.Close()

	stat := make(map[string]float64)

	err = m.fetchStatsMySQLGlobal(db, stat)
	if err != nil {
		return nil, err
	}

	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (m *ProxySQLPlugin) GraphDefinition() map[string]mp.Graphs {
	return m.proyxsqlGraphDef()
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "0.0.0.0", "Hostname")
	optPort := flag.String("port", "6032", "Port")
	optSocket := flag.String("socket", "", "Path to unix socket")
	optUser := flag.String("username", "radminuser", "Username")
	optPass := flag.String("password", os.Getenv("PROXYSQL_PASSWORD"), "Password")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optMetricKeyPrefix := flag.String("metric-key-prefix", "proxysql", "metric key prefix")
	flag.Parse()

	var proxysql ProxySQLPlugin

	if *optSocket != "" {
		proxysql.URI = fmt.Sprintf("unix:%s*stats/%s/%s", *optSocket, *optUser, *optPass)
	} else {
		proxysql.URI = fmt.Sprintf("tcp:%s:%s*stats/%s/%s", *optHost, *optPort, *optUser, *optPass)
	}

	proxysql.prefix = *optMetricKeyPrefix
	helper := mp.NewMackerelPlugin(&proxysql)
	helper.Tempfile = *optTempfile

	helper.Run()
}
