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
