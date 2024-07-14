package repository

import (
	"database/sql"
	"github.com/zavtra-na-rabotu/gophermart/internal/model"
)

type BalanceRepository struct {
	db *sql.DB
}

func NewBalanceRepository(db *sql.DB) *BalanceRepository {
	return &BalanceRepository{db: db}
}

func (r *BalanceRepository) WithdrawByUserID(tx *sql.Tx, userID int, sum float64) error {
	_, err := tx.Exec(`UPDATE balances SET current = balances.current - $1, withdrawn = withdrawn + $1 WHERE user_id = $2`, sum, userID)
	return err
}

func (r *BalanceRepository) AccrueByUserID(tx *sql.Tx, userID int, accrual float64) error {
	_, err := tx.Exec(`UPDATE balances SET current = current + $1 WHERE user_id = $2`, accrual, userID)
	return err
}

func (r *BalanceRepository) CreateBalance(tx *sql.Tx, userID int) (*model.Balance, error) {
	row := r.db.QueryRow(`INSERT INTO balances (user_id) VALUES ($1) RETURNING id, user_id, current, withdrawn`, userID)

	var balance model.Balance
	err := row.Scan(&balance.ID, &balance.UserID, &balance.Current, &balance.Withdrawn)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (r *BalanceRepository) GetBalanceByUserID(userID int) (*model.Balance, error) {
	row := r.db.QueryRow(`SELECT id, user_id, current, withdrawn FROM balances WHERE user_id = $1;`, userID)

	var balance model.Balance
	err := row.Scan(&balance.ID, &balance.UserID, &balance.Current, &balance.Withdrawn)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (r *BalanceRepository) GetBalanceForUpdateByUserID(tx *sql.Tx, userID int) (*model.Balance, error) {
	row := tx.QueryRow(`SELECT id, user_id, current, withdrawn FROM balances WHERE user_id = $1 FOR UPDATE`, userID)

	var balance model.Balance
	err := row.Scan(&balance.ID, &balance.UserID, &balance.Current, &balance.Withdrawn)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}
