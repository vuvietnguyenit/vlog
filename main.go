package main

import (
	"net/http"

	"github.com/rs/zerolog"
)

func main() {
	jsonLogger := NewJsonLogger(LoggerConfig{
		Level:          zerolog.DebugLevel,
		TimeFormat:     "2006-01-02T15:04:05Z07:00",
		IncludesCaller: true,
	})

	// Create HTTP handler log
	h := BuildHTTPHandleLog(*jsonLogger.Logger)
	http.Handle("/", h)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		print("Startup failed")
	}
}
