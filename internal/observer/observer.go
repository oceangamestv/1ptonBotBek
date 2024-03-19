package observer

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

type Observer struct {
	Logger *zerolog.Logger
}

func New(logLevel zerolog.Level, debug bool) *Observer {
	return &Observer{
		Logger: newLogger(logLevel, debug),
	}
}

func newLogger(level zerolog.Level, debug bool) *zerolog.Logger {
	var logger zerolog.Logger
	if debug == true {
		logger = zerolog.New(
			zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
		)
	} else {
		logger = zerolog.New(os.Stderr)
	}

	logger = logger.
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()
	return &logger
}
