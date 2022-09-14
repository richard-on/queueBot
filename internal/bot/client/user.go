package client

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/richard-on/QueueBot/internal"
	"github.com/richard-on/QueueBot/internal/db"
)

func (c *Client) handleGroupSelect(msg tgbotapi.MessageConfig) ([]db.Subject, tgbotapi.MessageConfig, error) {
	groupName, err := c.Db.GetGroupName(c.User)
	if err != nil {
		return nil, tgbotapi.MessageConfig{}, err
	}

	subGroupName, err := c.Db.GetSubGroupName(c.User)
	if err != nil {
		return nil, tgbotapi.MessageConfig{}, err
	}

	var group string
	var subjects []db.Subject
	if msg.Text == groupName {
		subjects, err = c.Db.GetSubjectList(c.User.GroupID)
		if err != nil {
			return nil, tgbotapi.MessageConfig{}, err
		}
		group = groupName

	} else if msg.Text == subGroupName {
		subjects, err = c.Db.GetSubjectList(c.User.SubgroupID)
		if err != nil {
			return nil, tgbotapi.MessageConfig{}, err
		}
		group = subGroupName

	} else {
		msg = tgbotapi.NewMessage(c.User.ChatID, "Неизвестная группа")
		return subjects, msg, nil
	}

	c.State = SubjectSelect
	msg = tgbotapi.NewMessage(c.User.ChatID, "Выбрана группа \""+group+"\"")
	msg.ReplyMarkup = internal.CreateSubjectReplyKeyboard(subjects)

	return subjects, msg, nil
}

func (c *Client) handleSubjectSelect(msg tgbotapi.MessageConfig) ([]db.Queue, tgbotapi.MessageConfig, error) {
	var queues []db.Queue
	var err error
	for _, subject := range c.Subject {
		if msg.Text == subject.SubjectName {
			queues, err = c.Db.GetQueueList(subject)
			if err != nil {
				return nil, msg, err
			}

			msg = tgbotapi.NewMessage(c.User.ChatID, "Выбран предмет \""+subject.SubjectName+"\"")
			msg.ReplyMarkup = internal.CreateQueueReplyKeyboard(queues)
		}
	}

	c.State = QueueSelect

	return queues, msg, nil
}

func (c *Client) handleQueueSelect(queues []db.Queue, msg tgbotapi.MessageConfig) (db.Queue, tgbotapi.MessageConfig, error) {
	for _, queue := range queues {
		if msg.Text == queue.Name {

			msg = tgbotapi.NewMessage(c.User.ChatID, "Выбрана очередь \""+queue.Name+"\"")
			msg.ReplyMarkup = internal.QueueActionKeyboard

			c.State = QueueAction

			return queue, msg, nil
		}
	}

	return db.Queue{}, msg, errors.New("no such queue")
}

func (c *Client) handleActionSelect(queue db.Queue, msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	switch msg.Text {
	case "Войти в очередь":
		err := c.Db.JoinQueue(c.User, &queue)
		if err != nil {
			msg = tgbotapi.NewMessage(c.User.ChatID, "Error entering the queue: "+err.Error())
			return msg, err
		}
		msg = tgbotapi.NewMessage(c.User.ChatID, "Entered queue")

	case "Выйти из очереди":
		err := c.Db.LeaveQueue(c.User, &queue)
		if err != nil {
			msg = tgbotapi.NewMessage(c.User.ChatID, "Error leaving the queue: "+err.Error())
			return msg, err
		}
		msg = tgbotapi.NewMessage(c.User.ChatID, "Left queue")

	case "Показать очередь":
		data, err := c.Db.PrintQueue(queue.ID)
		if err != nil {
			msg = tgbotapi.NewMessage(c.User.ChatID, "Error printing the queue: "+err.Error())
			return msg, err
		}
		msg = tgbotapi.NewMessage(c.User.ChatID, data)

	default:
		msg = tgbotapi.NewMessage(c.User.ChatID, "Неподдерживаемое действие")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("/groups"),
				tgbotapi.NewKeyboardButton("/back"),
			),
		)
	}

	c.State = QueueAction

	return msg, nil
}
