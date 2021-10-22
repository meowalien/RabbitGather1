package db_connect

import (
	"database/sql"
	"fmt"
)

func CreateMysqlDBConnection(dbconf MysqlConnectConfiguration) (*sql.DB, error) {
	dsn := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn = fmt.Sprintf(dsn, dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port, dbconf.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	//db.SetMaxOpenConns(20)
	//db.SetMaxIdleConns(10)
	//db.SetConnMaxLifetime(time.Minute * 10)

	return db, nil
}


type MysqlConnectConfiguration struct {
	Host     string
	Database string
	User     string
	Password string
	Port     string
}
