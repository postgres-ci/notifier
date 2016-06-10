package telegram

import (
	"github.com/flosch/pongo2"
	"github.com/postgres-ci/notifier/src/common"
	"github.com/tucnak/telebot"
)

var template = pongo2.Must(pongo2.FromString(`{% autoescape off %}
{% if build.Status == "failed" %}
<b>Build failed</b>
{% else %}
<b>Build success</b>
{% endif %}
<b>{{ build.ProjectName }}</b> (<i>{{ build.Branch }}</i>)

{{ build.CommitterName }} ({{ build.CommitterEmail }}) at {{ build.CommittedAt | time:"Mon, 02 Jan 2006 15:04:05 -0700" }}
 
{{ build.CommitMessage }}
{% if build.Error %}<pre>{{ build.Error }}</pre>{% endif %}
sha: {{ build.CommitSHA }}
{% endautoescape %}`))

func (b *bot) Send(build common.Build) error {

	message, err := template.Execute(pongo2.Context{
		"build": build,
	})

	if err != nil {

		return err
	}

	for _, recipient := range build.SendTo {

		if recipient.Method == Method && recipient.IntID != 0 {

			b.SendMessage(&user{TelegramID: recipient.IntID}, message, &telebot.SendOptions{ParseMode: telebot.ModeHTML})
		}
	}

	return nil
}
