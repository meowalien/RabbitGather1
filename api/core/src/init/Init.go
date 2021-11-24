package init

import (
	"core/src/conf"
	"core/src/db/mariadb"
	"core/src/db/redisdb"
	"core/src/lib/config"
	"core/src/module/log"
	"core/src/module/permission"
	"flag"
	"fmt"
	"github.com/kr/pretty"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}


func InitFlags() {
	fmt.Println("InitFlags ...")
	flag.BoolVar(&conf.DEBUG_MOD, "debug", false, "\"true\" to open debug mode")
	flag.StringVar(&ConfigFile , "config" , _DefaultConfigFile , "to set the config file")

	flag.Parse()

	fmt.Println("DEBUG_MOD: ", conf.DEBUG_MOD)
	fmt.Println("ConfigFile: ", ConfigFile)
}

const _DefaultConfigFile = "config/config.json"
var ConfigFile string

func InitConfig() {
	fmt.Println("InitConfig ...")
	err := config.JsonConfigModleMapping(&conf.GlobalConfig, ConfigFile)
	if err != nil {
		panic(err.Error())
	}
	_, e := pretty.Println(conf.GlobalConfig)
	if e != nil {
		fmt.Println("error when printing GlobalConfig by pretty:", err.Error())
		fmt.Println(conf.GlobalConfig)
	}
}


func init() {
	InitFlags()
	InitConfig()
	redisdb.InitRedis()
	mariadb.InitMariadbDBConnection()
	permission.InitRBAC(mariadb.Conn)
	log.InitLogger()

}