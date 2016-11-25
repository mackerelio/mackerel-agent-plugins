package mppostgres

import (
	"flag"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	// PostgreSQL Driver
	_ "github.com/lib/pq"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/mackerel-agent/logging"
)

var logger = logging.GetLogger("metrics.plugin.postgres")

var graphdef = map[string]mp.Graphs{
	"postgres.connections": {
		Label: "Postgres Connections",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "active", Label: "Active", Diff: false, Stacked: true},
			{Name: "waiting", Label: "Waiting", Diff: false, Stacked: true},
		},
	},
	"postgres.commits": {
		Label: "Postgres Commits",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "xact_commit", Label: "Xact Commit", Diff: true, Stacked: false},
			{Name: "xact_rollback", Label: "Xact Rollback", Diff: true, Stacked: false},
		},
	},
	"postgres.blocks": {
		Label: "Postgres Blocks",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "blks_read", Label: "Blocks Read", Diff: true, Stacked: false},
			{Name: "blks_hit", Label: "Blocks Hit", Diff: true, Stacked: false},
		},
	},
	"postgres.rows": {
		Label: "Postgres Rows",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "tup_returned", Label: "Returned Rows", Diff: true, Stacked: false},
			{Name: "tup_fetched", Label: "Fetched Rows", Diff: true, Stacked: true},
			{Name: "tup_inserted", Label: "Inserted Rows", Diff: true, Stacked: true},
			{Name: "tup_updated", Label: "Updated Rows", Diff: true, Stacked: true},
			{Name: "tup_deleted", Label: "Deleted Rows", Diff: true, Stacked: true},
		},
	},
	"postgres.size": {
		Label: "Postgres Data Size",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "total_size", Label: "Total Size", Diff: false, Stacked: false},
		},
	},
	"postgres.deadlocks": {
		Label: "Postgres Dead Locks",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "deadlocks", Label: "Deadlocks", Diff: true, Stacked: false},
		},
	},
	"postgres.iotime": {
		Label: "Postgres Block I/O time",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "blk_read_time", Label: "Block Read Time (ms)", Diff: true, Stacked: false},
			{Name: "blk_write_time", Label: "Block Write Time (ms)", Diff: true, Stacked: false},
		},
	},
	"postgres.tempfile": {
		Label: "Postgres Temporary file",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "temp_bytes", Label: "Temporary file size (byte)", Diff: true, Stacked: false},
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
	db = db.Unsafe()
	rows, err := db.Queryx(`SELECT * FROM pg_stat_database`)
	if err != nil {
		logger.Errorf("Failed to select pg_stat_database. %s", err)
		return nil, err
	}

	type pgStat struct {
		XactCommit   uint64   `db:"xact_commit"`
		XactRollback uint64   `db:"xact_rollback"`
		BlksRead     uint64   `db:"blks_read"`
		BlksHit      uint64   `db:"blks_hit"`
		BlkReadTime  *float64 `db:"blk_read_time"`
		BlkWriteTime *float64 `db:"blk_write_time"`
		TupReturned  uint64   `db:"tup_returned"`
		TupFetched   uint64   `db:"tup_fetched"`
		TupInserted  uint64   `db:"tup_inserted"`
		TupUpdated   uint64   `db:"tup_updated"`
		TupDeleted   uint64   `db:"tup_deleted"`
		Deadlocks    *uint64  `db:"deadlocks"`
		TempBytes    *uint64  `db:"temp_bytes"`
	}

	totalStat := pgStat{}
	for rows.Next() {
		p := pgStat{}
		if err := rows.StructScan(&p); err != nil {
			logger.Warningf("Failed to scan. %s", err)
			continue
		}
		totalStat.XactCommit += p.XactCommit
		totalStat.XactRollback += p.XactRollback
		totalStat.BlksRead += p.BlksRead
		totalStat.BlksHit += p.BlksHit
		if p.BlkReadTime != nil {
			if totalStat.BlkReadTime == nil {
				totalStat.BlkReadTime = p.BlkReadTime
			} else {
				*totalStat.BlkReadTime += *p.BlkReadTime
			}
		}
		if p.BlkWriteTime != nil {
			if totalStat.BlkWriteTime == nil {
				totalStat.BlkWriteTime = p.BlkWriteTime
			} else {
				*totalStat.BlkWriteTime += *p.BlkWriteTime
			}
		}
		totalStat.TupReturned += p.TupReturned
		totalStat.TupFetched += p.TupFetched
		totalStat.TupInserted += p.TupInserted
		totalStat.TupUpdated += p.TupUpdated
		totalStat.TupDeleted += p.TupDeleted
		if p.Deadlocks != nil {
			if totalStat.Deadlocks == nil {
				totalStat.Deadlocks = p.Deadlocks
			} else {
				*totalStat.Deadlocks += *p.Deadlocks
			}
		}
		if p.TempBytes != nil {
			if totalStat.TempBytes == nil {
				totalStat.TempBytes = p.TempBytes
			} else {
				*totalStat.TempBytes += *p.TempBytes
			}
		}
	}
	stat := make(map[string]interface{})
	stat["xact_commit"] = totalStat.XactCommit
	stat["xact_rollback"] = totalStat.XactRollback
	stat["blks_read"] = totalStat.BlksRead
	stat["blks_hit"] = totalStat.BlksHit
	if totalStat.BlkReadTime != nil {
		stat["blk_read_time"] = *totalStat.BlkReadTime
	}
	if totalStat.BlkWriteTime != nil {
		stat["blk_write_time"] = *totalStat.BlkWriteTime
	}
	stat["tup_returned"] = totalStat.TupReturned
	stat["tup_fetched"] = totalStat.TupFetched
	stat["tup_inserted"] = totalStat.TupInserted
	stat["tup_updated"] = totalStat.TupUpdated
	stat["tup_deleted"] = totalStat.TupDeleted
	if totalStat.Deadlocks != nil {
		stat["deadlocks"] = *totalStat.Deadlocks
	}
	if totalStat.TempBytes != nil {
		stat["temp_bytes"] = *totalStat.TempBytes
	}
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
func (p PostgresPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// Do the plugin
func Do() {
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
