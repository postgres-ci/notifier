package telegram

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/postgres-ci/notifier/src/common"
	"github.com/tucnak/telebot"

	"database/sql"
	"fmt"
	"strings"
	"time"
)

func New(token string, connect *sqlx.DB) (*bot, error) {

	telebot, err := telebot.NewBot(token)

	if err != nil {

		log.Warnf("Telegram Bot: %v", err)

		return nil, err
	}

	bot := &bot{
		token:   token,
		telebot: telebot,
		connect: connect,
	}

	go bot.listen()

	return bot, nil
}

type bot struct {
	token   string
	telebot *telebot.Bot
	connect *sqlx.DB
}

const (
	Method = "telegram"
)

func (b *bot) Send(build common.Build) error {

	log.Debugf("Telegram send notify, build [%d]", build.ID)

	for _, recipient := range build.SendTo {

		if recipient.Method != Method {

			log.Debugf("Telegram. Skip %s", recipient.Name)
		} else {

			b.telebot.SendMessage(&user{TelegramID: recipient.IntID}, build.CommitMessage, nil)
		}
	}

	return nil
}

func (b *bot) listen() {

	messages := make(chan telebot.Message)
	b.telebot.Listen(messages, time.Second)

	for message := range messages {

		switch strings.ToLower(message.Text) {

		case "/subscribe":
			b.subscribe(message)
		case "/unsubscribe":
			b.unsubscribe(message)
		default:

			if err := b.telebot.SendMessage(message.Chat, UsageMessage, nil); err != nil {

				log.Warnf("Telegram Bot can not send message: %v", err)
			}
		}
	}
}

type user struct {
	ID         int32 `db:"user_id"`
	TelegramID int64 `db:"telegram_id"`
}

func (u *user) Destination() string {

	return fmt.Sprint(u.TelegramID)
}

func (b *bot) subscribe(message telebot.Message) {

	var user user

	err := b.connect.Get(&user, `SELECT user_id, telegram_id FROM notification.find_user_by_telegram_username($1)`, message.Sender.Username)

	if err != nil {

		if err == sql.ErrNoRows {

			if err := b.telebot.SendMessage(message.Chat, "Not found", nil); err != nil {

				log.Warnf("Telegram Bot can not send message: %v", err)
			}
		} else {

			if err := b.telebot.SendMessage(message.Chat, err.Error(), nil); err != nil {

				log.Warnf("Telegram Bot can not send message: %v", err)
			}
		}

		return
	}

	_, err = b.connect.Exec(`SELECT notification.bind_with_telegram(
			$1,
			$2,
			$3
		)`,

		user.ID,
		message.Sender.Username,
		message.Sender.ID,
	)

	user.TelegramID = int64(message.Sender.ID)

	if err == nil {

		b.telebot.SendMessage(&user, "Ok, subscribed", nil)
	}
}

func (b *bot) unsubscribe(message telebot.Message) {

	var user user

	err := b.connect.Get(&user, `SELECT user_id, telegram_id FROM notification.find_user_by_telegram_username($1)`, message.Sender.Username)

	if err != nil {

		if err == sql.ErrNoRows {

			if err := b.telebot.SendMessage(message.Chat, "Not found", nil); err != nil {

				log.Warnf("Telegram Bot can not send message: %v", err)
			}
		} else {

			if err := b.telebot.SendMessage(message.Chat, err.Error(), nil); err != nil {

				log.Warnf("Telegram Bot can not send message: %v", err)
			}
		}

		return
	}

	_, err = b.connect.Exec(`SELECT notification.bind_with_telegram(
			$1,
			$2,
			$3
		)`,

		user.ID,
		message.Sender.Username,
		0,
	)

	user.TelegramID = int64(message.Sender.ID)

	if err == nil {

		b.telebot.SendMessage(&user, "Ok, unsubscribed", nil)
	}
}

const (
	UsageMessage = `
	Usage:
		/subscribe
		/unsubscribe
	`
)
