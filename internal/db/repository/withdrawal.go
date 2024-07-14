package repository

import (
	"database/sql"
	"errors"
	"github.com/zavtra-na-rabotu/gophermart/internal/model"
)

var (
	ErrNoWithdrawalsFound = errors.New("no withdrawals found")
)

type WithdrawalRepository struct {
	db *sql.DB
}

func NewWithdrawalRepository(db *sql.DB) *WithdrawalRepository {
	return &WithdrawalRepository{db: db}
}

func (r *WithdrawalRepository) GetWithdrawals(userID int) ([]model.Withdrawal, error) {
	var withdrawals []model.Withdrawal

	rows, err := r.db.Query(`SELECT id, user_id, order_number, sum, processed_at FROM withdrawals WHERE user_id = $1 ORDER BY processed_at;`, userID)
	if err != nil {
		return nil, err
	}

	err = rows.Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoWithdrawalsFound
		}
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var withdrawal model.Withdrawal

		err = rows.Scan(&withdrawal.ID, &withdrawal.UserID, &withdrawal.OrderNumber, &withdrawal.Sum, &withdrawal.ProcessedAt)
		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	return withdrawals, nil
}

func (r *WithdrawalRepository) CreateWithdrawal(tx *sql.Tx, userID int, orderNumber string, sum float64) error {
	_, err := tx.Exec(`INSERT INTO withdrawals (user_id, order_number, sum) VALUES ($1, $2, $3)`, userID, orderNumber, sum)
	return err
}
