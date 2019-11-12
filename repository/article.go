package repository

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hatena/go-Intern-Diary/model"
)

var articleNotFoundError = model.NotFoundError("article")

func (r *repository) CreateNewArticle(diaryID uint64, body string) (*model.Article, error) {
	id, err := r.generateID()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	_, err = r.db.Exec(
		`INSERT INTO article
			(id, diary_id, body, created_at, updated_at)
			VALUES(?, ?, ?, ?, ?)`,
		id, diaryID, body, now, now,
	)

	return &model.Article{
		ID:      id,
		DiaryID: diaryID,
		Body:    body,
		Created: now,
		Updated: now,
	}, nil
}

func (r *repository) FindArticlesByDiaryID(diaryID uint64) ([]*model.Article, error) {
	var articles []*model.Article

	err := r.db.Select(
		&articles,
		`SELECT * FROM article WHERE diary_id = ?`,
		diaryID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return articles, nil
		}
		return articles, err
	}

	return articles, nil
}

func (r *repository) FindArticleByID(articleID uint64) (*model.Article, error) {
	var article model.Article
	err := r.db.Get(
		&articleID,
		`SELECT * FROM article WHERE id = ? LIMIT 1`,
		articleID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, articleNotFoundError
		}
		return nil, err
	}

	return &article, nil
}
