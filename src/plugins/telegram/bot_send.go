package telegram

import (
	"github.com/flosch/pongo2"
	"github.com/postgres-ci/notifier/src/common"
)

var template = pongo2.Must(pongo2.FromString(`{% autoescape off %}
{% if build.Status == "failed" %}
Build failed
{% else %}
Build success
{% endif %}
{{ build.ProjectName }} ({{ build.Branch }})

{{ build.CommitterName }} ({{ build.CommitterEmail }}) at {{ build.CommittedAt | time:"Mon, 02 Jan 2006 15:04:05 -0700" }}
 
{{ build.CommitMessage }}
{% if build.Error %}{{ build.Error }}{% endif %}
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

		if recipient.Method == Method {

			b.SendMessage(&user{TelegramID: recipient.IntID}, message, nil)
		}
	}

	return nil
}
