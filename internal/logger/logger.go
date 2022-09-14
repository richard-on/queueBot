package logger

import (
	"github.com/rs/zerolog"
	"os"
)

const logDir = "./logs"
const logFile = "./logs/queueBot.log"

type Log struct {
	Logger zerolog.Logger
}

func NewLogger(level zerolog.Level, service string) *Log {
	if err := os.MkdirAll(logDir, 0744); err != nil {
		panic(err)
	}

	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	return &Log{Logger: zerolog.New(file /*zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC1123}*/).
		Level(level).
		With().
		Timestamp().
		Str("service", service).
		Int("pid", os.Getpid()).
		Logger(),
	}
}

func (l *Log) Fatal(err error, msg string) {
	l.Logger.Fatal().Err(err).Msg(msg)
}

func (l *Log) Error(err error, msg string) {
	l.Logger.Error().Err(err).Msg(msg)
}

func (l *Log) Info(msg string) {
	l.Logger.Info().Msg(msg)
}
