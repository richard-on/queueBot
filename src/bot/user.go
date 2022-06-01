package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleSubjectSelect(update tgbotapi.Update, msg tgbotapi.MessageConfig) ([]QueueInfo, tgbotapi.MessageConfig, error) {
	subjects, err := GetSubjects()
	if err != nil {
		return nil, msg, err
	}

	var queues []QueueInfo
	for _, subject := range subjects {
		if update.Message.Text == subject.Name || update.Message.Text == subject.Alias {
			botState = queueSelect
			queues, err = GetQueues(update.Message.Text)
			if err != nil {
				return nil, msg, err
			}
			text := subject.Name + " selected. Now select queue"
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
			msg.ReplyMarkup = createQueueReplyKeyboard(queues)
		}
	}

	return queues, msg, nil
}

func handleQueueSelect(queues []QueueInfo, update tgbotapi.Update, msg tgbotapi.MessageConfig) (QueueInfo, tgbotapi.MessageConfig, error) {
	for _, queue := range queues {
		if update.Message.Text == queue.Name {
			botState = queueAction
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Chosen queue "+queue.Name)
			msg.ReplyMarkup = queueActionKeyboard
			return queue, msg, nil
		}
	}

	return QueueInfo{}, msg, nil
}

func handleActionSelect(queue QueueInfo, update tgbotapi.Update, msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	switch update.Message.Text {
	case "Enter":
		err := JoinQueue(queue.SubjectId, queue.QueueId, update.Message.Chat.ID)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error entering the queue: "+err.Error())
			return msg, err
		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Entered queue")

	case "Leave":
		err := LeaveQueue(queue.SubjectId, queue.QueueId, update.Message.Chat.ID)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error leaving the queue: "+err.Error())
			return msg, err
		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Left queue")

	case "Print":
		data, err := PrintQueue(queue.QueueId, update.Message.Chat.ID)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error printing the queue: "+err.Error())
			return msg, err
		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, data)

	}

	return msg, nil
}
