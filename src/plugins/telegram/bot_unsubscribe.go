package telegram

import (
	"github.com/tucnak/telebot"
)

func (b *bot) unsubscribe(message telebot.Message) {

	currentUser, err := b.currentUser(message.Sender.Username)

	if err != nil {

		b.SendMessage(message.Chat, err.Error(), nil)

		return
	}

	if currentUser.TelegramID == 0 {

		b.SendMessage(message.Chat, NotSubscribed, nil)

		return
	}

	if currentUser.TelegramID != int64(message.Sender.ID) {

		b.SendMessage(message.Chat, "Sender ID is not matched", nil)

		return
	}

	_, err = b.connect.Exec(`SELECT notification.bind_with_telegram(
			$1,
			$2,
			$3
		)`,

		currentUser.ID,
		message.Sender.Username,
		0,
	)

	if err == nil {
		b.SendMessage(&user{TelegramID: int64(message.Sender.ID)}, "Ok, unsubscribed", nil)
	} else {
		b.SendMessage(message.Chat, "An error occurred please try again later", nil)
	}
}
