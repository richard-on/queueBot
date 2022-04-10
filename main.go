package main

import (
	"log"
	"queueBot/bot"
)

func main() {

	err := bot.InitDb()
	if err != nil {
		log.Fatal(err)
	}

	err = bot.CreateTables()
	if err != nil {
		log.Fatal(err)
	}

	bot.Bot()

}
