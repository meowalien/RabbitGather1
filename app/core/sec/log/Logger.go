package log

import (
	"core/sec/conf"
	"core/sec/lib/text"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

var Logger *logrus.Entry

func InitLogger() {
	logLevel := logrus.InfoLevel
	if conf.DEBUG_MOD {
		fmt.Println("DEBUG_MOD is on, set logLevel to DebugLevel")
		logLevel = logrus.DebugLevel
	}

	Logger = CreateLogger(&Config{
		LogLevel: logLevel,
		Fields: map[string]interface{}{
			"app": "PlatformProject/core",
		},
	})
}

type FormatterType string

const (
	JsonFormatterType FormatterType = "json"
)

type Config struct {
	DisableColorEncoding bool
	// FormatterType default json
	FormatterType          FormatterType
	Writers                []io.Writer
	DisablePrintOnTerminal bool
	LogLevel               logrus.Level
	Fields                 map[string]interface{}
}

func CreateLogger(cf *Config) *logrus.Entry {
	logger := logrus.New()

	var formatter logrus.Formatter
	switch cf.FormatterType {
	case "":
		fallthrough
	case JsonFormatterType:
		formatter = &logrus.JSONFormatter{
			// time格式
			TimestampFormat: time.StampNano,
			PrettyPrint:     true,
		}
	default:
		panic("unsupported formatter type: " + cf.FormatterType)
	}

	logger.SetFormatter(&MyFormatter{
		ColorEncoding: !cf.DisableColorEncoding,
		Formatter:     formatter,
	})

	logger.SetReportCaller(true)

	var writer io.Writer
	if cf.DisablePrintOnTerminal {
		writer = io.Discard
	} else if cf.Writers == nil || len(cf.Writers) == 0 {
		writer = os.Stdout
	} else {
		w := []io.Writer{
			os.Stdout,
		}
		writer = io.MultiWriter(append(w, cf.Writers...)...)
	}

	logger.SetOutput(writer)
	var logLevel logrus.Level
	if cf.LogLevel == 0 {
		logLevel = logrus.DebugLevel
	} else {
		logLevel = logrus.Level(cf.LogLevel)
	}
	logger.SetLevel(logLevel)
	return logger.WithFields(cf.Fields)
}

type MyFormatter struct {
	logrus.Formatter
	ColorEncoding bool
}

func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b, e := f.Formatter.Format(entry)
	if !f.ColorEncoding {
		return b, e
	}
	switch entry.Level {
	case logrus.PanicLevel:
		fallthrough
	case logrus.FatalLevel:
		fallthrough
	case logrus.ErrorLevel:
		return text.ColorByteSting(b, text.FgRed), e
	case logrus.WarnLevel:
		return text.ColorByteSting(b, text.FgYellow), e
	case logrus.InfoLevel:
		return text.ColorByteSting(b, text.FgGreen), e
	case logrus.DebugLevel:
		return text.ColorByteSting(b, text.FgBlue), e
	default:
		return b, e
	}
}
