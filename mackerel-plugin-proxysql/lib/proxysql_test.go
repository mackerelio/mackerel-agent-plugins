package mpproxysql

import (
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGraphDefinition(t *testing.T) {
	var proxysql ProxySQLPlugin

	graphdef := proxysql.GraphDefinition()
	if len(graphdef) != 1 {
		t.Errorf("GetTempfilename: %d should be 1", len(graphdef))
	}
}

func TestParseStatsMysqlGlobal(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Variable_Name", "Variable_Value"}).
		AddRow("ProxySQL_Uptime", "3600")

	query := "SELECT Variable_Name, Variable_Value FROM stats_mysql_global"
	mock.ExpectQuery(query).WillReturnRows(rows)

	result := make(map[string]float64)
	var proxysql ProxySQLPlugin
	err = proxysql.fetchStatsMySQLGlobal(db, result)
	if err != nil {
		t.Errorf("Failed to parse: %s", err)
	}

	expect := map[string]float64{
		"proxysql_uptime": 3600,
	}

	for k := range expect {
		if expect[k] != result[k] {
			t.Errorf("%s does not match\nExpect: %v\nResult: %v", k, expect[k], result[k])
		}
	}
}
