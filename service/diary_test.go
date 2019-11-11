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

func TestDiaryApp_FindDiariesByUserID(t *testing.T) {
	app := newApp()
	defer closeApp(app)

	user := createUser(app)
	other := createUser(app)

	createDiary(app, user, "diaryName1")
	createDiary(app, user, "diaryName2")
	createDiary(app, user, "diaryName3")
	createDiary(app, other, "diaryNameOther") // 他人なので引かれないはず

	diaries, err := app.FindDiariesByUserID(user.ID)
	assert.Nil(err)
	assert.Equal(t, 3, len(diaries))
}
