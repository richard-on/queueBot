package client

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/richard-on/QueueBot/internal"
	"github.com/richard-on/QueueBot/internal/bot/model"
	"github.com/richard-on/QueueBot/internal/db"
	"time"
)

type BotState uint

const (
	Initial = iota
	GroupSelect
	SubjectSelect
	QueueSelect
	QueueAction
	AdminMode
)

type Client struct {
	User     *model.User
	Group    *db.Group
	Subject  []db.Subject
	Queue    []db.Queue
	CurQueue db.Queue
	State    BotState
	Db       db.QueueDB
	IsActive bool
	LastConn time.Time
}

func NewClient(update tgbotapi.Update, database *sql.DB) (*Client, error) {
	var client Client
	client.Db = db.NewQueueDB(database)
	tgUser := update.SentFrom()

	user, err := client.Db.GetUser(tgUser.ID)
	if err != nil {
		client.User = &model.User{
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

		err = client.Db.AddUser(client.User)
		if err != nil {
			return &client, err
		}

		client.State = Initial
		client.IsActive = true
		client.LastConn = time.Now()

		return &client, err
	}

	client.User = user
	client.User.ChatID = client.User.UserID
	client.IsActive = true
	client.LastConn = time.Now()

	return &client, nil
}

func (c *Client) CheckTimeout() {
	if time.Now().Sub(c.LastConn) > time.Minute*1 {
		c.IsActive = false
	}
}

func (c *Client) HandleState(msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	var err error
	switch c.State {
	case Initial:
		msg, err = c.handleCommand(msg)

	case GroupSelect:
		c.Subject, msg, err = c.handleGroupSelect(msg)

	case SubjectSelect:
		c.Queue, msg, err = c.handleSubjectSelect(msg)

	case QueueSelect:
		c.CurQueue, msg, err = c.handleQueueSelect(c.Queue, msg)

	case QueueAction:
		msg, err = c.handleActionSelect(c.CurQueue, msg)

	case AdminMode:
		/*argsString := update.Message.CommandArguments()
		if argsString == "" {
			msg = tgbotapi.NewMessage(c.User.ChatID, "Admin usage:\nadd_subject\nrm_subject\nadd_queue\nrm_queue")
		} else {
			queueBot.BotState = queueBot.AdminMode
			args := strings.Split(argsString, " ")
			msg, err = handleAdmin(update, msg, args)
			if err != nil {
				msg = tgbotapi.NewMessage(c.User.ChatID, "Error: "+err.Error())
			}
		}*/

	default:
		c.State = Initial
		msg = tgbotapi.NewMessage(c.User.ChatID, "Unknown text")
		msg.ReplyMarkup = internal.StartKeyboard
	}

	if err != nil {
		msg = tgbotapi.NewMessage(c.User.ChatID, "Error: "+err.Error())
	}

	return msg, nil
}

func (c *Client) handleCommand(msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	switch msg.Text {
	case "/start":
		greet, err := c.CreateGreeting(c.User)
		if err != nil {
			msg = tgbotapi.NewMessage(c.User.ChatID, err.Error())
		}
		msg = tgbotapi.NewMessage(c.User.ChatID, greet)

		groupName, err := c.Db.GetGroupName(c.User)
		if err != nil {
			return tgbotapi.MessageConfig{}, err
		}

		subGroupName, err := c.Db.GetSubGroupName(c.User)
		if err != nil {
			return tgbotapi.MessageConfig{}, err
		}

		msg.ReplyMarkup = internal.CreateGroupReplyKeyboard(groupName, subGroupName)
		c.State = GroupSelect

	case "/groups":
		if c.User.GroupID == 0 && !c.User.SubgroupID.Valid {
			msg = tgbotapi.NewMessage(c.User.ChatID, "Недостаточно информации")
			c.State = Initial
		}

		/*internal.BotState = internal.SubjectSelect
		subjects, err := db2.GetSubjects(user)
		if err != nil {
			return msg, err
		}
		var sb strings.Builder
		sb.WriteString("Choose a subject. Available subjects are:\n")
		for _, subject := range subjects {
			sb.WriteString("• " + subject.SubjectName + "\n")
		}*/

		c.State = GroupSelect

		groupName, err := c.Db.GetGroupName(c.User)
		if err != nil {
			return tgbotapi.MessageConfig{}, err
		}

		subGroupName, err := c.Db.GetSubGroupName(c.User)
		if err != nil {
			return tgbotapi.MessageConfig{}, err
		}

		msg = tgbotapi.NewMessage(c.User.ChatID, "Выберите группу")
		msg.ReplyMarkup = internal.CreateGroupReplyKeyboard(groupName, subGroupName)

	case "/back":
		c.State -= 1
		msg = tgbotapi.NewMessage(c.User.ChatID, "Назад")

	case "/admin":
		/*if db.CheckAdmin(c.User.ChatID, c.User.TgUsername) {
			c.State = AdminMode
			msg = tgbotapi.NewMessage(c.User.ChatID, "Auth success. Entered admin mode")
		} else {
			c.State = Initial
			msg = tgbotapi.NewMessage(c.User.ChatID, "You are not an admin")
		}*/

	default:
		c.State = Initial
		msg = tgbotapi.NewMessage(c.User.ChatID, "Unknown command")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("/start"),
			),
		)
	}

	return msg, nil
}
