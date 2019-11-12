package model

import (
	"strconv"
	"time"
)

type Article struct {
	ID      uint64    `db:"id"`
	DiaryID uint64    `db:"diary_id"`
	Body    string    `db:"body"`
	Created time.Time `db:"created_at"`
	Updated time.Time `db:"updated_at"`
}

func (a *Article) CreateArticleURL() string {
	return "/diary/" + strconv.FormatUint(a.DiaryID, 10) + "/" + strconv.FormatUint(a.ID, 10)
}
