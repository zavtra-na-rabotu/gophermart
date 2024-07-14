package service

import (
	"github.com/zavtra-na-rabotu/gophermart/internal/db/repository"
	"github.com/zavtra-na-rabotu/gophermart/internal/model"
)

type BalanceService struct {
	balanceRepository *repository.BalanceRepository
}

func NewBalanceService(balanceRepository *repository.BalanceRepository) *BalanceService {
	return &BalanceService{balanceRepository: balanceRepository}
}

func (s *BalanceService) GetBalance(userID int) (*model.Balance, error) {
	balance, err := s.balanceRepository.GetBalanceByUserID(userID)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
