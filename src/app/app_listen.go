package app

import (
	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"
	"github.com/postgres-ci/notifier/src/common"

	"database/sql"
	"time"
)

const (
	minReconnectInterval = time.Second
	maxReconnectInterval = 5 * time.Second
	channelNotifications = "postgres-ci::notification"
)

func (a *app) listen() {

	listener := pq.NewListener(a.config.Connect.DSN(), minReconnectInterval, maxReconnectInterval, func(event pq.ListenerEventType, err error) {

		if err != nil {

			log.Errorf("Postgres listen: %v", err)

			return
		}

		log.Debugf("Postgres notify send event: %v", event)
	})

	listener.Listen(channelNotifications)

	var (
		events             = listener.NotificationChannel()
		checkNotifications = time.Tick(time.Minute)
	)

	for {

		select {

		case event := <-events:

			if event == nil {

				continue
			}

			log.Debugf("Received from [%s] playload: %s", event.Channel, event.Extra)

			switch event.Channel {
			case channelNotifications:
				a.checkNotifications()
			}

		case <-checkNotifications:
			a.checkNotifications()
		}
	}
}

func (a *app) checkNotifications() {

	for {

		var build common.Build

		if err := a.connect.Get(&build, checkNotificationsSql); err != nil {

			if err == sql.ErrNoRows {

				log.Debug("All notifications fetched")

			} else {

				log.Errorf("Error when fetching notifications: %v", err)
			}

			return
		}

		for _, plugin := range a.plugins {

			plugin.Send(build)
		}
	}
}

const checkNotificationsSql = `
	SELECT 
		build_id,
		build_status,
		project_id,
		project_name,
		branch,
		build_error,
		build_created_at,
		build_started_at,
		build_finished_at,
		commit_sha,
		commit_message,
		committed_at,
		committer_name,
		committer_email,
		commit_author_name,
		commit_author_email,
		send_to
	FROM notification.fetch()
`
