//go:build windows

package mpmssql

//go:generate wmi2struct -n -p mpmssql -o wmi.go Win32_PerfRawData_MSSQLSQLEXPRESS_MSSQLSQLEXPRESSBufferManager Win32_PerfRawData_MSSQLSQLEXPRESS_MSSQLSQLEXPRESSSQLStatistics Win32_PerfRawData_MSSQLSQLEXPRESS_MSSQLSQLEXPRESSGeneralStatistics Win32_PerfRawData_MSSQLSQLEXPRESS_MSSQLSQLEXPRESSAccessMethods

import (
	"flag"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/StackExchange/wmi"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

// MSSQLPlugin store the name of servers
type MSSQLPlugin struct {
	prefix   string
	instance string
}

// MetricKeyPrefix retruns the metrics key prefix
func (m MSSQLPlugin) MetricKeyPrefix() string {
	if m.prefix == "" {
		m.prefix = "mssql"
	}
	return m.prefix
}

func (m MSSQLPlugin) name(n string) string {
	if m.instance == "SQLSERVER" {
		return "Win32_PerfRawData_MSSQLSERVER_SQLServer" + n
	} else {
		return "Win32_PerfRawData_MSSQLSQLEXPRESS_MSSQLSQLEXPRESS" + n
	}
}

// FetchMetrics interface for mackerelplugin
func (m MSSQLPlugin) FetchMetrics() (map[string]float64, error) {
	var bufferManager []Win32_PerfRawData_MSSQLSQLEXPRESS_MSSQLSQLEXPRESSBufferManager
	var sqlStatistics []Win32_PerfRawData_MSSQLSQLEXPRESS_MSSQLSQLEXPRESSSQLStatistics
	var generalStatistics []Win32_PerfRawData_MSSQLSQLEXPRESS_MSSQLSQLEXPRESSGeneralStatistics
	var accessMethods []Win32_PerfRawData_MSSQLSQLEXPRESS_MSSQLSQLEXPRESSAccessMethods

	var err error
	err = wmi.Query("select * from "+m.name("BufferManager"), &bufferManager)
	if err != nil {
		return nil, err
	}
	err = wmi.Query("select * from "+m.name("SQLStatistics"), &sqlStatistics)
	if err != nil {
		return nil, err
	}
	err = wmi.Query("select * from "+m.name("GeneralStatistics"), &generalStatistics)
	if err != nil {
		return nil, err
	}
	err = wmi.Query("select * from "+m.name("AccessMethods"), &accessMethods)
	if err != nil {
		return nil, err
	}

	stat := make(map[string]float64)
	stat["buffer_cache_hit_ratio"] = float64(bufferManager[0].Buffercachehitratio)
	stat["buffer_page_life_expectancy"] = float64(bufferManager[0].Pagelifeexpectancy)
	stat["buffer_checkpoint_pages"] = float64(bufferManager[0].CheckpointpagesPersec)
	stat["stats_batch_requests"] = float64(sqlStatistics[0].BatchRequestsPersec)
	stat["stats_sql_compilations"] = float64(sqlStatistics[0].SQLCompilationsPersec)
	stat["stats_sql_recompilations"] = float64(sqlStatistics[0].SQLReCompilationsPersec)
	stat["stats_connections"] = float64(generalStatistics[0].UserConnections)
	stat["stats_lock_waits"] = float64(generalStatistics[0].SQLTraceIOProviderLockWaits)
	stat["stats_procs_blocked"] = float64(generalStatistics[0].Processesblocked)
	stat["access_page_splits"] = float64(accessMethods[0].PageSplitsPersec)
	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (m MSSQLPlugin) GraphDefinition() map[string](mp.Graphs) {
	labelPrefix := cases.Title(language.Und, cases.NoLower).String(strings.Replace(m.MetricKeyPrefix(), "mssql", "MSSQL", -1))
	return map[string](mp.Graphs){
		"buffer": mp.Graphs{
			Label: labelPrefix + " Buffer",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{
					Name:    "buffer_cache_hit_ratio",
					Label:   "Cache Hit Ratio",
					Diff:    false,
					Stacked: true,
				},
				{
					Name:    "buffer_page_life_expectancy",
					Label:   "Page Life Expectancy",
					Diff:    true,
					Stacked: true,
				},
				{
					Name:    "buffer_checkpoint_pages",
					Label:   "Checkpoint Pages",
					Diff:    false,
					Stacked: true,
				},
			},
		},
		"stats": mp.Graphs{
			Label: labelPrefix + " Stats",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{
					Name:    "stats_batch_requests",
					Label:   "Batch Requests",
					Diff:    true,
					Stacked: true,
				},
				{
					Name:    "stats_sql_compilations",
					Label:   "SQL Compilations",
					Diff:    true,
					Stacked: true,
				},
				{
					Name:    "stats_sql_recompilations",
					Label:   "SQL Re-Compilations",
					Diff:    true,
					Stacked: true,
				},
				{
					Name:    "stats_connections",
					Label:   "Connections",
					Diff:    false,
					Stacked: true,
				},
				{
					Name:    "stats_lock_waits",
					Label:   "Lock Waits",
					Diff:    false,
					Stacked: true,
				},
				{
					Name:    "stats_procs_blocked",
					Label:   "Procs Blocked",
					Diff:    false,
					Stacked: true,
				},
			},
		},
		"access": mp.Graphs{
			Label: labelPrefix + " Access",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{
					Name:    "access_page_splits",
					Label:   "Page Splits",
					Diff:    true,
					Stacked: true,
				},
			},
		},
	}
}

// Do the plugin
func Do() {
	optPrefix := flag.String("metric-key-prefix", "mssql", "Metric key prefix")
	optInstance := flag.String("instance", "SQLSERVER", "Instance name of MSSQL(SQLSERVER or SQLEXPRESS)")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	instance := strings.ToUpper(*optInstance)
	if instance != "SQLSERVER" && instance != "SQLEXPRESS" {
		flag.Usage()
		os.Exit(2)
	}
	plugin := MSSQLPlugin{
		prefix:   *optPrefix,
		instance: instance,
	}
	helper := mp.NewMackerelPlugin(plugin)
	helper.Tempfile = *optTempfile
	helper.Run()
}
