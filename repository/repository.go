package repository

import (
	"fmt"
	"time"

	//"../model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/hatena/go-Intern-Diary/model"
)

type Repository interface {
	CreateNewUser(name string, passwordHash string) error
	CreateNewToken(userID uint64, token string, expiresAt time.Time) error
	FindUserByName(name string) (*model.User, error)
	FindPasswordHashByName(name string) (string, error)
	//CreateNewToken(userID uint64, token string, expiresAt time.Time) error
	//DeleteToken(userID uint64) error

	Close() error
}

func New(dsn string) (Repository, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Opening mysql failed: %v", err)
	}
	return &repository{db: db}, nil
}

type repository struct {
	db *sqlx.DB
}

func (r *repository) generateID() (uint64, error) {
	var id uint64
	err := r.db.Get(&id, "SELECT UUID_SHORT()")
	return id, err
}

func (r *repository) Close() error {
	return r.db.Close()
}
