package service

import (
	"errors"

	"github.com/hatena/go-Intern-Diary/model"
)

func (app *diaryApp) CreateNewArticle(diaryID uint64, body string) (*model.Article, error) {
	if diaryID == 0 {
		return nil, errors.New("empty diary id")
	}
	return app.repo.CreateNewArticle(diaryID, body)
}

func (app *diaryApp) FindArticlesByDiaryID(diaryID uint64) ([]*model.Article, error) {
	if diaryID == 0 {
		return nil, errors.New("empty diary id")
	}
	return app.repo.FindArticlesByDiaryID(diaryID)
}

func (app *diaryApp) FindArticleByID(articleID uint64) (*model.Article, error) {
	if articleID == 0 {
		return nil, errors.New("empty diary id")
	}
	return app.repo.FindArticleByID(articleID)
}
