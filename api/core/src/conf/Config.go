package conf

import (
	"core/src/lib/db/config_modules"
	"time"
)

var DEBUG_MOD = false


var GlobalConfig Config

type DB struct {
	MarinaDB config_modules.MysqlConnectConfiguration
	Redis    config_modules.RedisConfiguration
}

type RBAC struct {
	ModelConfFile   string
	PolicyTableName string
	DriverName      string
}

type JWT struct {
	//SignMethod         string
	TokenExpiresAt     int64
	TokenNotBefore     int64
	LongTokenExpiresAt int64
	Issuer             string
}

type Log struct {
	LogFile    string
	LogLevel   int
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

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
	HTTPServer    HTTPServerConfig
	RPCServer     RPCServerConfig
	AllowOrigin     []string
	AllowHeaders  []string
	ExposeHeaders []string
}

type SMTP struct {
	MailAddr string
	UserName string
	Password string
}

type Files struct {
	UploadSaveFilePath  string
	ServeStaticFilePath string
	ServeFileURL        string
}

type VCCode struct {
	Timeout int64
	Length int
}

type Config struct {
	ProjectName string
	DB DB
	RBAC RBAC
	JWT JWT
	Log Log
	Pepper             string
	Server             ServerConfig
	SMTP SMTP
	Files Files
	VCCode VCCode
}
