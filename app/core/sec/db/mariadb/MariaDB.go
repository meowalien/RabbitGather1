package mariadb

import (
	"core/sec/conf"
	"core/sec/lib/db_connect"
	"database/sql"
	"fmt"

	"time"
)

var Conn *sql.DB

func InitMariadbDBConnection() {
	fmt.Println("InitMariadbDBConnection ...")
	var err error

	Conn, err = db_connect.CreateMysqlDBConnection(conf.GlobalConfig.DB.MarinaDB)
	// 有可能資料庫還沒開好
	waitTime := time.Second * 3
	for err != nil {
		fmt.Println("error when open MarinaDB connection with: ", conf.GlobalConfig.DB.MarinaDB, "error msg: ", err.Error())
		fmt.Printf("try again after %.2f sec.\n", waitTime.Seconds())
		time.Sleep(waitTime)
		Conn, err = db_connect.CreateMysqlDBConnection(conf.GlobalConfig.DB.MarinaDB)
	}

}
