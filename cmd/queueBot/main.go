package main

import (
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/richard-on/QueueBot/config"
	"github.com/richard-on/QueueBot/internal/logger"
	"github.com/richard-on/QueueBot/pkg/bot"
	"github.com/rs/zerolog"
	"os"
	"time"
)

func main() {
	log := logger.NewLogger(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC1123},
		zerolog.TraceLevel,
		"queueBot-setup")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err, "cannot load env variables")
	}

	config.Init()

	log.Info(config.SentryDsn)
	err = sentry.Init(sentry.ClientOptions{
		Dsn:              config.SentryDsn,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatal(err, "Sentry init failed")
	}
	defer sentry.Flush(2 * time.Second)

	/*log.Info("initializing database")
	err = db.InitDb()
	if err != nil {
		sentry.CaptureException(err)
		log.Error(err, "")
		return
	}

	log.Info("creating tables")
	err = db.CreateTables()
	if err != nil {
		sentry.CaptureException(err)
		log.Error(err, "")
		return
	}*/

	log.Info("starting bot")
	bot.Run()
}
