package main

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/richard-on/QueueBot/cmd/queueBot/initEnv"
	_ "github.com/richard-on/QueueBot/pkg/queueBot"
	"github.com/richard-on/QueueBot/pkg/queueBot/bot"
	"github.com/richard-on/QueueBot/pkg/queueBot/db"
	"log"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	initEnv.Init()

	fmt.Println(initEnv.SentryDsn)
	err = sentry.Init(sentry.ClientOptions{
		Dsn:              initEnv.SentryDsn,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	log.Println("Initializing Database")
	err = db.InitDb()
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("reported to Sentry: %s", err)
		return
	}

	log.Println("Creating Tables")
	err = db.CreateTables()
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("reported to Sentry: %s", err)
		return
	}

	log.Println("Starting bot")
	bot.Start()

}
