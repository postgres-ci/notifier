package telegram

import (
	"github.com/tucnak/telebot"

	"database/sql"
	"fmt"
)

func (b *bot) subscribe(message telebot.Message) {

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

	if currentUser.TelegramID != 0 {

		b.SendMessage(message.Chat, "Already subscribed", nil)

		return
	}

	_, err = b.connect.Exec(`SELECT notification.bind_with_telegram(
			$1,
			$2,
			$3
		)`,

		currentUser.ID,
		message.Sender.Username,
		message.Sender.ID,
	)

	if err == nil {
		b.SendMessage(&user{TelegramID: int64(message.Sender.ID)}, "Ok, subscribed", nil)
	} else {
		b.SendMessage(message.Chat, "An error occurred please try again later", nil)
	}
}
