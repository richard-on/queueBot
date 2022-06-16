package main

import (
	"log"
	bot2 "queueBot/cmd/queueBot/bot"
)

func main() {
	log.Println("Initializing Database...")
	err := bot2.InitDb()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Creating Tables...")
	err = bot2.CreateTables()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting bot...")
	bot2.Bot()

}
