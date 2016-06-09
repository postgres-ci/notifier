package common

import (
	"encoding/json"
	"fmt"
	"time"
)

type Build struct {
	ID                int32      `db:"build_id"`
	Status            string     `db:"build_status"`
	Branch            string     `db:"branch"`
	Error             string     `db:"build_error"`
	CreatedAt         time.Time  `db:"build_created_at"`
	StartedAt         time.Time  `db:"build_started_at"`
	FinishedAt        time.Time  `db:"build_finished_at"`
	CommitSHA         string     `db:"commit_sha"`
	CommitMessage     string     `db:"commit_message"`
	CommittedAt       time.Time  `db:"committed_at"`
	CommitterName     string     `db:"committer_name"`
	CommitterEmail    string     `db:"committer_email"`
	CommitAuthorName  string     `db:"commit_author_name"`
	CommitAuthorEmail string     `db:"commit_author_email"`
	SendTo            recipients `db:"send_to"`
}

type recipient struct {
	Name   string `json:"user_name"`
	Method string `json:"notify_method"`
	TextID string `json:"notify_text_id"`
	IntID  int64  `json:"notify_int_id"`
}

type recipients []recipient

func (r *recipients) Scan(src interface{}) error {

	var source []byte

	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	default:
		return fmt.Errorf("Incompatible type for recipients")
	}

	return json.Unmarshal(source, r)
}
