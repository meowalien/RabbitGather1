package mariadb

import (
	"core/src/conf"
	"core/src/lib/db/db_connect"
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/source/file"
	"time"
)

var Conn *sql.DB

func InitMariadbDBConnection() {
	fmt.Println("InitMariadbDBConnection ...")
	var err error

	Conn, err = db_connect.CreateMysqlDBConnection(conf.GlobalConfig.DB.MarinaDB)
	waitTime := time.Second * 3
	for err != nil {
		fmt.Println("error when open MarinaDB connection with: ", conf.GlobalConfig.DB.MarinaDB, "error msg: ", err.Error())
		fmt.Printf("try again after %.2f sec.\n", waitTime.Seconds())
		time.Sleep(waitTime)
		Conn, err = db_connect.CreateMysqlDBConnection(conf.GlobalConfig.DB.MarinaDB)
	}
}
