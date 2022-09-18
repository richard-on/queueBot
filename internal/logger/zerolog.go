package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
)

const logDir = "./logs"
const logFile = "./logs/queueBot.log"

type Logger struct {
	log zerolog.Logger
}

func NewLogger(out io.Writer, level zerolog.Level, service string) Logger {
	/*if err := os.MkdirAll(logDir, 0744); err != nil {
		panic(err)
	}

	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}*/

	return Logger{log: zerolog.New(out).
		Level(level).
		With().
		Timestamp().
		Str("service", service).
		Int("pid", os.Getpid()).
		Logger(),
	}
}

func (l Logger) Debug(i ...interface{}) {
	l.log.Debug().Msgf(fmt.Sprint(i...))
}

func (l Logger) Dedugf(format string, i ...interface{}) {
	l.log.Debug().Msgf(format, i...)
}

func (l Logger) Info(i ...interface{}) {
	l.log.Info().Msgf(fmt.Sprint(i...))
}

func (l Logger) Infof(format string, i ...interface{}) {
	l.log.Info().Msgf(format, i...)
}

func (l Logger) Error(err error, msg string) {
	l.log.Error().Err(err).Msg(msg)
}

func (l Logger) Errorf(err error, format string, i ...interface{}) {
	l.log.Error().Err(err).Msgf(format, i...)
}

func (l Logger) Fatal(err error, msg string) {
	l.log.Fatal().Err(err).Msg(msg)
}

func (l Logger) Fatalf(err error, format string, i ...interface{}) {
	l.log.Fatal().Err(err).Msgf(format, i...)
}
