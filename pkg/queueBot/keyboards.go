package queueBot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

/*var CommandKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("admin"),
		tgbotapi.NewKeyboardButton("subjects"),
	),
)*/

var QueueActionKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Enter"),
		tgbotapi.NewKeyboardButton("Leave"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Print"),
	),
)

func CreateGroupReplyKeyboard(user User) tgbotapi.ReplyKeyboardMarkup {
	GroupActionKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(user.GroupName),
			tgbotapi.NewKeyboardButton(user.SubGroupName),
		),
	)

	return GroupActionKeyboard
}

func CreateSubjectReplyKeyboard(data []Subjects) tgbotapi.ReplyKeyboardMarkup {
	var subjectRows [][]tgbotapi.KeyboardButton
	var subjectRow []tgbotapi.KeyboardButton
	var subjectButtons []tgbotapi.KeyboardButton

	for i := 0; i < len(data); i++ {
		subjectButton := tgbotapi.NewKeyboardButton(data[i].Name)
		subjectButtons = append(subjectButtons, subjectButton)
		if (i+1)%3 == 0 || i == len(data)-1 {
			subjectRow = tgbotapi.NewKeyboardButtonRow(subjectButtons...)
			subjectRows = append(subjectRows, subjectRow)
			subjectButtons = nil
		}
	}

	var subjectNewKeyboard = tgbotapi.NewReplyKeyboard(subjectRows...)

	return subjectNewKeyboard
}

func CreateQueueReplyKeyboard(data []QueueInfo) tgbotapi.ReplyKeyboardMarkup {
	var queueRows [][]tgbotapi.KeyboardButton
	var subjectRow []tgbotapi.KeyboardButton
	var subjectButtons []tgbotapi.KeyboardButton

	for i := 0; i < len(data); i++ {
		subjectButton := tgbotapi.NewKeyboardButton(data[i].Name)
		subjectButtons = append(subjectButtons, subjectButton)
		if (i+1)%3 == 0 || i == len(data)-1 {
			subjectRow = tgbotapi.NewKeyboardButtonRow(subjectButtons...)
			queueRows = append(queueRows, subjectRow)
			subjectButtons = nil
		}
	}

	var newKeyboard = tgbotapi.NewReplyKeyboard(queueRows...)

	return newKeyboard
}

//TODO: Think about better keyboard creation options
/*type Keyboard interface {
	Subjects | QueueInfo
}

func createReplyKeyboard[T any](data []T) tgbotapi.ReplyKeyboardMarkup {
	var button tgbotapi.KeyboardButton
	var buttonSlice []tgbotapi.KeyboardButton
	var row []tgbotapi.KeyboardButton
	var rows [][]tgbotapi.KeyboardButton

	for i := 0; i < len(data); i++ {
		dataConv := data[i]
		switch any(&dataConv).(type) {
		case *Subjects:
			button = tgbotapi.NewKeyboardButton(dataConv.(Subjects).Name)
		case *QueueInfo:
			button = tgbotapi.NewKeyboardButton(dataConv.(QueueInfo).Name)
		}
		buttonSlice = append(buttonSlice, button)
		if (i+1)%3 == 0 || i == len(data)-1 {
			row = tgbotapi.NewKeyboardButtonRow(buttonSlice...)
			rows = append(rows, row)
			buttonSlice = nil
		}
	}

	return tgbotapi.NewReplyKeyboard(rows...)
}*/
