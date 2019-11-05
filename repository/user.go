package repository

import "time"

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
