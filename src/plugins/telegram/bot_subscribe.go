package telegram

import (
	"github.com/tucnak/telebot"
)

func (b *bot) subscribe(message telebot.Message) {

	currentUser, err := b.currentUser(message.Sender.Username)

	if err != nil {

		b.SendMessage(message.Chat, err.Error(), nil)

		return
	}

	if currentUser.TelegramID != 0 {

		b.SendMessage(message.Chat, AlreadySubscribed, nil)

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
