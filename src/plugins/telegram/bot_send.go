package telegram

import (
	"github.com/flosch/pongo2"
	"github.com/postgres-ci/notifier/src/common"
	"github.com/tucnak/telebot"
)

var template = pongo2.Must(pongo2.FromString(`{% autoescape off %}
{% if build.Status == "success" %}
<b>Build #{{ build.ID }} has passed</b>
{% else %}
<b>Build #{{ build.ID }} has failed</b>
{% endif %}
<b>{{ build.ProjectName }}</b> (<i>{{ build.Branch }}</i>)
{% if build.Error %}<pre>{{ build.Error }}</pre>{% endif %}
{{ build.CommitterName }} ({{ build.CommitterEmail }}) at {{ build.CommittedAt | time:"Mon, 02 Jan 2006 15:04:05 -0700" }}
{{ build.CommitMessage }}

sha: {% if APP_ADDRESS %}<a href="{{APP_ADDRESS}}/project-{{ build.ProjectID }}/build-{{ build.ID }}/">{{ build.CommitSHA }}</a>{% else %}{{ build.CommitSHA }}{% endif %}
{% endautoescape %}`))

func (b *bot) Send(build common.Build) error {

	message, err := template.Execute(pongo2.Context{
		"APP_ADDRESS": b.config.AppAddress,
		"build":       build,
	})

	if err != nil {

		return err
	}

	for _, recipient := range build.SendTo {

		if recipient.Method == Method && recipient.IntID != 0 {

			b.SendMessage(&user{TelegramID: recipient.IntID}, message, &telebot.SendOptions{ParseMode: telebot.ModeHTML, DisableWebPagePreview: true})
		}
	}

	return nil
}
