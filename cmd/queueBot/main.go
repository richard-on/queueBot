package main

import (
	"github.com/richard-on/QueueBot/pkg/queueBot/bot"
	"github.com/richard-on/QueueBot/pkg/queueBot/db"
	"log"
)

func main() {
	log.Println("Initializing Database")
	err := db.InitDb()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Creating Tables")
	err = db.CreateTables()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting bot")
	bot.Start()

}
