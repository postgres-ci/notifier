package telegram

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/tucnak/telebot"

	"strings"
	"time"
)

const (
	Method = "telegram"
)

func New(token string, connect *sqlx.DB) (*bot, error) {

	if token == "" {

		return &bot{
			connect: connect,
		}, nil
	}

	telegramBot, err := telebot.NewBot(token)

	if err != nil {

		log.Warnf("Telegram Bot: %v", err)

		return nil, err
	}

	bot := &bot{
		telebot:  telegramBot,
		connect:  connect,
		messages: make(chan telebot.Message),
	}

	telegramBot.Listen(bot.messages, time.Second)

	for i := 0; i < 5; i++ {

		go bot.listen()
	}

	return bot, nil
}

type bot struct {
	connect *sqlx.DB
	telebot interface {
		SendMessage(telebot.Recipient, string, *telebot.SendOptions) error
	}
	messages chan telebot.Message
}

func (b *bot) listen() {

	for message := range b.messages {

		switch strings.ToLower(message.Text) {
		case "/status":
			b.status(message)
		case "/subscribe":
			b.subscribe(message)
		case "/unsubscribe":
			b.unsubscribe(message)
		default:
			b.SendMessage(message.Chat, UsageMessage, &telebot.SendOptions{ParseMode: telebot.ModeHTML})
		}
	}
}

func (b *bot) SendMessage(recipient telebot.Recipient, message string, options *telebot.SendOptions) {

	if b.telebot == nil {

		return
	}

	log.Debugf("telegram send: %s to %s", message, recipient.Destination())

	if err := b.telebot.SendMessage(recipient, message, options); err != nil {

		log.Errorf("Error when sending telegram message: %v", err)
	}
}

const (
	UsageMessage = `
Hello, i'm a <b>Postgres-CI</b> notifier

commands:
	/status
	/subscribe
	/unsubscribe
	`
)
