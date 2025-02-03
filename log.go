package vlog

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog"
)

func init() {
	// Change to more short name can be ingested by log
	zerolog.TimestampFieldName = "ts"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "msg"

}

type LoggerConfig struct {
	Level          zerolog.Level
	TimeFormat     string
	IncludesCaller bool
	// HTTPLoggerConfig  *HTTPLoggerConfig
}

type HTTPLoggerConfig struct {
	EnableDurationReq bool
	ReqTraceLog       struct {
		Enable           bool
		Base64DataEncode string
	}
	RespTraceLog struct {
		Enable           bool
		Base64DataEncode string
	}
}

type JsonLogger struct {
	Logger *zerolog.Logger
}

type ConsoleLogger struct {
	logger *zerolog.Logger
}

func createZeroLogInst(c LoggerConfig) zerolog.Logger {
	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Logger()
	logger = logger.With().Timestamp().Logger()
	if c.IncludesCaller {
		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			return filepath.Base(file) + ":" + strconv.Itoa(line)
		}
		logger = logger.With().Caller().Logger()
	}
	logger = logger.Level(c.Level)
	return logger
}

func NewConsoleLogger(c LoggerConfig) *ConsoleLogger {
	return &ConsoleLogger{
		logger: nil,
	}
}
func NewJsonLogger(c LoggerConfig) JsonLogger {
	l := createZeroLogInst(c)
	return JsonLogger{
		Logger: &l,
	}
}
