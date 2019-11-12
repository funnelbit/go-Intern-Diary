package service

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDiaryApp_CreateNewArticle(t *testing.T) {
	app := newApp()
	defer closeApp(app)

	user := createUser(app)
	diary := createDiary(app, user, "test1")

	article, err := app.CreateNewArticle(diary.ID, "aaaaaaaaaaaaaaaaa")
	assert.Nil(t, err)
	assert.NotNil(t, article)
}

func TestDiaryApp_FindArticlesByDiaryID(t *testing.T) {
	
}
