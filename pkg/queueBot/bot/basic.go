package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/richard-on/QueueBot/pkg/queueBot"
	"github.com/richard-on/QueueBot/pkg/queueBot/db"
)

var queueSlice []queueBot.QueueInfo
var queue queueBot.QueueInfo
var user queueBot.User

func (c *Client) HandleState(msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	var err error
	switch c.State {
	case queueBot.Initial:
		msg, err = c.handleCommand(msg)
		if err != nil {
			msg = tgbotapi.NewMessage(c.User.ChatID, "Error: "+err.Error())
		}

	case queueBot.SubjectSelect:
		queueSlice, msg, err = c.handleSubjectSelect(msg)
		if err != nil {
			msg = tgbotapi.NewMessage(c.User.ChatID, "Error: "+err.Error())
		}

	case queueBot.QueueSelect:
		queue, msg, err = c.handleQueueSelect(queueSlice, msg)
		if err != nil {
			msg = tgbotapi.NewMessage(c.User.ChatID, "Error: "+err.Error())
		}

	case queueBot.QueueAction:
		msg, err = c.handleActionSelect(queue, msg)
		if err != nil {
			msg = tgbotapi.NewMessage(c.User.ChatID, "Error: "+err.Error())
		}

	case queueBot.AdminMode:
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
		queueBot.BotState = queueBot.Initial
		msg = tgbotapi.NewMessage(c.User.ChatID, "Unknown text")
		msg.ReplyMarkup = queueBot.StartKeyboard
	}

	return msg, nil
}

func (c *Client) handleCommand(msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	var err error
	switch msg.Text {
	case "/start":
		/*err = db.CollectUserData(c.User.ChatID, c.User.TgUsername,
			c.User.TgFirstName, c.User.TgFirstName)
		if err != nil {
			msg = tgbotapi.NewMessage(c.User.ChatID, "Error, we can't identify you.")
			return msg, err
		}*/
		user, err = db.GetUserData(c.User.ChatID, c.User.TgUsername)
		if err != nil {
			msg = tgbotapi.NewMessage(c.User.ChatID, "Error, we can't identify you.")
			return msg, err
		}
		greet, err := createGreeting(user)
		if err != nil {
			msg = tgbotapi.NewMessage(c.User.ChatID, err.Error())
		}
		msg = tgbotapi.NewMessage(c.User.ChatID, greet)
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("13 группа"),
				tgbotapi.NewKeyboardButton("МСС"),
			),
		)

	case "admin":
		if db.CheckAdmin(c.User.ChatID, c.User.TgUsername) {
			queueBot.BotState = queueBot.AdminMode
			msg = tgbotapi.NewMessage(c.User.ChatID, "Auth success. Entered admin mode")
		} else {
			queueBot.BotState = queueBot.Initial
			msg = tgbotapi.NewMessage(c.User.ChatID, "You are not an admin")
		}

	case "/subjects":
		user, err = db.GetUserData(c.User.ChatID, c.User.TgUsername)
		queueBot.BotState = queueBot.SubjectSelect
		subjects, err := db.GetSubjects(user)
		if err != nil {
			return msg, err
		}
		/*var sb strings.Builder
		sb.WriteString("Choose a subject. Available subjects are:\n")
		for _, subject := range subjects {
			sb.WriteString("• " + subject.SubjectName + "\n")
		}*/

		msg = tgbotapi.NewMessage(c.User.ChatID, "Выберите предмет")
		msg.ReplyMarkup = queueBot.CreateSubjectReplyKeyboard(subjects)

	case "/back":
		queueBot.BotState -= 1

		/*msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите предмет")
		msg.ReplyMarkup = queueBot.CreateSubjectReplyKeyboard(subjects)*/

	default:
		queueBot.BotState = queueBot.Initial
		msg = tgbotapi.NewMessage(c.User.ChatID, "Unknown command")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("/start"),
			),
		)
	}

	return msg, nil
}
