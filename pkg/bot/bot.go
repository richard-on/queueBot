package bot

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/richard-on/QueueBot/config"
	"github.com/richard-on/QueueBot/internal/bot/client"
	"github.com/richard-on/QueueBot/internal/logger"
	"github.com/rs/zerolog"
)

func Run() {
	log := logger.NewLogger(zerolog.TraceLevel, "queueBot-bot")

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

	connectedUsers := make(map[int64]*client.Client, 5)
	var msg tgbotapi.MessageConfig

	updates := bot.GetUpdatesChan(upd)
	for update := range updates {
		if update.Message == nil {
			go func() {
				for _, connected := range connectedUsers {
					connected.CheckTimeout()
				}
			}()

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
			msg.Text = "Неполная информация о пользователе"
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
