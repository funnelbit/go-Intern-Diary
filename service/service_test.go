package service

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/hatena/go-Intern-Diary/config"
	"github.com/hatena/go-Intern-Diary/model"
	"github.com/hatena/go-Intern-Diary/repository"
)

func newApp() DiaryApp {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}
	repo, err := repository.New(conf.DbDsn)
	if err != nil {
		panic(err)
	}
	return NewApp(repo)
}

func closeApp(app DiaryApp) {
	err := app.Close()
	if err != nil {
		panic(err)
	}
}

func randomString() string {
	return strconv.FormatInt(time.Now().Unix()^rand.Int63(), 16)
}

func createUser(app DiaryApp) *model.User {
	name := "test name "
	password := "a"
	err := app.CreateNewUser(name, password)
	user, err := app.FindUserByName(name)
	if err != nil {
		panic(err)
	}
	return user
}

func createDiary(app DiaryApp, user *model.User, diaryName string) *model.Diary {
	diary, err := app.CreateNewDiary(user, diaryName)
	if err != nil {
		panic(err)
	}
	return diary
}
