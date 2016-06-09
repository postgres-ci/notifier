package telegram

import (
	"fmt"
)

type user struct {
	ID         int32 `db:"user_id"`
	TelegramID int64 `db:"telegram_id"`
}

func (u *user) Destination() string {

	return fmt.Sprint(u.TelegramID)
}
