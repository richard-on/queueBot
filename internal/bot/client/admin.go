package client

/*func (c *Client) handleAdmin(update tgbotapi.Update, msg tgbotapi.MessageConfig, args []string) (tgbotapi.MessageConfig, error) {
	switch args[0] {
	case "add_subject":
		err := db.AddSubject(args[1], args[2], args[3])
		if err != nil {
			return msg, err
		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Added subject "+args[2])

	case "rm_subject":
		if err := db.RmSubject(args[1], args[2]); err != nil {
			return msg, err
		}

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Removed subject "+args[2])

	case "add_queue":
		if err := db.AddQueue(args[1], args[2], args[3]); err != nil {
			return msg, err
		}

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Added queue "+args[2])

	case "rm_queue":
		if err := db.RmQueue(args[1], args[2], args[3]); err != nil {
			return msg, err
		}

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Removed queue "+args[2])

	case "exit":
		c.State = Initial
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Exited admin mode")

	}

	return msg, nil
}*/
