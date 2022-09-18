package bot

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/richard-on/QueueBot/config"
	"github.com/richard-on/QueueBot/internal/bot/client"
	"github.com/richard-on/QueueBot/internal/logger"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"
)

func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "QueueBot is up and running!")
}

func Run() {
	log := logger.NewLogger(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC1123},
		zerolog.TraceLevel,
		"queueBot-bot")

	bot, err := tgbotapi.NewBotAPI(config.TgToken)
	if err != nil {
		log.Fatal(err, "failed to connect to Telegram bot")
	}

	bot.Debug = true
	log.Info("Authorized on account " + bot.Self.UserName)
	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 60

	database, err := sql.Open("mysql", config.DbInfo)
	if err != nil {
		log.Fatal(err, "failed to open database connection")
	}
	defer func(database *sql.DB) {
		err = database.Close()
		if err != nil {
			log.Fatal(err, "failed to gracefully close database connection")
		}
	}(database)

	http.HandleFunc("/queue-bot/healthcheck", healthcheck)
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err, "cannot establish healthcheck")
	}

	connectedUsers := make(map[int64]*client.Client, 5)
	var msg tgbotapi.MessageConfig

	updates := bot.GetUpdatesChan(upd)
	for update := range updates {
		if update.Message == nil {
			for _, connected := range connectedUsers {
				if connected.CheckTimeout() {
					delete(connectedUsers, connected.User.UserID)
				}
			}

			continue
		}

		var c *client.Client
		u := update.SentFrom()
		if conn, ok := connectedUsers[u.ID]; !ok || !conn.User.IsRegistered {
			c, err = client.NewClient(update, u, database)
			connectedUsers[u.ID] = c
		} else {
			c = conn
		}

		if !c.User.IsRegistered {
			msg.Text = client.NeedMoreInfo
			msg.ChatID = c.User.ChatID
			if _, err = bot.Send(msg); err != nil {
				log.Fatal(err, "failed to send message")
				log.Info("sent message")
			}

			continue
		} else if err != nil {
			log.Fatal(err, "Register error")
		}

		msg.Text = update.Message.Text
		msg, err = c.HandleCommand(msg)
		msg.ChatID = c.User.ChatID
		if err != nil {
			log.Fatal(err, "failed to perform bot action")
		}

		if _, err = bot.Send(msg); err != nil {
			log.Fatal(err, "failed to send message")
		}
		log.Info("sent message")
	}
}
