package main

import (
	"log"
	bot "queueBot/src/bot"
)

func main() {
	log.Println("Initializing Database...")
	err := bot.InitDb()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Creating Tables...")
	err = bot.CreateTables()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting bot...")
	bot.Bot()

}
