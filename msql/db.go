package msql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mysql_stress/utils"
	//_ "github.com/go-msql-driver/mysql"
	//"xorm.io/xorm"
	"mysql_stress/config"
)

const INSERT_SQL = "INSERT INTO `db_test_%02d`.`u_player_%02d` VALUES (%d,%d,2,'%s','省略...',NULL,'大家好！',1,NULL,4,1,20,1,2,0,0,0,'2019-10-25 17:21:29','2019-10-25 16:35:27','2019-10-25 17:21:28','2019-10-25 17:21:29','2019-10-25 16:37:52',NULL,'2019-10-25 16:25:32','2019-10-25 17:21:29',NULL,NULL,35,1571991932,'Controller_Api_XXX',1571995289,40001);"
const QUERY_SQL = "select * from `db_test_%02d`.`u_player_%02d` where player_id = %d;"
const UPDATE_SQL = "update `db_test_%02d`.`u_player_%02d` set player_name = '%s' where player_id = %d;"

func cleanDBs() (sql []string) {
	var tmpl = "DROP DATABASE if EXISTS db_test_%02d;"
	for i := uint(0); i < 100; i++ {
		sql = append(sql, fmt.Sprintf(tmpl, i))
	}
	return
}

func getCreateDBSQL(num uint) (sql []string) {
	var tmpl = "CREATE DATABASE if NOT EXISTS db_test_%02d;"
	for i := uint(0); i < num; i++ {
		sql = append(sql, fmt.Sprintf(tmpl, i))
	}
	return
}

func GetCreateTableSQL(DBNum uint, tbNum uint) (sql []string) {
	var tmpl = "CREATE TABLE `db_test_%02d`.`u_player_%02d` (\n\t`player_id` BIGINT(20) UNSIGNED NOT NULL,\n\t`public_id` BIGINT(20) UNSIGNED NOT NULL,\n\t`player_rank` SMALLINT(5) UNSIGNED NOT NULL,\n\t`player_name` TINYTEXT NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci',\n\t`heroine_name` TINYTEXT NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci',\n\t`birthday` DATE NULL DEFAULT NULL,\n\t`self_introduction` TINYTEXT NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci',\n\t`title_id` SMALLINT(5) UNSIGNED NOT NULL,\n\t`title_id2` SMALLINT(5) UNSIGNED NULL DEFAULT NULL,\n\t`platform_id` TINYINT(3) UNSIGNED NOT NULL,\n\t`task_rookie_banner_type` TINYINT(3) UNSIGNED NOT NULL,\n\t`main_tutorial_id` SMALLINT(5) NOT NULL,\n\t`main_tutorial_clear_flag` TINYINT(1) NOT NULL,\n\t`scene_tutorial_id` BIGINT(20) NOT NULL,\n\t`clear_section_id` TINYINT(3) UNSIGNED NOT NULL,\n\t`test_device_flag` TINYINT(1) NOT NULL,\n\t`account_status` TINYINT(3) UNSIGNED NOT NULL,\n\t`last_access_date` DATETIME NULL DEFAULT NULL,\n\t`last_login_date` DATETIME NULL DEFAULT NULL,\n\t`last_open_check_date` DATETIME NULL DEFAULT NULL,\n\t`last_event_check_date` DATETIME NULL DEFAULT NULL,\n\t`first_purchase_date` DATETIME NULL DEFAULT NULL,\n\t`dormant_return_days` TINYINT(3) UNSIGNED NULL DEFAULT NULL,\n\t`in_date` DATETIME NOT NULL,\n\t`up_date` DATETIME NOT NULL,\n\t`thaw_date` DATETIME NULL DEFAULT NULL COMMENT '解冻时间',\n\t`frozen_content` TINYTEXT NULL DEFAULT NULL COMMENT '冻结说明' COLLATE 'utf8mb4_general_ci',\n\t`channel` INT(10) UNSIGNED NOT NULL COMMENT '渠道',\n\t`reg_time` INT(10) UNSIGNED NOT NULL COMMENT '注册时间戳',\n\t`last_ac` VARCHAR(255) NOT NULL COMMENT '最后一次行为' COLLATE 'utf8mb4_general_ci',\n\t`last_ac_time` INT(10) UNSIGNED NOT NULL COMMENT '最后一次执行时间',\n\t`server_id` INT(10) UNSIGNED NOT NULL COMMENT '区服id',\n\tPRIMARY KEY (`player_id`) USING BTREE,\n\tINDEX `u_player_index1` (`in_date`) USING BTREE,\n\tINDEX `idx_last_ac_time` (`last_ac_time`) USING BTREE\n)\nCOLLATE='utf8mb4_general_ci'\nENGINE=InnoDB\n;\n"
	for j := uint(0); j < tbNum; j++ {
		sql = append(sql, fmt.Sprintf(tmpl, j%DBNum, j))
	}
	return
}
func InitDB(DBNum uint) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8",
		config.Cfg.Username,
		config.Cfg.Password,
		config.Cfg.Host,
		config.Cfg.Port,
	)

	//db, err := xorm.NewEngine("mysql", dsn)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("db connect failed:" + err.Error())
	}

	cleanSQL := cleanDBs()
	for _, s := range cleanSQL {
		tx := db.Exec(s)
		if tx.Error != nil {
			panic("create dbs failed:" + tx.Error.Error())
		}
		//db.Exec(s)
	}
	//time.Sleep(5 * time.Second)
	createSQL := getCreateDBSQL(DBNum)
	for _, s := range createSQL {
		tx := db.Exec(s)
		if tx.Error != nil {
			panic("create dbs failed:" + tx.Error.Error())
		}
		//db.Exec(s)
	}
	utils.AppLog.Infof("DB init complete,db:%d", DBNum)
}

func GetDBConnects() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8",
		config.Cfg.Username,
		config.Cfg.Password,
		config.Cfg.Host,
		config.Cfg.Port,
	)
	//for i := uint(0); i < num; i++ {
	//db, err := xorm.NewEngine("mysql", fmt.Sprintf(dsn, i))
	dbConn, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		panic("db connect failed:" + err.Error())
	}
	//db, err := dbConn.DB()
	//if err != nil {
	//	panic("db connect failed:" + err.Error())
	//}
	//db.SetMaxOpenConns(1)
	//db.SetMaxIdleConns(0)

	//}
	return dbConn
}
