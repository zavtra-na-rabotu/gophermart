package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDBStorage(databaseURI string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURI)
	return db, err
}
