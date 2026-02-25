package repository

import "database/sql"

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Create(email, password string) error {
	_, err := r.DB.Exec(
		"INSERT INTO users (email, password) VALUES ($1, $2)",
		email,
		password,
	)
	return err
}

func (r *UserRepository) FindByEmail(email string) (int, string, error) {
	var id int
	var password string

	err := r.DB.QueryRow(
		"SELECT id, password FROM users WHERE email=$1",
		email,
	).Scan(&id, &password)

	return id, password, err
}
