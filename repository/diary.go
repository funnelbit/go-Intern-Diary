package repository

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hatena/go-Intern-Diary/model"
)

func (r *repository) CreateNewDiary(userID uint64, diaryName string) (*model.Diary, error) {
	id, err := r.generateID()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	_, err = r.db.Exec(
		`INSERT INTO diary
			(id, name, user_id, created_at, updated_at)
			VALUES(?, ?, ?, ?, ?)`,
		id, diaryName, userID, now, now,
	)
	if err != nil {
		return nil, err
	}

	return &model.Diary{ID: id, Name: diaryName, UserID: userID, Created: now}, nil
}

func (r *repository) FindDiariesByUserID(userID uint64) ([]*model.Diary, error) {
	var diaries []*model.Diary

	err := r.db.Select(
		&diaries,
		`SELECT * FROM diary WHERE user_id = ?`,
		userID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return diaries, nil
		}
		return nil, err
	}

	return diaries, err
}
