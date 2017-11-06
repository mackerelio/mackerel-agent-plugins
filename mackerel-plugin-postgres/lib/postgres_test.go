package mppostgres

import (
	"testing"

	"github.com/erikstmartin/go-testdb"
	"github.com/jmoiron/sqlx"
)

func TestFetchStatDatabase(t *testing.T) {
	db, _ := sqlx.Connect("testdb", "")

	columns := []string{"xact_commit", "xact_rollback", "blks_read", "blks_hit", "blk_read_time", "blk_write_time",
		"tup_returned", "tup_fetched", "tup_inserted", "tup_updated", "tup_deleted", "deadlocks", "temp_bytes"}

	testdb.StubQuery(`SELECT * FROM pg_stat_database`, testdb.RowsFromCSVString(columns, `
	1,2,3,4,5,6,7,8,9,10,11,12,13
	10,20,30,40,50,60,70,80,90,100,110,120,130
	`))

	stat, err := fetchStatDatabase(db)

	expected := map[string]interface{}{
		"xact_commit":  uint64(11),
		"blks_hit":     uint64(44),
		"tup_returned": uint64(77),
	}

	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
	}
	if err = db.Close(); err != nil {
		t.Errorf("Error '%s' was not expected while closing the database", err)
	}
	if stat["xact_commit"] != expected["xact_commit"] {
		t.Error("should be 11")
	}
	if stat["blks_hit"] != expected["blks_hit"] {
		t.Error("should be 44")
	}
	if stat["tup_returned"] != expected["tup_returned"] {
		t.Error("should be 77")
	}
}

var fetchVersionTests = []struct {
	response string
	expected version
}{
	{
		`
		PostgreSQL 9.6.4 on x86_64-redhat-linux-gnu, compiled by gcc (GCC) 4.8.5 20150623 (Red Hat 4.8.5-11), 64-bit
		`,
		version{uint(9), uint(6), uint(4)},
	},
	{
		// Azure Database for PostgreSQL
		`
		PostgreSQL 9.6.5, compiled by Visual C++ build 1800, 64-bit
		`,
		version{uint(9), uint(6), uint(5)},
	},
	{
		`
		PostgreSQL 10.0 on x86_64-pc-linux-gnu, compiled by gcc (Debian 6.3.0-18) 6.3.0 20170516, 64-bit
		`,
		version{uint(10), uint(0), uint(0)},
	},
}

func TestFetchVersion(t *testing.T) {
	for _, tc := range fetchVersionTests {
		db, _ := sqlx.Connect("testdb", "")

		columns := []string{"version"}

		testdb.StubQuery(`SELECT version()`, testdb.RowsFromCSVString(columns, tc.response, '|'))

		v, err := fetchVersion(db)

		if err != nil {
			t.Errorf("Expected no error, but got %s instead", err)
		}
		if err = db.Close(); err != nil {
			t.Errorf("Error '%s' was not expected while closing the database", err)
		}
		if v.first != tc.expected.first {
			t.Errorf("should be %d", tc.expected.first)
		}
		if v.second != tc.expected.second {
			t.Errorf("should be %d", tc.expected.second)
		}
		if v.third != tc.expected.third {
			t.Errorf("should be %d", tc.expected.third)
		}
	}
}
