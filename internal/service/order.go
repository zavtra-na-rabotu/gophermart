package service

import (
	"database/sql"
	"errors"
	"github.com/zavtra-na-rabotu/gophermart/internal/db"
	"github.com/zavtra-na-rabotu/gophermart/internal/db/repository"
	"github.com/zavtra-na-rabotu/gophermart/internal/model"
)

var (
	ErrOrderCreatedByAnotherUser = errors.New("order created by another user")
)

type OrderService struct {
	transactionManager *db.TransactionManager
	orderRepository    *repository.OrderRepository
	balanceRepository  *repository.BalanceRepository
}

func NewOrderService(
	transactionManager *db.TransactionManager,
	orderRepository *repository.OrderRepository,
	balanceRepository *repository.BalanceRepository,
) *OrderService {
	return &OrderService{
		transactionManager: transactionManager,
		orderRepository:    orderRepository,
		balanceRepository:  balanceRepository,
	}
}

func (s *OrderService) CreateOrder(orderNumber string, userID int) error {
	var finalError error

	err := s.orderRepository.CreateOrder(orderNumber, userID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderAlreadyExists) {
			finalError = err
		} else {
			return err
		}
	}

	order, err := s.orderRepository.GetOrder(orderNumber)
	if err != nil {
		return err
	}

	if finalError != nil {
		if userID == order.UserID {
			return finalError
		} else {
			return ErrOrderCreatedByAnotherUser
		}
	}

	return nil
}

func (s *OrderService) UpdateOrder(orderNumber string, accrual float64, status string) error {
	_, err := s.transactionManager.RunInTransaction(func(tx *sql.Tx) (any, error) {
		order, err := s.orderRepository.UpdateOrderByNumber(tx, accrual, status, orderNumber)
		if err != nil {
			return nil, err
		}

		err = s.balanceRepository.AccrueByUserID(tx, order.UserID, accrual)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	return err
}

func (s *OrderService) GetOrders(userID int) ([]model.Order, error) {
	orders, err := s.orderRepository.GetOrders(userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *OrderService) GetAllNotTerminated() ([]model.Order, error) {
	orders, err := s.orderRepository.GetAllNotTerminated()
	if err != nil {
		return nil, err
	}

	return orders, nil
}
