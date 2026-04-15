package repository

import "database/sql"

type DebiturRepository struct {
	DB *sql.DB
}

func (r *DebiturRepository) FindByID(id string) (map[string]interface{}, error) {
	var name string

	err := r.DB.QueryRow("SELECT name FROM debitur WHERE id=$1", id).Scan(&name)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":   id,
		"name": name,
	}, nil
}

func (r *DebiturRepository) Create(id int, name string) error {
	_, err := r.DB.Exec(
		"INSERT INTO debitur (id, name) VALUES ($1, $2)",
		id, name,
	)
	return err
}
