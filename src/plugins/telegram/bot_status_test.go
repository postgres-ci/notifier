package telegram

import (
	"github.com/erikstmartin/go-testdb"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/tucnak/telebot"

	"database/sql/driver"
	"fmt"
	"testing"
)

func TestStatusNotFound(t *testing.T) {

	defer testdb.Reset()

	var sendedMessage string

	bot := bot{
		connect: sqlx.MustOpen("testdb", ""),
		telebot: &testSender{
			send: func(recipient telebot.Recipient, message string, _ *telebot.SendOptions) error {

				sendedMessage = message

				return nil
			},
		},
	}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {

		return testdb.RowsFromSlice(
			[]string{"user_id", "telegram_id"},
			[][]driver.Value{},
		), nil
	})

	bot.status(telebot.Message{
		Sender: telebot.User{
			Username: "test",
		},
	})

	if assert.NotEmpty(t, sendedMessage) {

		assert.Equal(t, fmt.Sprintf(UsernameNotFound, "test"), sendedMessage)
	}
}

func TestStatusNotSubscribed(t *testing.T) {

	defer testdb.Reset()

	var sendedMessage string

	bot := bot{
		connect: sqlx.MustOpen("testdb", ""),
		telebot: &testSender{
			send: func(recipient telebot.Recipient, message string, _ *telebot.SendOptions) error {

				sendedMessage = message

				return nil
			},
		},
	}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {

		return testdb.RowsFromSlice(
			[]string{"user_id", "telegram_id"},
			[][]driver.Value{
				{
					"1",
					"0",
				},
			},
		), nil
	})

	bot.status(telebot.Message{
		Sender: telebot.User{
			Username: "test",
		},
	})

	if assert.NotEmpty(t, sendedMessage) {

		assert.Equal(t, NotSubscribed, sendedMessage)
	}
}

func TestStatusAlreadySubscribed(t *testing.T) {

	defer testdb.Reset()

	var sendedMessage string

	bot := bot{
		connect: sqlx.MustOpen("testdb", ""),
		telebot: &testSender{
			send: func(recipient telebot.Recipient, message string, _ *telebot.SendOptions) error {

				sendedMessage = message

				return nil
			},
		},
	}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {

		return testdb.RowsFromSlice(
			[]string{"user_id", "telegram_id"},
			[][]driver.Value{
				{
					"1",
					"42",
				},
			},
		), nil
	})

	bot.status(telebot.Message{
		Sender: telebot.User{
			Username: "test",
		},
	})

	if assert.NotEmpty(t, sendedMessage) {

		assert.Equal(t, AlreadySubscribed, sendedMessage)
	}
}
