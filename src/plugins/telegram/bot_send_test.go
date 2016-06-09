package telegram

import (
	log "github.com/Sirupsen/logrus"
	"github.com/postgres-ci/notifier/src/common"
	"github.com/stretchr/testify/assert"
	"github.com/tucnak/telebot"

	"bytes"
	"testing"
	"time"
)

type testSender struct {
	send func(recipient telebot.Recipient, message string, _ *telebot.SendOptions) error
}

func (t *testSender) SendMessage(recipient telebot.Recipient, message string, options *telebot.SendOptions) error {
	return t.send(recipient, message, options)
}

func TestSendSuccess(t *testing.T) {

	log.SetOutput(&bytes.Buffer{})

	var (
		telegramID    string
		sendedMessage string
	)

	sender := &testSender{
		send: func(recipient telebot.Recipient, message string, _ *telebot.SendOptions) error {

			telegramID = recipient.Destination()
			sendedMessage = message

			return nil
		},
	}

	bot := bot{
		telebot: sender,
	}

	now := time.Now()

	bot.Send(common.Build{
		Status:         "success",
		ProjectName:    "Postgres-CI test project",
		Branch:         "testbranch",
		CommitSHA:      "6ba05bc0e064d5ad9b6044199edb99c1aca7f024",
		CommitMessage:  "Git commit message",
		CommitterName:  "Elephant Sam",
		CommitterEmail: "samelephant82@gmail.com",
		CommittedAt:    now,
		SendTo: []common.Recipient{
			{
				Method: "telegram",
				IntID:  42,
			},
		},
	})

	if assert.Equal(t, "42", telegramID) {

		for _, contains := range []string{
			"success",
			"Postgres-CI test project",
			"testbranch",
			"6ba05bc0e064d5ad9b6044199edb99c1aca7f024",
			"Git commit message",
			now.Format("Mon, 02 Jan 2006 15:04:05 -0700"),
			"Elephant Sam",
			"samelephant82@gmail.com",
		} {

			assert.Contains(t, sendedMessage, contains)
		}

	}

	t.Log(sendedMessage)
}

func TestSendFailedBuild(t *testing.T) {

	buildError := `
ERROR:  unrecognized GET DIAGNOSTICS item at or near "table_name"
	`

	log.SetOutput(&bytes.Buffer{})

	var (
		telegramID    string
		sendedMessage string
	)

	sender := &testSender{
		send: func(recipient telebot.Recipient, message string, _ *telebot.SendOptions) error {

			telegramID = recipient.Destination()
			sendedMessage = message

			return nil
		},
	}

	bot := bot{
		telebot: sender,
	}

	now := time.Now()

	bot.Send(common.Build{
		Status:         "failed",
		Error:          buildError,
		ProjectName:    "Postgres-CI test project",
		Branch:         "testbranch",
		CommitSHA:      "6ba05bc0e064d5ad9b6044199edb99c1aca7f024",
		CommitMessage:  "Git commit message",
		CommitterName:  "Elephant Sam",
		CommitterEmail: "samelephant82@gmail.com",
		CommittedAt:    now,
		SendTo: []common.Recipient{
			{
				Method: "telegram",
				IntID:  42,
			},
		},
	})

	if assert.Equal(t, "42", telegramID) {

		for _, contains := range []string{
			"failed",
			"Postgres-CI test project",
			"testbranch",
			"6ba05bc0e064d5ad9b6044199edb99c1aca7f024",
			"Git commit message",
			now.Format("Mon, 02 Jan 2006 15:04:05 -0700"),
			"Elephant Sam",
			"samelephant82@gmail.com",
		} {

			assert.Contains(t, sendedMessage, contains)
		}

	}

	t.Log(sendedMessage)
}
