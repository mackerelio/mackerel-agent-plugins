package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/mackerel-agent/logging"
)

var logger = logging.GetLogger("metrics.plugin.postgres")

var graphdef = map[string](mp.Graphs){
	"postgres.connections": mp.Graphs{
		Label: "Postgres Connections",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "active", Label: "Active", Diff: false, Stacked: true},
			mp.Metrics{Name: "waiting", Label: "Waiting", Diff: false, Stacked: true},
		},
	},
	"postgres.commits": mp.Graphs{
		Label: "Postgres Commits",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "xact_commit", Label: "Xact Commit", Diff: true, Stacked: false},
			mp.Metrics{Name: "xact_rollback", Label: "Xact Rollback", Diff: true, Stacked: false},
		},
	},
	"postgres.blocks": mp.Graphs{
		Label: "Postgres Blocks",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "blks_read", Label: "Blocks Read", Diff: true, Stacked: false},
			mp.Metrics{Name: "blks_hit", Label: "Blocks Hit", Diff: true, Stacked: false},
		},
	},
	"postgres.rows": mp.Graphs{
		Label: "Postgres Rows",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "tup_returned", Label: "Returned Rows", Diff: true, Stacked: false},
			mp.Metrics{Name: "tup_fetched", Label: "Fetched Rows", Diff: true, Stacked: true},
			mp.Metrics{Name: "tup_inserted", Label: "Inserted Rows", Diff: true, Stacked: true},
			mp.Metrics{Name: "tup_updated", Label: "Updated Rows", Diff: true, Stacked: true},
			mp.Metrics{Name: "tup_deleted", Label: "Deleted Rows", Diff: true, Stacked: true},
		},
	},
	"postgres.size": mp.Graphs{
		Label: "Postgres Data Size",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "total_size", Label: "Total Size", Diff: false, Stacked: false},
		},
	},
	"postgres.relpages.#": mp.Graphs{
		Label: "Postgres Relpages",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "*", Label: "%1", Diff: false, Stacked: false},
		},
	},
	"postgres.deadlocks": mp.Graphs{
		Label: "Postgres Dead Locks",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "deadlocks", Label: "Deadlocks", Diff: true, Stacked: false},
		},
	},
	"postgres.iotime": mp.Graphs{
		Label: "Postgres Block I/O time",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "blk_read_time", Label: "Block Read Time (ms)", Diff: true, Stacked: false},
			mp.Metrics{Name: "blk_write_time", Label: "Block Write Time (ms)", Diff: true, Stacked: false},
		},
	},
	"postgres.tempfile": mp.Graphs{
		Label: "Postgres Temporary file",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "temp_bytes", Label: "Temporary file size (byte)", Diff: true, Stacked: false},
		},
	},
}

// PostgresPlugin mackerel plugin for PostgreSQL
type PostgresPlugin struct {
	Host          string
	Port          string
	Username      string
	Password      string
	SSLmode       string
	Timeout       int
	Tempfile      string
	Database      string
	RelpagesLimit int
}

func fetchStatDatabase(db *sql.DB) (map[string]float64, error) {
	rows, err := db.Query(`
		select xact_commit, xact_rollback, blks_read, blks_hit, blk_read_time, blk_write_time,
		tup_returned, tup_fetched, tup_inserted, tup_updated, tup_deleted, deadlocks, temp_bytes
		from pg_stat_database
	`)
	if err != nil {
		logger.Errorf("Failed to select pg_stat_database. %s", err)
		return nil, err
	}

	stat := make(map[string]float64)

	for rows.Next() {
		var xactCommit, xactRollback, blksRead, blksHit, blkReadTime, blkWriteTime, tupReturned, tupFetched, tupInserted, tupUpdated, tupDeleted, deadlocks, tempBytes int

		if err := rows.Scan(&xactCommit, &xactRollback, &blksRead, &blksHit, &blkReadTime, &blkWriteTime, &tupReturned, &tupFetched, &tupInserted, &tupUpdated, &tupDeleted, &deadlocks, &tempBytes); err != nil {
			logger.Warningf("Failed to scan. %s", err)
			continue
		}

		stat["xact_commit"] += float64(xactCommit)
		stat["xact_rollback"] += float64(xactRollback)
		stat["blks_read"] += float64(blksRead)
		stat["blks_hit"] += float64(blksHit)
		stat["blk_read_time"] += float64(blksHit)
		stat["blk_write_time"] += float64(blksHit)
		stat["tup_returned"] += float64(tupReturned)
		stat["tup_fetched"] += float64(tupFetched)
		stat["tup_inserted"] += float64(tupInserted)
		stat["tup_updated"] += float64(tupUpdated)
		stat["tup_deleted"] += float64(tupDeleted)
		stat["deadlocks"] += float64(deadlocks)
		stat["temp_bytes"] += float64(tempBytes)
	}

	return stat, nil
}

func fetchConnections(db *sql.DB) (map[string]float64, error) {
	rows, err := db.Query(`
		select count(*), waiting from pg_stat_activity group by waiting
	`)
	if err != nil {
		logger.Errorf("Failed to select pg_stat_activity. %s", err)
		return nil, err
	}

	stat := make(map[string]float64)

	for rows.Next() {
		var count int
		var waiting string
		if err := rows.Scan(&count, &waiting); err != nil {
			logger.Warningf("Failed to scan %s", err)
			continue
		}
		if waiting != "" {
			stat["active"] += float64(count)
		} else {
			stat["waiting"] += float64(count)
		}
	}

	return stat, nil
}

func fetchDatabaseSize(db *sql.DB) (map[string]float64, error) {
	rows, err := db.Query("select sum(pg_database_size(datname)) as dbsize from pg_database")
	if err != nil {
		logger.Errorf("Failed to select pg_database_size. %s", err)
		return nil, err
	}

	stat := make(map[string]float64)

	for rows.Next() {
		var dbsize int
		if err := rows.Scan(&dbsize); err != nil {
			logger.Warningf("Failed to scan %s", err)
			continue
		}
		stat["total_size"] += float64(dbsize)
	}

	return stat, nil
}

func (p PostgresPlugin) fetchRelpages(db *sql.DB) (map[string]float64, error) {
	if p.RelpagesLimit <= 0 {
		return nil, nil
	}
	_, err := db.Query("analyze")
	if err != nil {
		logger.Warningf("Failed to ANALYZE", err)
		return nil, err
	}
	rows, err := db.Query("select relname, relpages from pg_class order by relpages desc limit $1", p.RelpagesLimit)
	if err != nil {
		logger.Errorf("Failed to SELECT pg_class.relpages. %s", err)
		return nil, err
	}

	stat := make(map[string]float64)

	for rows.Next() {
		var relname string
		var relpages int64
		if err := rows.Scan(&relname, &relpages); err != nil {
			logger.Warningf("Failed to scan %s", err)
			continue
		}
		stat["postgres.relpages."+p.Database+"."+relname] = float64(relpages)
	}

	return stat, nil
}
func mergeStat(dst map[string]interface{}, src map[string]float64) {
	for k, v := range src {
		dst[k] = v
	}
}

func (p PostgresPlugin) option() string {
	option := ""
	if p.Database != "" {
		option = fmt.Sprintf("dbname=%s", p.Database)
	}
	return option
}

// FetchMetrics interface for mackerelplugin
func (p PostgresPlugin) FetchMetrics() (map[string]interface{}, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=%s connect_timeout=%d %s", p.Username, p.Password, p.Host, p.Port, p.SSLmode, p.Timeout, p.option()))
	if err != nil {
		logger.Errorf("FetchMetrics: %s", err)
		return nil, err
	}
	defer db.Close()

	statStatDatabase, err := fetchStatDatabase(db)
	if err != nil {
		return nil, err
	}
	statConnections, err := fetchConnections(db)
	if err != nil {
		return nil, err
	}
	statDatabaseSize, err := fetchDatabaseSize(db)
	if err != nil {
		return nil, err
	}
	statRelpages, err := p.fetchRelpages(db)
	if err != nil {
		return nil, err
	}

	stat := make(map[string]interface{})
	mergeStat(stat, statStatDatabase)
	mergeStat(stat, statConnections)
	mergeStat(stat, statDatabaseSize)
	mergeStat(stat, statRelpages)

	return stat, err
}

// GraphDefinition interface for mackerelplugin
func (p PostgresPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optHost := flag.String("hostname", "localhost", "Hostname to login to")
	optPort := flag.String("port", "5432", "Database port")
	optUser := flag.String("user", "", "Postgres User")
	optDatabase := flag.String("database", "", "Database name")
	optPass := flag.String("password", "", "Postgres Password")
	optSSLmode := flag.String("sslmode", "disable", "Whether or not to use SSL")
	optConnectTimeout := flag.Int("connect_timeout", 5, "Maximum wait for connection, in seconds.")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optRelpagesLimit := flag.Int("relpages_limit", 0, "Outputs only top `relpages_limit` relations with the highest relpages.")
	flag.Parse()

	if *optUser == "" {
		logger.Warningf("user is required")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *optPass == "" {
		logger.Warningf("password is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var postgres PostgresPlugin
	postgres.Host = *optHost
	postgres.Port = *optPort
	postgres.Username = *optUser
	postgres.Password = *optPass
	postgres.SSLmode = *optSSLmode
	postgres.Timeout = *optConnectTimeout
	postgres.Database = *optDatabase
	postgres.RelpagesLimit = *optRelpagesLimit

	helper := mp.NewMackerelPlugin(postgres)

	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-postgres-%s-%s", *optHost, *optPort)
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
