package redisdb

import (
	"core/sec/conf"
	"core/sec/lib/db_connect"
	"fmt"
)

var Conn *db_connect.RedisClientWrapper

func InitRedis() {
	fmt.Println("InitRedis ...")

	var err error
	Conn, err = db_connect.CreateRedisConnection(conf.GlobalConfig.DB.Redis)
	if err != nil {
		panic(fmt.Sprint("error when open Redis connection with: ", conf.GlobalConfig.DB.Redis, "error msg: ", err.Error()))
	}

}
