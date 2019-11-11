package service

import (
	"math/rand"
	"time"

	"github.com/hatena/go-Intern-Diary/model"
	"github.com/hatena/go-Intern-Diary/repository"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type DiaryApp interface {
	CreateNewUser(name string, passwordHash string) error
	FindUserByName(name string) (*model.User, error)
	LoginUser(name string, password string) (bool, error)
	CreateNewToken(userID uint64, expiresAt time.Time) (string, error)
	FindUserByToken(token string) (*model.User, error)
	CreateNewDiary(user *model.User, name string) (*model.Diary, error)
	Close() error
}

func NewApp(repo repository.Repository) DiaryApp {
	return &diaryApp{repo}
}

type diaryApp struct {
	repo repository.Repository
}

func (app *diaryApp) Close() error {
	return app.repo.Close()
}
