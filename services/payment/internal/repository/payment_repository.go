package repository

import "database/sql"

type PaymentRepository struct {
	DB *sql.DB
}
