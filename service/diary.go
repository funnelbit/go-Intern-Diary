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
