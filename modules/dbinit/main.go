package main

import (
	"fmt"
	"github.com/meowalien/rabbitgather-lib/db_connect"
	"gorm.io/gorm"
	"time"
)

/* -------- RBAC 權限控制 -------- */

type Role struct {
	gorm.Model
	Users []User
	Permissions []Permission`gorm:"many2many:role_permission;"`
	Title       string `gorm:"type:VARCHAR(75) NOT NULL;"`
	Slug        string `gorm:"type:VARCHAR(100) NOT NULL;,uniqueIndex,sort:asc"`
	Description string `gorm:"type:TINYTEXT NULL"`
	Active      string `gorm:"type:TINYINT(1) NOT NULL DEFAULT 0"`
	Content     string `gorm:"type:TEXT NULL DEFAULT NULL"`

}

type Permission struct {
	gorm.Model
	Title       string `gorm:"type:VARCHAR(75) NOT NULL;"`
	Slug        string `gorm:"type:VARCHAR(100) NOT NULL;,uniqueIndex,sort:asc"`
	Description string `gorm:"type:TINYTEXT NULL;"`
	Active      string `gorm:"type:TINYINT(1) NOT NULL DEFAULT 0;"`
	Content     string `gorm:"type:TEXT NULL DEFAULT NULL;"`
}


type User struct {
	gorm.Model
	RoleID uint

	FirstName string `gorm:"type:VARCHAR(50) DEFAULT NULL;"`
	MiddleName string `gorm:"type:VARCHAR(50) DEFAULT NULL;"`
	LastName string `gorm:"type:VARCHAR(50) DEFAULT NULL;"`

	Mobile string `gorm:"type:VARCHAR(15) NULL;,uniqueIndex,sort:asc"`
	Email string `gorm:"type:VARCHAR(50) NULL;,uniqueIndex,sort:asc"`

	PasswordHash string `gorm:"type:char(60) NOT NULL;"`
	PasswordSalt string `gorm:"type:char(24) NOT NULL;"`

	RegisteredAt time.Time
	LastLogin time.Time
	Intro string `gorm:"type:TINYTEXT DEFAULT NULL;"`
	Profile string `gorm:"type:TEXT DEFAULT NULL;"`
}

/* -------- 文章相關 -------- */


type Article struct {
	gorm.Model
	ArticleTags []Tag  `gorm:"many2many:article_tag;"`
	Title       string `gorm:"type:VARCHAR(75) NOT NULL;"`
	Content  string `gorm:"type:MEDIUMTEXT NOT NULL;"`
	Coords string `gorm:"type:POINT NOT NULL;"`
}
type TagType struct {
	gorm.Model
	Tags []Tag
	Name string `gorm:"type:CHAR(24) NOT NULL;,uniqueIndex,sort:asc"`
}


type Tag struct {
	gorm.Model
	Name string `gorm:"type:CHAR(24) NOT NULL;,uniqueIndex,sort:asc"`
	TagTypeID uint
}






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

	err = db.AutoMigrate(&Role{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&Permission{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err.Error())
	}

	err = db.AutoMigrate(&Article{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&TagType{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(&Tag{})
	if err != nil {
		panic(err.Error())
	}
}
