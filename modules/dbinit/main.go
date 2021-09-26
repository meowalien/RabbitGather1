package main

import (
	"fmt"
	"github.com/meowalien/rabbitgather-lib/db_connect"
	"github.com/meowalien/rabbitgather-lib/table_struct"
)

func main() {
	fmt.Println("InitMariadbDB ...")
	conf := db_connect.MysqlConnectConfiguration{
		Host:     "localhost",
		Database: "rabbit_gather",
		User:     "rabbit_gather",
		Password: "5678",
		Port:     "3306",
	}
	GlobalConn, err := db_connect.CreateGormDBConnection(conf)
	if err != nil {
		panic(fmt.Sprint("error when open MarinaDB connection with: ", conf, "error msg: ", err.Error()))
	}
	db := GlobalConn.Set("gorm:table_options", "ENGINE=InnoDB")

	err = db.AutoMigrate(&table_struct.Role{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&table_struct.Permission{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&table_struct.User{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&table_struct.Article{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&table_struct.TagType{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&table_struct.Tag{})
	if err != nil {
		panic(err.Error())
	}
}
