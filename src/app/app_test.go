package app

import (
	"github.com/erikstmartin/go-testdb"
	"github.com/jmoiron/sqlx"
	"github.com/postgres-ci/notifier/src/common"
	"github.com/stretchr/testify/assert"

	"database/sql/driver"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

type testPlugin struct {
	send func(common.Build) error
}

func (t *testPlugin) Send(build common.Build) error {

	return t.send(build)
}

func TestCheckNotifications(t *testing.T) {

	var returnErrNoRows bool

	testdb.EnableTimeParsing(true)

	testdb.SetQueryFunc(func(query string) (result driver.Rows, err error) {

		columns := `
			build_id
			build_status
			project_id
			project_name
			branch
			build_error
			build_created_at
			build_started_at
			build_finished_at
			commit_sha
			commit_message
			committed_at
			committer_name
			committer_email
			commit_author_name
			commit_author_email
			send_to
		`

		sendto, _ := json.Marshal([]common.Recipient{
			{
				Name:   "Name",
				Method: "Method",
				TextID: "TextID",
				IntID:  42,
			},
		})

		if returnErrNoRows {

			return testdb.RowsFromSlice(
				strings.Fields(columns),
				[][]driver.Value{},
			), nil
		}

		returnErrNoRows = true

		return testdb.RowsFromSlice(
			strings.Fields(columns),
			[][]driver.Value{
				{
					"1",
					"pending",
					"1",
					"Project",
					"master",
					"error",
					time.Now(),
					time.Now(),
					time.Now(),
					"7c140d12eb6b02552e1df13e62d4b47514c93d3b",
					"Message",
					time.Now(),
					"Elephant Sam",
					"samelephant82@gmail.com",
					"Elephant Sam",
					"samelephant82@gmail.com",
					string(sendto),
				},
			},
		), nil
	})

	var build *common.Build

	app := &app{
		connect: sqlx.MustOpen("testdb", ""),
		plugins: []plugin{
			&testPlugin{
				send: func(b common.Build) error {

					build = &b

					return nil
				},
			},
		},
	}

	app.checkNotifications()

	if assert.NotNil(t, build) {

		assert.Equal(t, int32(1), build.ProjectID)
		assert.Equal(t, "pending", build.Status)
		assert.Equal(t, "Project", build.ProjectName)
		assert.Equal(t, "master", build.Branch)

		if assert.Len(t, build.SendTo, 1) {

			assert.Equal(t, "Name", build.SendTo[0].Name)
			assert.Equal(t, "Method", build.SendTo[0].Method)
			assert.Equal(t, "TextID", build.SendTo[0].TextID)
			assert.Equal(t, int64(42), build.SendTo[0].IntID)
		}
	}
}
