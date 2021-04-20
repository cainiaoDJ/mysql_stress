package msql

import (
	"fmt"
	"mysql_stress/config"
	"testing"
)

type Result struct {
	Variable_name string
	Value         string
}

func TestGetDBEngines(t *testing.T) {
	config.LoadAppConfig("../config")
	db := GetDBConnects()
	var ret []Result
	tx := db.Raw("show variables like '%wait_timeout%'").Scan(&ret)
	if tx.Error != nil {
		t.Fail()
	}
	fmt.Println(ret)
}
