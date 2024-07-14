package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDBStorage(databaseUri string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseUri)
	return db, err
}
