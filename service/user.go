package service

import (
	"time"
	"math/rand"
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

func (app *diaryApp) CreateNewToken(userID uint64, expiresAt time.Time) (string, error) {
	token := generateToken()
	err := app.repo.CreateNewToken(userID, token, expiresAt)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (app *diaryApp) FindUserByName(name string) (*model.User, error) {
	if name == "" {
		return nil, errors.New("empty user name")
	}

	return app.repo.FindUserByName(name)
}

func (app *diaryApp) LoginUser(name string, password string) (bool, error) {
	if name == "" {
		return false, errors.New("empty user name")
	}

	if password == "" {
		return false, errors.New("empty user passowrd")
	}

	passwordHash, err := app.repo.FindPasswordHashByName(name)
	if err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (app *diaryApp) FindUserByToken(token string) (*model.User, error) {
	if token == "" {
		return nil, errors.New("empty user token")
	}

	user, err := app.repo.FindUserByToken(token)
	return user, err
}

func generateToken() string {
	table := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_@"
	l := len(table)
	ret := make([]byte, 128)
	src := make([]byte, 128)
	rand.Read(src)
	for i := 0; i < 128; i++ {
		ret[i] = table[int(src[i])%l]
	}
	return string(ret)
}
