package bot

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/richard-on/QueueBot/pkg/queueBot"
	"github.com/richard-on/QueueBot/pkg/queueBot/db"
	"log"
	"os"
	"time"
)

type Client struct {
	User     *queueBot.User
	State    uint
	IsActive bool
	LastConn time.Time
}

func NewClient(update tgbotapi.Update) *Client {
	var client Client
	tgUser := update.SentFrom()

	user, err := db.CheckUserData(tgUser)
	if err != nil {
		client.User = &queueBot.User{
			UserID:      tgUser.ID,
			ChatID:      update.Message.Chat.ID,
			TgUsername:  tgUser.UserName,
			GroupID:     0,
			SubgroupID:  sql.NullInt64{},
			TgFirstName: sql.NullString{String: tgUser.FirstName},
			TgLastName:  sql.NullString{String: tgUser.LastName},
			FirstName:   sql.NullString{},
			LastName:    sql.NullString{},
		}
		//client.Message = update.Message.Text
		client.State = queueBot.Initial
		client.IsActive = true
		client.LastConn = time.Now()

		db.AddUser(client.User)

		return &client
	}

	client.User = &user
	client.User.ChatID = client.User.UserID
	client.IsActive = true
	client.LastConn = time.Now()

	return &client
}

func (c *Client) CheckTimeout() {
	if time.Now().Sub(c.LastConn) > time.Minute*1 {
		c.IsActive = false
	}
}

func Start() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	var connectedUsers []*Client
	var msg tgbotapi.MessageConfig
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			for _, u := range connectedUsers {
				u.CheckTimeout()
			}
			continue
		}

		client := NewClient(update)
		connectedUsers = append(connectedUsers, client)

		msg.Text = update.Message.Text
		msg, err = client.HandleState(msg)
		if err != nil {
			log.Fatal(err)
		}

		if _, err = bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
