package main

import (
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/richard-on/queueBot/config"
	"github.com/richard-on/queueBot/pkg/bot"
	"github.com/richard-on/queueBot/pkg/logger"
	"github.com/rs/zerolog"
	"os"
	"time"
)

func main() {
	var err error
	log := logger.NewLogger(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC1123},
		zerolog.TraceLevel,
		"queueBot-setup")

	err = godotenv.Load()
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

	log.Info("starting bot")
	bot.Run()
}
