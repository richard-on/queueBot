package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strings"
)

func handleCommand(update tgbotapi.Update, msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	switch update.Message.Command() {
	case "start":
		if err := CollectUserData(update.Message.Chat.ID, update.Message.Chat.UserName,
			update.Message.Chat.FirstName, update.Message.Chat.LastName); err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Database error, we can't identify you.")
			return msg, err
		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Hello! This is queue bot")
		msg.ReplyMarkup = commandKeyboard

	case "admin":
		if CheckAdmin(update.Message.Chat.ID, update.Message.Chat.UserName) {
			botState = adminMode
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Auth success. Entered admin mode")
		} else {
			botState = initial
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "You are not an admin")
		}

	case "subjects":
		botState = subjectSelect
		subjects, err := GetSubjects()
		if err != nil {
			return msg, err
		}
		var sb strings.Builder
		sb.WriteString("Choose a subject. Available subjects are:\n")
		for _, subject := range subjects {
			sb.WriteString("• " + subject.Name + "\n")
		}

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, sb.String())
		msg.ReplyMarkup = createSubjectReplyKeyboard(subjects)

	default:
		botState = initial
		msg.ReplyMarkup = commandKeyboard
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
	}

	return msg, nil
}

func handleBotState(update tgbotapi.Update, msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	var err error
	switch botState {
	case initial:
		msg, err = handleCommand(update, msg)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error: "+err.Error())
		}

	case subjectSelect:
		queues, msg, err = handleSubjectSelect(update, msg)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error: "+err.Error())
		}

	case queueSelect:
		queue, msg, err = handleQueueSelect(queues, update, msg)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error: "+err.Error())
		}

	case queueAction:
		msg, err = handleActionSelect(queue, update, msg)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error: "+err.Error())
		}

	case adminMode:
		argsString := update.Message.CommandArguments()
		if argsString == "" {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Admin usage:\nadd_subject\nrm_subject\nadd_queue\nrm_queue")
		} else {
			botState = adminMode
			args := strings.Split(argsString, " ")
			msg, err = handleAdmin(update, msg, args)
			if err != nil {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error: "+err.Error())
			}
		}

	default:
		botState = initial
		msg.ReplyMarkup = commandKeyboard
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown text")
	}

	return msg, nil
}

var queues []QueueInfo
var queue QueueInfo

func Bot() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	var msg tgbotapi.MessageConfig
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg, err = handleBotState(update, msg)
		if err != nil {
			log.Fatal(err)
		}

		if _, err = bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
