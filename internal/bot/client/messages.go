package client

import (
	"fmt"
	"github.com/richard-on/QueueBot/internal/bot/model"
)

var greeting = `Привет! Это бот для ведения очередей на сдачу лаб.

Прежде чем начать, проверим информацию и убедимся что ты состоишь в нужных группах:
Telegram: %v
Имя: %v
Фамилия: %v
Группа: %v
Кафедра: %v`

func (c *Client) CreateGreeting(user *model.User) (string, error) {
	group, err := c.Db.GetGroupName(user)
	if err != nil {
		return "", err
	}

	subgroup, err := c.Db.GetSubGroupName(user)
	if err != nil {
		return "", err
	}

	/*user.GroupName = group
	user.SubGroupName = subgroup*/

	return fmt.Sprintf(greeting,
		"@"+user.TgUsername, user.FirstName, user.LastName, group, subgroup), nil
}
