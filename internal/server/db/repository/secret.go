package repository

import "database/sql"

type SecretRepository struct {
	db *sql.DB
}

func NewDataRepository(db *sql.DB) *SecretRepository {
	return &SecretRepository{db: db}
}
