package repository

import "database/sql"

type BlockchainRepository struct {
	DB *sql.DB
}
