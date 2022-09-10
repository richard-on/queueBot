package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/richard-on/QueueBot/pkg/queueBot"
	"github.com/richard-on/QueueBot/pkg/queueBot/db"
)

/*func (c *Client) handleGroupSelect(msg tgbotapi.MessageConfig) ([]queueBot.QueueInfo, tgbotapi.MessageConfig, error) {

	var queues []queueBot.QueueInfo
	for _, subject := range subjects {
		if update.Message.Text == subject.SubjectName {
			c.State = queueBot.QueueSelect
			queues, err = db.GetQueues(update.Message.Text)
			if err != nil {
				return nil, msg, err
			}

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Выбран предмет \""+subject.SubjectName+"\"")
			msg.ReplyMarkup = queueBot.CreateQueueReplyKeyboard(queues)
		}
	}

	return queues, msg, nil
}*/

func (c *Client) handleSubjectSelect(msg tgbotapi.MessageConfig) ([]queueBot.QueueInfo, tgbotapi.MessageConfig, error) {
	subjects, err := db.GetSubjects(user)
	if err != nil {
		return nil, msg, err
	}

	var queues []queueBot.QueueInfo
	for _, subject := range subjects {
		if msg.Text == subject.SubjectName {
			c.State = queueBot.QueueSelect
			queues, err = db.GetQueues(msg.Text)
			if err != nil {
				return nil, msg, err
			}

			msg = tgbotapi.NewMessage(c.User.ChatID, "Выбран предмет \""+subject.SubjectName+"\"")
			msg.ReplyMarkup = queueBot.CreateQueueReplyKeyboard(queues)
		}
	}

	return queues, msg, nil
}

func (c *Client) handleQueueSelect(queues []queueBot.QueueInfo, msg tgbotapi.MessageConfig) (queueBot.QueueInfo, tgbotapi.MessageConfig, error) {
	for _, queue := range queues {
		if msg.Text == queue.Name {
			c.State = queueBot.QueueAction
			msg = tgbotapi.NewMessage(c.User.ChatID, "Выбрана очередь \""+queue.Name+"\"")
			msg.ReplyMarkup = queueBot.QueueActionKeyboard
			return queue, msg, nil
		}
	}

	return queueBot.QueueInfo{}, msg, nil
}

func (c *Client) handleActionSelect(queue queueBot.QueueInfo, msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	switch msg.Text {
	case "Войти в очередь":
		err := db.JoinQueue(queue.SubjectId, queue.QueueId, c.User.ChatID)
		if err != nil {
			msg = tgbotapi.NewMessage(c.User.ChatID, "Error entering the queue: "+err.Error())
			return msg, err
		}
		msg = tgbotapi.NewMessage(c.User.ChatID, "Entered queue")

	case "Выйти из очереди":
		err := db.LeaveQueue(queue.SubjectId, queue.QueueId, c.User.ChatID)
		if err != nil {
			msg = tgbotapi.NewMessage(c.User.ChatID, "Error leaving the queue: "+err.Error())
			return msg, err
		}
		msg = tgbotapi.NewMessage(c.User.ChatID, "Left queue")

	case "Показать очередь":
		data, err := db.PrintQueue(queue.QueueId, c.User.ChatID)
		if err != nil {
			msg = tgbotapi.NewMessage(c.User.ChatID, "Error printing the queue: "+err.Error())
			return msg, err
		}
		msg = tgbotapi.NewMessage(c.User.ChatID, data)

	}

	c.State = queueBot.QueueAction

	return msg, nil
}
