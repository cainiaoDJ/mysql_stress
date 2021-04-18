package msql

import (
	"mysql_stress/config"
	"testing"
)

func TestGetDBEngines(t *testing.T) {
	config.LoadAppConfig("../config")
	DBs := GetDBConnects(1)
	result, err := DBs[0].Exec("show databases;")
	if err != nil {
		t.Fail()
	}
	if n, _ := result.RowsAffected(); n > 0 {
		t.Fail()
	}

}
