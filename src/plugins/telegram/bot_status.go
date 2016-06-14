package telegram

import (
	"github.com/tucnak/telebot"
)

func (b *bot) status(message telebot.Message) {

	currentUser, err := b.currentUser(message.Sender.Username)

	if err != nil {

		b.SendMessage(message.Chat, err.Error(), nil)

		return
	}

	if currentUser.TelegramID == 0 {

		b.SendMessage(message.Chat, NotSubscribed, nil)

		return
	}

	b.SendMessage(message.Chat, AlreadySubscribed, nil)
}
