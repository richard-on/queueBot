package main

import (
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/richard-on/QueueBot/config"
	"github.com/richard-on/QueueBot/pkg/bot"
	"github.com/rs/zerolog"
	"os"
	"time"
)

func main() {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC1123}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Logger()

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	config.Init()
	/*config.TgToken = os.Getenv("TOKEN")
	config.SentryDsn = os.Getenv("SENTRY_DSN")
	config.Host = os.Getenv("HOST")
	config.Port = os.Getenv("PORT")
	config.User = os.Getenv("USER")
	config.Password = os.Getenv("PASSWORD")
	config.DbName = os.Getenv("DBNAME")

	config.DbInfo = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.User, config.Password, config.Host, config.Port, config.DbName)

	config.InitDbInfo = fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		config.User, config.Password, config.Host, config.Port)*/

	log.Info().Msg(config.SentryDsn)
	err = sentry.Init(sentry.ClientOptions{
		Dsn:              config.SentryDsn,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Sentry init failed")
	}
	defer sentry.Flush(2 * time.Second)

	/*log.Info().Msg("initializing database")
	err = db.InitDb()
	if err != nil {
		sentry.CaptureException(err)
		log.Info().Msgf("reported to Sentry: %s", err)
		return
	}

	log.Info().Msg("creating tables")
	err = db.CreateTables()
	if err != nil {
		sentry.CaptureException(err)
		log.Info().Msgf("reported to Sentry: %s", err)
		return
	}*/

	log.Info().Msg("starting bot")
	bot.Run(log)
}
