package bot

import (
	"database/sql"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/richard-on/QueueBot/config"
	"github.com/richard-on/QueueBot/internal"
	"github.com/richard-on/QueueBot/internal/bot/client"
	"github.com/rs/zerolog"
)

func Run(log zerolog.Logger) {
	bot, err := tgbotapi.NewBotAPI(config.TgToken)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to Telegram bot")
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 60

	database, err := sql.Open("mysql", config.DbInfo)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open database connection")
	}
	defer func(database *sql.DB) {
		err = database.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to gracefully close database connection")
		}
	}(database)

	connectedUsers := make(map[int64]*client.Client)
	var msg tgbotapi.MessageConfig

	updates := bot.GetUpdatesChan(upd)
	for update := range updates {
		if update.Message == nil {
			for _, connected := range connectedUsers {
				connected.CheckTimeout()
			}
			continue
		}

		var c *client.Client
		u := update.SentFrom()
		if conn, ok := connectedUsers[u.ID]; !ok {
			c, err = client.NewClient(update, database)
			connectedUsers[u.ID] = c
		} else {
			c = conn
		}

		if err == errors.New("user not initialised") {
			msg.Text = "Неполная информация о пользователе"
			if _, err = bot.Send(msg); err != nil {
				log.Fatal().Err(err).Msg("failed to send message")
			}

			continue
		}

		msg.Text = update.Message.Text
		if msg.Text == "/back" {
			c.State = 0
			msg.ReplyMarkup = internal.StartKeyboard
		} else {
			msg, err = c.HandleState(msg)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to perform bot action")
			}
		}

		if _, err = bot.Send(msg); err != nil {
			log.Fatal().Err(err).Msg("failed to send message")
		}
	}
}
