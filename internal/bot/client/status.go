package client

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func (c *Client) CheckTimeout() bool {
	if time.Since(c.LastConn) > time.Second*30 {
		c.IsActive = false

		return true
	}

	return false
}

func (c *Client) Disconnect() (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	if c.IsActive {
		return tgbotapi.MessageConfig{}, errors.New("user is active")
	}

	c.State = Initial

	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/start"),
			tgbotapi.NewKeyboardButton("/groups"),
		),
	)

	msg.Text = "Sent disconnect"
	msg.ChatID = c.User.ChatID

	return msg, nil
}
