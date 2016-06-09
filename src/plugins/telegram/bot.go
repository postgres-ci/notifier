package telegram

import (
	"github.com/jmoiron/sqlx"
	"github.com/postgres-ci/notifier/src/common"
)

func New(token string, connect *sqlx.DB) *bot {
	return &bot{
		token:   token,
		connect: connect,
	}
}

type bot struct {
	token   string
	connect *sqlx.DB
}

func (b *bot) Send(common.Build) error {
	return nil
}
