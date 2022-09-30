package logger

import (
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/richard-on/queueBot/internal/bot/model"
	"github.com/rs/zerolog"
	"io"
	"os"
	"time"
)

//const logDir = "./logs"
//const logFile = "./logs/queueBot.log"

type TGUpdate struct {
}

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

func (l Logger) TgUpdate(upd tgbotapi.Update) {
	jsonUpd, err := json.Marshal(upd)
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	l.log.Info().
		Str("endpoint", "getUpdates").
		RawJSON("update", jsonUpd).
		Msg("")
}

func (l Logger) TgSend(msg tgbotapi.Message) {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	l.log.Info().
		Str("endpoint", "sendMessage").
		RawJSON("message", jsonMsg).
		Msg("")
}

func (l Logger) ClientDisconnect(user model.User, lastConn time.Time) {
	jsonUser, err := json.Marshal(user)
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}

	l.log.Info().
		Str("endpoint", "clientDisconnect").
		RawJSON("user", jsonUser).
		Str("lastConn", lastConn.String()).
		Msg("client disconnected")
}

func (l Logger) Debug(i ...interface{}) {
	l.log.Debug().Caller().Msgf(fmt.Sprint(i...))
}

func (l Logger) Dedugf(format string, i ...interface{}) {
	l.log.Debug().Caller().Msgf(format, i...)
}

func (l Logger) Info(i ...interface{}) {
	l.log.Info().Msgf(fmt.Sprint(i...))
}

func (l Logger) Infof(format string, i ...interface{}) {
	l.log.Info().Msgf(format, i...)
}

func (l Logger) Error(err error, msg string) {
	l.log.Error().Caller().Err(err).Msg(msg)
}

func (l Logger) Errorf(err error, format string, i ...interface{}) {
	l.log.Error().Caller().Err(err).Msgf(format, i...)
}

func (l Logger) Fatal(err error, msg string) {
	l.log.Fatal().Caller().Err(err).Msg(msg)
}

func (l Logger) Fatalf(err error, format string, i ...interface{}) {
	l.log.Fatal().Caller().Err(err).Msgf(format, i...)
}
