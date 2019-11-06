package repository

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hatena/go-Intern-Diary/model"
)

var userNotFoundError = model.NotFoundError("user")

func (r *repository) CreateNewUser(name string, passwordHash string) error {
	id, error := r.generateID()
	if error != nil {
		return error
	}
	now := time.Now()
	_, error = r.db.Exec(
		`INSERT INTO user
			(id, name, password_hash, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?)`,
		id, name, passwordHash, now, now,
	)
	return error
}

func (r *repository) FindUserByName(name string) (*model.User, error) {
	var user model.User
	err := r.db.Get(
		&user,
		`SELECT id, name FROM user WHERE name = ? LIMIT 1`,
		name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, userNotFoundError
		}
		return nil, err
	}

	return &user, nil
}
