package repository

import "database/sql"

type LogRepository struct {
	DB *sql.DB
}
