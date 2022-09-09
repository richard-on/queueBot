package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/richard-on/QueueBot/pkg/queueBot"
	"github.com/richard-on/QueueBot/pkg/queueBot/db"
	"strings"
)

var queueSlice []queueBot.QueueInfo
var queue queueBot.QueueInfo
var user queueBot.User

func HandleState(update tgbotapi.Update, msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	var err error
	switch queueBot.BotState {
	case queueBot.Initial:
		msg, err = handleCommand(update, msg)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error: "+err.Error())
		}

	case queueBot.SubjectSelect:
		queueSlice, msg, err = handleSubjectSelect(update, msg)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error: "+err.Error())
		}

	case queueBot.QueueSelect:
		queue, msg, err = handleQueueSelect(queueSlice, update, msg)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error: "+err.Error())
		}

	case queueBot.QueueAction:
		msg, err = handleActionSelect(queue, update, msg)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error: "+err.Error())
		}

	case queueBot.AdminMode:
		argsString := update.Message.CommandArguments()
		if argsString == "" {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Admin usage:\nadd_subject\nrm_subject\nadd_queue\nrm_queue")
		} else {
			queueBot.BotState = queueBot.AdminMode
			args := strings.Split(argsString, " ")
			msg, err = handleAdmin(update, msg, args)
			if err != nil {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error: "+err.Error())
			}
		}

	default:
		queueBot.BotState = queueBot.Initial
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown text")
		msg.ReplyMarkup = queueBot.StartKeyboard
	}

	return msg, nil
}

func handleCommand(update tgbotapi.Update, msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	var err error
	switch update.Message.Text {
	case "/start":
		err = db.CollectUserData(update.Message.Chat.ID, update.Message.Chat.UserName,
			update.Message.Chat.FirstName, update.Message.Chat.LastName)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error, we can't identify you.")
			return msg, err
		}
		user, err = db.GetUserData(update.Message.Chat.ID, update.Message.Chat.UserName)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error, we can't identify you.")
			return msg, err
		}
		greet, err := createGreeting(user)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
		}
		msg.ReplyMarkup = queueBot.CreateGroupReplyKeyboard(user)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, greet)

	case "admin":
		if db.CheckAdmin(update.Message.Chat.ID, update.Message.Chat.UserName) {
			queueBot.BotState = queueBot.AdminMode
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Auth success. Entered admin mode")
		} else {
			queueBot.BotState = queueBot.Initial
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "You are not an admin")
		}

	case "/subjects":
		user, err = db.GetUserData(update.Message.Chat.ID, update.Message.Chat.UserName)
		queueBot.BotState = queueBot.SubjectSelect
		subjects, err := db.GetSubjects(user)
		if err != nil {
			return msg, err
		}
		var sb strings.Builder
		sb.WriteString("Choose a subject. Available subjects are:\n")
		for _, subject := range subjects {
			sb.WriteString("• " + subject.SubjectName + "\n")
		}

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, sb.String())
		msg.ReplyMarkup = queueBot.CreateSubjectReplyKeyboard(subjects)

	default:
		queueBot.BotState = queueBot.Initial
		msg.ReplyMarkup = queueBot.StartKeyboard
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
	}

	return msg, nil
}
