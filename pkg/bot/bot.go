package bot

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/richard-on/queueBot/config"
	"github.com/richard-on/queueBot/internal/bot/client"
	"github.com/richard-on/queueBot/pkg/logger"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"
)

// healthHandler is a handler for healthcheck
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "queueBot is up and running!")
}

// healthcheck creates an endpoint with bot status info
func healthcheck(log logger.Logger) {
	http.HandleFunc("/queue-bot/health", healthHandler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err, "cannot establish healthcheck")
	}
}

// Run creates and runs telegram bot
func Run() {
	// Set up logger
	log := logger.NewLogger(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC1123},
		zerolog.TraceLevel,
		"queueBot-bot")

	// Connect to telegram bot with token
	bot, err := tgbotapi.NewBotAPI(config.TgToken)
	if err != nil {
		log.Fatal(err, "failed to connect to Telegram bot")
	}
	log.Info("Authorized on account " + bot.Self.UserName)

	// Set bot modes and create update config
	bot.Debug = false
	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 60

	// Connect to a database
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
	log.Info("Connected to database")

	// Set up healthcheck endpoint in a separate goroutine
	go healthcheck(log)

	// Create map of connected clients, where key is userID
	connectedUsers := make(map[int64]*client.Client, 5)
	var msg tgbotapi.MessageConfig

	// Create UpdatesChan
	updates := bot.GetUpdatesChan(upd)

	// Monitor connections
	go func(connPool *map[int64]*client.Client) {
		for {
			for _, connected := range *connPool {
				// Disconnect if client have not sent a message in a period
				if connected.CheckTimeout() {
					msg, err := connected.Disconnect()
					if err != nil {
						log.Error(err, "cannot disconnect user")
					}

					msgSend, err := bot.Send(msg)
					if err != nil {
						log.Error(err, "failed to send message")
					}
					log.TgSend(msgSend)

					delete(connectedUsers, connected.User.UserID)
					log.ClientDisconnect(*connected.User, connected.LastConn)
				}
			}
			time.Sleep(time.Second * 1)
		}

	}(&connectedUsers)

	for update := range updates {
		log.TgUpdate(update)

		// If no message has been received, ignore update
		if update.Message == nil {
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
			msgSend, err := bot.Send(msg)
			if err != nil {
				log.Fatal(err, "failed to send message")
				log.Info("sent message")
			}
			log.TgSend(msgSend)

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

		msgSend, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err, "failed to send message")
		}
		log.TgSend(msgSend)
	}
}
