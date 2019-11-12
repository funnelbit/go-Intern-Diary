package service

import (
	"errors"

	"github.com/hatena/go-Intern-Diary/model"
)

func (app *diaryApp) CreateNewDiary(user *model.User, name string) (*model.Diary, error) {
	if name == "" {
		return nil, errors.New("empty diary name")
	}
	return app.repo.CreateNewDiary(user.ID, name)
}

func (app *diaryApp) FindDiariesByUserID(userID uint64) ([]*model.Diary, error) {
	return app.repo.FindDiariesByUserID(userID)
}

func (app *diaryApp) FindDiaryByID(diaryID uint64) (*model.Diary, error) {
	if diaryID == 0 {
		return nil, errors.New("empty diary id")
	}
	return app.repo.FindDiaryByID(diaryID)
}
