package email

import (
	log "github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	"github.com/go-gomail/gomail"
	"github.com/postgres-ci/notifier/src/common"

	"fmt"
)

const Method = "email"

var template = pongo2.Must(pongo2.FromString(emailBodyTpl))

func New(config common.Config) *email {

	return &email{
		config: config,
	}
}

type email struct {
	config common.Config
}

func (e *email) Send(build common.Build) error {

	if e.config.SMTP.Host == "" {

		return nil
	}

	smtp := gomail.NewDialer(e.config.SMTP.Host, int(e.config.SMTP.Port), e.config.SMTP.Username, e.config.SMTP.Password)

	connect, err := smtp.Dial()

	if err != nil {

		log.Warnf("SMTP: %v", err)

		return err
	}

	body, err := template.Execute(pongo2.Context{
		"APP_ADDRESS": e.config.AppAddress,
		"build":       build,
	})

	if err != nil {

		return err
	}

	message := gomail.NewMessage()

	for _, recipient := range build.SendTo {

		if recipient.Method == Method && recipient.TextID != "" {

			message.SetHeader("From", e.config.SMTP.Username)
			message.SetAddressHeader("To", recipient.TextID, recipient.Name)

			if build.Status == "success" {

				message.SetHeader("Subject", fmt.Sprintf("Build #%d has passed. %s", build.ID, build.ProjectName))

			} else {

				message.SetHeader("Subject", fmt.Sprintf("Build #%d has failed. %s", build.ID, build.ProjectName))
			}

			message.SetBody("text/html", body)

			if err := gomail.Send(connect, message); err != nil {

				log.Warnf("An error occurred when sending email: %v", err)
			}

			message.Reset()
		}
	}

	return connect.Close()
}
