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

func TestSubscribeNotFound(t *testing.T) {

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

	bot.subscribe(telebot.Message{
		Sender: telebot.User{
			Username: "test",
		},
	})

	if assert.NotEmpty(t, sendedMessage) {

		assert.Equal(t, fmt.Sprintf(UsernameNotFound, "test"), sendedMessage)
	}
}

func TestSubscribeAlreadySubscribed(t *testing.T) {

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

	bot.subscribe(telebot.Message{
		Sender: telebot.User{
			Username: "test",
		},
	})

	if assert.NotEmpty(t, sendedMessage) {

		assert.Equal(t, AlreadySubscribed, sendedMessage)
	}
}

func TestSubscribeOk(t *testing.T) {

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

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (driver.Result, error) {

		if assert.Len(t, args, 3) {

			assert.Equal(t, int64(1), args[0].(int64))
			assert.Equal(t, "test", args[1].(string))
			assert.Equal(t, int64(42), args[2].(int64))
		}

		return testdb.NewResult(0, nil, 0, nil), nil
	})

	bot.subscribe(telebot.Message{
		Sender: telebot.User{
			ID:       42,
			Username: "test",
		},
	})

	if assert.NotEmpty(t, sendedMessage) {

		assert.Equal(t, "Ok, subscribed", sendedMessage)
	}
}
