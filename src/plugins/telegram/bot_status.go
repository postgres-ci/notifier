package telegram

import (
	"github.com/tucnak/telebot"

	"database/sql"
	"fmt"
)

func (b *bot) status(message telebot.Message) {

	var currentUser user

	err := b.connect.Get(&currentUser, `SELECT user_id, telegram_id FROM notification.find_user_by_telegram_username($1)`, message.Sender.Username)

	if err != nil {

		if err == sql.ErrNoRows {

			b.SendMessage(message.Chat, fmt.Sprintf("Username \"%s\" not found", message.Sender.Username), nil)

		} else {

			b.SendMessage(message.Chat, err.Error(), nil)
		}

		return
	}

	if currentUser.TelegramID == 0 {

		b.SendMessage(message.Chat, "You are not subscribed", nil)

		return
	}

	b.SendMessage(message.Chat, "You are already subscribed", nil)
}
