package vlog

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func init() {
	// Change to more short name can be ingested by log
	zerolog.TimestampFieldName = "ts"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "msg"

}

type LoggerConfig struct {
	Level             zerolog.Level
	TimeFormat        string
	IncludesCaller    bool
	IncludeHTTPReqLog bool
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

// Build HTTP Log integrate
func BuildHTTPHandleLog(log zerolog.Logger) http.Handler {
	c := alice.New()

	// Install the logger handler with default output on the console
	c = c.Append(hlog.NewHandler(log))

	// Install some provided extra handler to set some request's context fields.
	// Thanks to that handler, all our logs will come with some prepopulated fields.
	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Str("type", "response").
			Dur("duration", duration).
			Msg("")
	}))
	c = c.Append(hlog.RemoteAddrHandler("remote_addr"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	// Here is your final handler
	h := c.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the logger from the request's context. You can safely assume it
		// will be always there: if the handler is removed, hlog.FromRequest
		// will return a no-op logger.
		hlog.FromRequest(r).Info().
			Str("user", "current user").
			Str("status", "ok").
			Str("type", "request").
			Msg("")

	}))
	return h
}
