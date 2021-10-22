package conf

import (
	"core/sec/lib/db_connect"
	"time"
)

var DEBUG_MOD = false

type HTTPServerConfig struct {
	Name             string
	Host             string
	Port             string
	ShutdownWaitTime time.Duration
}

type RPCServerConfig struct {
	Name string
	Host string
	Port string
}

type ServerConfig struct {
	HTTPServer HTTPServerConfig
	RPCServer  RPCServerConfig
}

type DatabaseConfig struct {
	MarinaDB db_connect.MysqlConnectConfiguration
	Redis    db_connect.RedisConfiguration
}

var GlobalConfig Config

type RBACConfig struct {
	ModelConfFile   string
	PolicyTableName string
	DriverName      string
}

type JWTConfig struct {
	SignMethod string
	Pepper     string
}

type SZFU struct {
	Hashid string
	Hashkey string
	URL string
}

// ArticleConfig is the root config struct
type Config struct {
	Server             ServerConfig
	DB                 DatabaseConfig
	RBAC               RBACConfig
	JWT                JWTConfig
	Pepper             string
	PasswordEncryption bool
	SZFU SZFU
}
