package email

import (
	"github.com/postgres-ci/notifier/src/common"
)

func New(config common.Config) *email {
	return &email{}
}

type email struct {
}

func (e *email) Send(common.Build) error {
	return nil
}
