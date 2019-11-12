package model

import (
	"time"
	"strconv"
)

type Diary struct {
	ID      uint64    `db:"id"`
	Name    string    `db:"name"`
	UserID  uint64    `db:"user_id"`
	Created time.Time `db:"created_at"`
	Updated time.Time `db:"updated_at"`
}

func (d *Diary) CreateDiaryURL() string {
	return "/diary/" + strconv.FormatUint(d.ID, 10)
}
