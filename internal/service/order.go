package service

import (
	"errors"
	"github.com/zavtra-na-rabotu/gophermart/internal/db/repository"
	"github.com/zavtra-na-rabotu/gophermart/internal/model"
)

var (
	ErrOrderCreatedByAnotherUser = errors.New("order created by another user")
)

type OrderService struct {
	orderRepository *repository.OrderRepository
}

func NewOrderService(orderRepository *repository.OrderRepository) *OrderService {
	return &OrderService{orderRepository: orderRepository}
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

func (s *OrderService) GetOrders(userID int) ([]model.Order, error) {
	orders, err := s.orderRepository.GetOrders(userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
