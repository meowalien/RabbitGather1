package log

import (
	"core/src/conf"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *MyLog

type MyLog struct {
	*zap.SugaredLogger
}

func (m *MyLog)Skip(deap int ) *MyLog{
	return &MyLog{SugaredLogger:m.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(deap)).Sugar()}
}

func InitLogger() {
	fmt.Println("InitLogger ...")
	var level zapcore.Level
	switch conf.GlobalConfig.Log.LogLevel {
	case -1:
		level = zapcore.DebugLevel
	case 1:
		level = zapcore.InfoLevel
	case 2:
		level = zapcore.WarnLevel
	case 3:
		level = zapcore.ErrorLevel
	default:
		panic("not supported error")
	}
	Logger = &MyLog{SugaredLogger: NewLogger(level)}
}

func NewLogger(level zapcore.Level) *zap.SugaredLogger {
	// 限制日誌輸出級別, >= DebugLevel 會打印所有級別的日誌
	// 生產環境中一般使用 >= ErrorLevel
	levelEnabler := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= level
	})

	// 使用 JSON 格式日誌
	jsonEnc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	stdCore := zapcore.NewCore(jsonEnc, zapcore.Lock(os.Stdout), levelEnabler)

	lumberJackLogger := &lumberjack.Logger{
		Filename:   conf.GlobalConfig.Log.LogFile,
		MaxSize:    conf.GlobalConfig.Log.MaxSize,    // 最大文件大小 500 MB
		MaxBackups: conf.GlobalConfig.Log.MaxBackups, // 最多保留 5個文件
		MaxAge:     conf.GlobalConfig.Log.MaxAge,     // 保存30天內的日誌
		//Compress:   false, //壓縮
	}
	// addSync 將 io.Writer 裝飾爲 WriteSyncer
	//	// 故只需要一個實現 io.Writer 接口的對象即可
	syncer := zapcore.AddSync(lumberJackLogger)

	redisCore := zapcore.NewCore(jsonEnc, zapcore.Lock(syncer), levelEnabler)

	// 集成多個 core
	core := zapcore.NewTee(redisCore,stdCore) //stdCore

	// logger 輸出到 console 且標識調用代碼行
	return zap.New(core).WithOptions(zap.AddCaller()).Sugar()
}
