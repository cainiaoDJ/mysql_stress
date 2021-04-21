package main

import (
	"fmt"
	"github.com/op/go-logging"
	"gorm.io/gorm"
	"math/rand"
	"mysql_stress/config"
	"mysql_stress/excel"
	"mysql_stress/msql"
	"mysql_stress/utils"
	"os"
	"strconv"
	"sync"
	"time"
)

// 1,100
// 10,100
// 100,100
// 500万初始化
// 测试写入100万数据花费时间，（顺序，并发100，并发1000）
// 随机读取100万数据花费时间（随机5次，求平均数据）（顺序，并发100，并发1000）

var format = logging.MustStringFormatter(
	`%{level:.4s} <%{shortfunc}> %{id:03x} %{message}`,
)

var dbObj *gorm.DB

func main() {
	filepath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	config.LoadAppConfig(filepath + "/config")
	logging.SetFormatter(format)

	utils.AppLog.Info("init Test >", "DB:", config.Cfg.Host)
	dbObj = msql.GetDBConnects()
	//insertDataTest(config.Cfg.DBNum, config.Cfg.TableNum)
	var xlsData []excel.RowData
	var tmp excel.RowData

	mdb, err := dbObj.DB()
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			st := mdb.Stats()
			utils.AppLog.Infof("Idle:%3d, InUse:%3d, OpenConnections:%3d", st.Idle, st.InUse, st.OpenConnections)
			time.Sleep(time.Second)
		}
	}()
	mdb.SetMaxOpenConns(int(config.Cfg.MaxDBConns))
	mdb.SetConnMaxLifetime(time.Duration(config.Cfg.MaxDBConnsLifetime) * time.Second)
	defer msql.CleanTestDB(dbObj)

	// 先从高并发往低并发走
	for _, n := range []uint{300, 200, 100, 50, 10} {
		for _, dn := range []uint{1, 2, 10, 20, 100} {
			for _, tn := range []uint{1, 10, 20, 100, 200} {
				if tn < dn {
					continue
				}
				utils.AppLog.Infof("Start test. DB:%d,TB:%d,Routine:%d", dn, tn, n)
				tmp = insertDataTest(dn, tn, n)
				tmp.Func = "InsertDataTest"
				xlsData = append(xlsData, tmp)
				tmp = ReadDataTest(dn, tn, n)
				tmp.Func = "ReadDataTest"
				xlsData = append(xlsData, tmp)
				tmp = UpdateDataTest(dn, tn, n)
				tmp.Func = "UpdateDataTest"
				xlsData = append(xlsData, tmp)
			}
		}
	}
	utils.AppLog.Info("test end!")
	err = excel.WriteToExcel(xlsData, "mysql_stress")
	if err != nil {
		utils.AppLog.Errorf("write to excel failed. err:%s ", err.Error())
		return
	}
	utils.AppLog.Info("write result to excel")
}

func insertDataTest(DBNum uint, tableNum uint, connNum uint) excel.RowData {
	msql.InitDB(DBNum)
	createSQL := msql.GetCreateTableSQL(DBNum, tableNum)
	//dbs = msql.GetDBConnects(DBNum)
	//// init db
	//defer func() {
	//	for _, db := range dbs {
	//		err := db.Close()
	//		if err != nil {
	//			utils.AppLog.Error(err.Error())
	//		}
	//	}
	//	time.Sleep(1 * time.Second)
	//}()
	//for i := uint(0); i < DBNum; i++ {
	//	for j := uint(0); j < tableNum; j++ {
	//		if j%DBNum == i {
	//			_, err := mdb.Exec(createSQL[j])
	//			if err != nil {
	//				panic("create table err:" + err.Error())
	//			}
	//			// time.Sleep(1 * time.Millisecond)
	//		}
	//	}
	//}
	for _, csql := range createSQL {
		tx := dbObj.Exec(csql)
		if tx.Error != nil {
			panic("create table err:" + tx.Error.Error())
		}
	}
	utils.AppLog.Infof("table create complete! tb:%d", tableNum)

	startTime := time.Now()
	var step = config.Cfg.InitRows / connNum
	var wg sync.WaitGroup

	for pid := uint(0); pid < connNum; pid++ {
		wg.Add(1)
		startID := 100000000 + pid*step
		endId := startID + step
		//utils.AppLog.Info("insert data,pid:", pid)
		go func(start uint, end uint) {
			defer wg.Done()
			for uid := start; uid < end; uid++ {
				db := uid % DBNum
				tb := uid % tableNum
				tbSQL := fmt.Sprintf(msql.INSERT_SQL, db, tb, uid, uid, "用户名"+strconv.Itoa(int(uid)))
				tx := dbObj.Exec(tbSQL)
				if tx.Error != nil || tx.RowsAffected != 1 {
					panic("insert data err:" + tx.Error.Error() +
						", rows_affected:" + strconv.Itoa(int(tx.RowsAffected)))
				}
				// time.Sleep(1 * time.Millisecond)
			}
		}(startID, endId)
	}
	wg.Wait()
	cost := time.Since(startTime)
	utils.AppLog.Debugf("db:%2d, tb:%2d, routine:%2d, cost:%8.4f s, speed:%7.4f tps",
		DBNum, tableNum, connNum, cost.Seconds(), float64(config.Cfg.InitRows)/cost.Seconds())
	return excel.RowData{
		Routine: connNum,
		DBNum:   DBNum,
		TbNum:   tableNum,
		Cost:    cost.Seconds(),
		Speed:   float64(config.Cfg.InitRows) / cost.Seconds(),
	}
}

func ReadDataTest(DBNum uint, tableNum uint, connNum uint) excel.RowData {
	//dbs := msql.GetDBConnects(DBNum)
	//defer func() {
	//	for _, db := range dbs {
	//		err := db.Close()
	//		if err != nil {
	//			utils.AppLog.Error(err.Error())
	//		}
	//	}
	//	time.Sleep(1 * time.Second)
	//}()

	var wg sync.WaitGroup
	startTime := time.Now()
	for pid := uint(0); pid < connNum; pid++ {
		wg.Add(1)
		//utils.AppLog.Info("insert data,pid:", pid)
		go func() {
			defer wg.Done()
			for i := uint(0); i < config.Cfg.InitRows/connNum; i++ {
				random := rand.New(rand.NewSource(time.Now().UnixNano()))
				uid := uint(random.Int63n(int64(config.Cfg.InitRows))) + 100000000
				db := uid % DBNum
				tb := uid % tableNum
				tbSQL := fmt.Sprintf(msql.QUERY_SQL, db, tb, uid)
				tx := dbObj.Exec(tbSQL)
				if tx.Error != nil {
					panic("query data err:" + tx.Error.Error())
				}
				// time.Sleep(1 * time.Millisecond)
			}
		}()
	}
	wg.Wait()
	cost := time.Since(startTime)
	utils.AppLog.Debugf("db:%2d, tb:%2d, routine:%2d, cost:%8.4f s, speed:%8.4f tps",
		DBNum, tableNum, connNum, cost.Seconds(), float64(config.Cfg.InitRows)/cost.Seconds())
	return excel.RowData{
		Routine: connNum,
		DBNum:   DBNum,
		TbNum:   tableNum,
		Cost:    cost.Seconds(),
		Speed:   float64(config.Cfg.InitRows) / cost.Seconds(),
	}
}

func UpdateDataTest(DBNum uint, tableNum uint, connNum uint) excel.RowData {
	//dbs := msql.GetDBConnects(DBNum)
	//defer func() {
	//	for _, db := range dbs {
	//		err := db.Close()
	//		if err != nil {
	//			utils.AppLog.Error(err.Error())
	//		}
	//	}
	//	time.Sleep(1 * time.Second)
	//}()

	var wg sync.WaitGroup
	startTime := time.Now()
	for pid := uint(0); pid < connNum; pid++ {
		wg.Add(1)
		//utils.AppLog.Info("insert data,pid:", pid)
		go func() {
			defer wg.Done()
			for i := uint(0); i < config.Cfg.InitRows/connNum; i++ {
				random := rand.New(rand.NewSource(time.Now().UnixNano()))
				uid := uint(random.Int63n(int64(config.Cfg.InitRows))) + 100000000
				db := uid % DBNum
				tb := uid % tableNum
				tbSQL := fmt.Sprintf(msql.UPDATE_SQL, db, tb, "update test", uid)
				tx := dbObj.Exec(tbSQL)
				if tx.Error != nil {
					panic("update data err:" + tx.Error.Error())
				}
			}
		}()
	}
	wg.Wait()
	cost := time.Since(startTime)
	utils.AppLog.Debugf("db:%2d, tb:%2d, routine:%2d, cost:%8.4f s, speed:%8.4f tps",
		DBNum, tableNum, connNum, cost.Seconds(), float64(config.Cfg.InitRows)/cost.Seconds())
	return excel.RowData{
		Routine: connNum,
		DBNum:   DBNum,
		TbNum:   tableNum,
		Cost:    cost.Seconds(),
		Speed:   float64(config.Cfg.InitRows) / cost.Seconds(),
	}
}
