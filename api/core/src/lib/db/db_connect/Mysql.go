package db_connect

import (
	"core/src/lib/db/config_modules"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func CreateMysqlDBConnection(dbconf config_modules.MysqlConnectConfiguration) (*sql.DB, error) {
	dsn := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn = fmt.Sprintf(dsn, dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port, dbconf.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateGORMMysqlDBConnection(db *sql.DB) (*gorm.DB, error) {

	dbconn, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}
	rawdb, err := dbconn.DB()
	if err != nil {
		return nil, err
	}
	err = rawdb.Ping()
	if err != nil {
		return nil, err
	}
	return dbconn, err
}

func CloseConn(dbconn *gorm.DB) error {
	rawDbconn, err := dbconn.DB()
	if err != nil {
		return err
	}
	err = rawDbconn.Close()
	if err != nil {
		return err
	}
	return nil
}
