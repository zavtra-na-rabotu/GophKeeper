package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewDBStorage creates new DB connection
func NewDBStorage(databaseDSN string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseDSN)
	return db, err
}
