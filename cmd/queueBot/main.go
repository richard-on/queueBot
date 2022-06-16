package main

import (
	"log"
	"queueBot/pkg/queueBot/bot"
	"queueBot/pkg/queueBot/db"
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
