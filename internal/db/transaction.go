package db

import (
	"database/sql"
	"go.uber.org/zap"
)

type TransactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

type TxFunc func(tx *sql.Tx) (interface{}, error)

func (tm *TransactionManager) RunInTransaction(txFunc TxFunc) (interface{}, error) {
	tx, err := tm.db.Begin()
	if err != nil {
		return nil, err
	}

	result, err := txFunc(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			zap.L().Error("Failed to rollback transaction", zap.Error(rbErr))
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}
