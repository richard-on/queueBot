package client

import (
	"fmt"
	"github.com/richard-on/queueBot/internal/bot/model"
)

var greeting = `Привет! Это бот для ведения очередей на сдачу лаб.

Прежде чем начать, проверим информацию и убедимся что ты состоишь в нужных группах:
Telegram: %v
Имя: %v
Фамилия: %v
Группа: %v
Кафедра: %v`

var NeedMoreInfo = `Недостаточно информации о пользователе.

Для работы бота необходима информация о группе и кафедре. Эти данные может предоставить только админ.

Пожалуйста, подождите, пока вас зарегистрируют.`

var NoSubjects = `Не найдено ни одного предмета, относящегося к этой группе.`

var NoQueues = `Очереди не найдены.`

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
