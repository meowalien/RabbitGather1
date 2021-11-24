package redisdb

import (
	"core/src/conf"
	"core/src/lib/db/db_connect"
	"fmt"
	"strings"
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

const Split =":"
const ProjectName = "core"

// 格式化redis key
func FormatKey(keys ...string) string {
	k := append([]string{ProjectName}, keys ...)
	return strings.Join(k , Split)
}


