package telegram

import (
	"github.com/postgres-ci/notifier/src/common"

	"fmt"
)

const (
	successMessage = `
	Build success
		project: %s
		branch: %s
		commit: %s
		message: %s
	`
	errorMessage = `
	Build failed
		project: %s
		branch: %s
		commit: %s
		message: %s
	`
)

func (b *bot) Send(build common.Build) error {

	for _, recipient := range build.SendTo {

		if recipient.Method == Method {

			var message string

			if build.Status == "success" {

				message = fmt.Sprintf(successMessage,
					build.ProjectName,
					build.Branch,
					build.CommitSHA,
					build.CommitMessage,
				)

			} else {
				message = fmt.Sprintf(errorMessage,
					build.ProjectName,
					build.Branch,
					build.CommitSHA,
					build.CommitMessage,
				)
			}

			b.SendMessage(
				&user{TelegramID: recipient.IntID},
				message,
				nil,
			)
		}
	}

	return nil
}
