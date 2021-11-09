package config_modules

import (
	_ "github.com/go-sql-driver/mysql"
)

type MysqlConnectConfiguration struct {
	Host     string
	Database string
	User     string
	Password string
	Port     string
}
