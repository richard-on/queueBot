package bot

import (
	"fmt"
	"github.com/richard-on/QueueBot/pkg/queueBot"
	"github.com/richard-on/QueueBot/pkg/queueBot/db"
)

var greeting = `Привет! Это бот для ведения очередей на сдачу лаб.

Прежде чем начать, проверим информацию и убедимся что ты состоишь в нужных группах:
Telegram: %v
Имя: %v
Фамилия: %v
Группа: %v
Кафедра: %v`

func createGreeting(user queueBot.User) (string, error) {
	group, err := db.GetGroup(user)
	if err != nil {
		return "", err
	}

	subgroup, err := db.GetSubGroup(user)
	if err != nil {
		return "", err
	}

	user.GroupName = group
	user.SubGroupName = subgroup

	return fmt.Sprintf(greeting,
		"@"+user.TgUsername, user.FirstName.String, user.LastName.String, group, subgroup), nil
}
