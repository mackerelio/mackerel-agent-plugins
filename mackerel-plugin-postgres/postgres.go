package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
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
	Host     string
	Port     string
	Username string
	Password string
	SSLmode  string
	Timeout  int
	Tempfile string
	Option   string
}

func fetchStatDatabase(db *sqlx.DB) (map[string]interface{}, error) {
	rows, err := db.Query(`
		select xact_commit, xact_rollback, blks_read, blks_hit, blk_read_time, blk_write_time,
		tup_returned, tup_fetched, tup_inserted, tup_updated, tup_deleted, deadlocks, temp_bytes
		from pg_stat_database
	`)
	if err != nil {
		logger.Errorf("Failed to select pg_stat_database. %s", err)
		return nil, err
	}

	type pgStat struct {
		XactCommit   *float64 `db:"xact_commit"`
		XactRollback *float64 `db:"xact_rollback"`
		BlksRead     *float64 `db:"blks_read"`
		BlksHit      *float64 `db:"blks_hit"`
		BlkReadTime  *float64 `db:"blk_read_time"`
		BlkWriteTime *float64 `db:"blk_write_time"`
		TupReturned  *float64 `db:"tup_returned"`
		TupFetched   *float64 `db:"tup_fetched"`
		TupInserted  *float64 `db:"tup_inserted"`
		TupUpdated   *float64 `db:"tup_updated"`
		TupDeleted   *float64 `db:"tup_deleted"`
		Deadlocks    *float64 `db:"deadlocks"`
		TempBytes    *float64 `db:"temp_bytes"`
	}

	var totalXactCommit, totalXactRollback, totalBlksRead, totalBlksHit, totalBlkReadTime, totalBlkWriteTime, totalTupReturned, totalTupFetched, totalTupInserted, totalTupUpdated, totalTupDeleted, totalDeadlocks, totalTempBytes float64

	for rows.Next() {
		var xactCommit, xactRollback, blksRead, blksHit, blkReadTime, blkWriteTime, tupReturned, tupFetched, tupInserted, tupUpdated, tupDeleted, deadlocks, tempBytes float64

		if err := rows.Scan(&xactCommit, &xactRollback, &blksRead, &blksHit, &blkReadTime, &blkWriteTime, &tupReturned, &tupFetched, &tupInserted, &tupUpdated, &tupDeleted, &deadlocks, &tempBytes); err != nil {
			logger.Warningf("Failed to scan. %s", err)
			continue
		}

		totalXactCommit += xactCommit
		totalXactRollback += xactRollback
		totalBlksRead += blksRead
		totalBlksHit += blksHit
		totalBlkReadTime += blkReadTime
		totalBlkWriteTime += blkWriteTime
		totalTupReturned += tupReturned
		totalTupFetched += tupFetched
		totalTupInserted += tupInserted
		totalTupUpdated += tupUpdated
		totalTupDeleted += tupDeleted
		totalDeadlocks += deadlocks
		totalTempBytes += tempBytes
	}
	stat := make(map[string]interface{})
	stat["xact_commit"] = totalXactCommit
	stat["xact_rollback"] = totalXactRollback
	stat["blks_read"] = totalBlksRead
	stat["blks_hit"] = totalBlksHit
	stat["blk_read_time"] = totalBlkReadTime
	stat["blk_write_time"] = totalBlkWriteTime
	stat["tup_returned"] = totalTupReturned
	stat["tup_fetched"] = totalTupFetched
	stat["tup_inserted"] = totalTupInserted
	stat["tup_updated"] = totalTupUpdated
	stat["tup_deleted"] = totalTupDeleted
	stat["deadlocks"] = totalDeadlocks
	stat["temp_bytes"] = totalTempBytes

	return stat, nil
}

func fetchConnections(db *sqlx.DB) (map[string]interface{}, error) {
	rows, err := db.Query(`
		select count(*), waiting from pg_stat_activity group by waiting
	`)
	if err != nil {
		logger.Errorf("Failed to select pg_stat_activity. %s", err)
		return nil, err
	}

	var totalActive, totalWaiting float64
	for rows.Next() {
		var count float64
		var waiting string
		if err := rows.Scan(&count, &waiting); err != nil {
			logger.Warningf("Failed to scan %s", err)
			continue
		}
		if waiting != "" {
			totalActive += count
		} else {
			totalWaiting += count
		}
	}

	return map[string]interface{}{
		"active":  totalActive,
		"waiting": totalWaiting,
	}, nil
}

func fetchDatabaseSize(db *sqlx.DB) (map[string]interface{}, error) {
	rows, err := db.Query("select sum(pg_database_size(datname)) as dbsize from pg_database")
	if err != nil {
		logger.Errorf("Failed to select pg_database_size. %s", err)
		return nil, err
	}

	var totalSize float64
	for rows.Next() {
		var dbsize float64
		if err := rows.Scan(&dbsize); err != nil {
			logger.Warningf("Failed to scan %s", err)
			continue
		}
		totalSize += dbsize
	}

	return map[string]interface{}{
		"total_size": totalSize,
	}, nil
}

func mergeStat(dst, src map[string]interface{}) {
	for k, v := range src {
		dst[k] = v
	}
}

// FetchMetrics interface for mackerelplugin
func (p PostgresPlugin) FetchMetrics() (map[string]interface{}, error) {

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=%s connect_timeout=%d %s", p.Username, p.Password, p.Host, p.Port, p.SSLmode, p.Timeout, p.Option))
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

	stat := make(map[string]interface{})
	mergeStat(stat, statStatDatabase)
	mergeStat(stat, statConnections)
	mergeStat(stat, statDatabaseSize)

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
	option := ""
	if *optDatabase != "" {
		option = fmt.Sprintf("dbname=%s", *optDatabase)
	}

	var postgres PostgresPlugin
	postgres.Host = *optHost
	postgres.Port = *optPort
	postgres.Username = *optUser
	postgres.Password = *optPass
	postgres.SSLmode = *optSSLmode
	postgres.Timeout = *optConnectTimeout
	postgres.Option = option

	helper := mp.NewMackerelPlugin(postgres)

	helper.Tempfile = *optTempfile
	helper.Run()
}
