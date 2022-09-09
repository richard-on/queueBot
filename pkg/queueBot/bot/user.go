package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/richard-on/QueueBot/pkg/queueBot"
	"github.com/richard-on/QueueBot/pkg/queueBot/db"
)

func handleSubjectSelect(update tgbotapi.Update, msg tgbotapi.MessageConfig) ([]queueBot.QueueInfo, tgbotapi.MessageConfig, error) {
	subjects, err := db.GetSubjects(user)
	if err != nil {
		return nil, msg, err
	}

	var queues []queueBot.QueueInfo
	for _, subject := range subjects {
		if update.Message.Text == subject.SubjectName {
			queueBot.BotState = queueBot.QueueSelect
			queues, err = db.GetQueues(update.Message.Text)
			if err != nil {
				return nil, msg, err
			}
			text := subject.SubjectName + " selected. Now select queue"
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
			msg.ReplyMarkup = queueBot.CreateQueueReplyKeyboard(queues)
		}
	}

	return queues, msg, nil
}

func handleQueueSelect(queues []queueBot.QueueInfo, update tgbotapi.Update, msg tgbotapi.MessageConfig) (queueBot.QueueInfo, tgbotapi.MessageConfig, error) {
	for _, queue := range queues {
		if update.Message.Text == queue.Name {
			queueBot.BotState = queueBot.QueueAction
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Chosen queue "+queue.Name)
			msg.ReplyMarkup = queueBot.QueueActionKeyboard
			return queue, msg, nil
		}
	}

	return queueBot.QueueInfo{}, msg, nil
}

func handleActionSelect(queue queueBot.QueueInfo, update tgbotapi.Update, msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	switch update.Message.Text {
	case "Войти в очередь":
		err := db.JoinQueue(queue.SubjectId, queue.QueueId, update.Message.Chat.ID)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error entering the queue: "+err.Error())
			return msg, err
		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Entered queue")

	case "Выйти из очереди":
		err := db.LeaveQueue(queue.SubjectId, queue.QueueId, update.Message.Chat.ID)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error leaving the queue: "+err.Error())
			return msg, err
		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Left queue")

	case "Показать очередь":
		data, err := db.PrintQueue(queue.QueueId, update.Message.Chat.ID)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error printing the queue: "+err.Error())
			return msg, err
		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, data)

	}

	queueBot.BotState = queueBot.QueueAction

	return msg, nil
}
