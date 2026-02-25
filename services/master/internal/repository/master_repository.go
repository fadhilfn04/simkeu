package repository

import "database/sql"

type MasterRepository struct {
	DB *sql.DB
}
