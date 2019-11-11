package service

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDiaryApp_CreateNewDiary(t *testing.T) {
	app := newApp()
	defer closeApp(app)

	user := createUser(app)
	diaryName := "newDiary!!"
	diary, err := app.CreateNewDiary(user, diaryName)
	assert.Nil(t, err)
	assert.NotNil(t, diary)
}
