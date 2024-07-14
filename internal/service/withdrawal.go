package service

import (
	"database/sql"
	"errors"
	"github.com/zavtra-na-rabotu/gophermart/internal/db"
	"github.com/zavtra-na-rabotu/gophermart/internal/db/repository"
	"github.com/zavtra-na-rabotu/gophermart/internal/model"
)

var (
	ErrNotEnoughBalance = errors.New("not enough balance")
)

type WithdrawalService struct {
	transactionManager   *db.TransactionManager
	withdrawalRepository *repository.WithdrawalRepository
	orderRepository      *repository.OrderRepository
	balanceRepository    *repository.BalanceRepository
}

func NewWithdrawalService(
	transactionManager *db.TransactionManager,
	withdrawalRepository *repository.WithdrawalRepository,
	orderRepository *repository.OrderRepository,
	balanceRepository *repository.BalanceRepository,
) *WithdrawalService {
	return &WithdrawalService{
		transactionManager:   transactionManager,
		withdrawalRepository: withdrawalRepository,
		orderRepository:      orderRepository,
		balanceRepository:    balanceRepository,
	}
}

func (s *WithdrawalService) GetWithdrawals(userID int) ([]model.Withdrawal, error) {
	withdrawals, err := s.withdrawalRepository.GetWithdrawals(userID)
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}

func (s *WithdrawalService) CreateWithdrawal(userID int, orderNumber string, sum float64) error {
	_, err := s.transactionManager.RunInTransaction(func(tx *sql.Tx) (any, error) {
		balance, err := s.balanceRepository.GetBalanceForUpdateByUserID(tx, userID)
		if err != nil {
			return nil, err
		}

		if balance.Current < sum {
			return nil, ErrNotEnoughBalance
		}

		err = s.orderRepository.CreateOrderInTransaction(tx, orderNumber, userID)
		if err != nil {
			return nil, err
		}

		err = s.withdrawalRepository.CreateWithdrawal(tx, userID, orderNumber, sum)
		if err != nil {
			return nil, err
		}

		err = s.balanceRepository.WithdrawByUserID(tx, userID, sum)
		if err != nil {
			return nil, err
		}

		return nil, err
	})

	return err
}
