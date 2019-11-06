package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/hatena/go-Intern-Diary/model"
)

func (app *diaryApp) CreateNewUser(name string, password string) (err error) {
	if name == "" {
		return errors.New("empty user name")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return app.repo.CreateNewUser(name, string(passwordHash))
}

func (app *diaryApp) FindUserByName(name string) (*model.User, error) {
	if name == "" {
		return nil, errors.New("empty user name")
	}

	return app.repo.FindUserByName(name)
}
