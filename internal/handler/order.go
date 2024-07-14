package handler

import (
	"encoding/json"
	"errors"
	"github.com/zavtra-na-rabotu/gophermart/internal/db/repository"
	"github.com/zavtra-na-rabotu/gophermart/internal/dto"
	"github.com/zavtra-na-rabotu/gophermart/internal/middleware"
	"github.com/zavtra-na-rabotu/gophermart/internal/service"
	"github.com/zavtra-na-rabotu/gophermart/internal/utils/luhn"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType != "text/plain" {
			http.Error(w, "Invalid request content type", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Incorrect input json", http.StatusInternalServerError)
			return
		}

		orderNumber := string(body)
		if !luhn.Valid(orderNumber) {
			http.Error(w, "Bad order number", http.StatusUnprocessableEntity)
			return
		}

		userID := r.Context().Value(middleware.UserIDKey).(int)
		err = h.orderService.CreateOrder(orderNumber, userID)
		if err != nil {
			if errors.Is(err, repository.ErrOrderAlreadyExists) {
				http.Error(w, "Order already exists", http.StatusOK)
				return
			}
			if errors.Is(err, service.ErrOrderCreatedByAnotherUser) {
				http.Error(w, "Order created by another user", http.StatusConflict)
				return
			}
			zap.L().Error("Failed to create order", zap.Error(err))
			http.Error(w, "Failed to create order", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func (h *OrderHandler) GetOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(int)

		orders, err := h.orderService.GetOrders(userID)
		if err != nil {
			if errors.Is(err, repository.ErrNoOrdersFound) {
				http.Error(w, "No orders found", http.StatusNoContent)
				return
			}
			zap.L().Error("Failed to get orders", zap.Error(err))
			http.Error(w, "Failed to get orders", http.StatusInternalServerError)
			return
		}

		response := make([]dto.GetOrdersResponse, len(orders))
		for i, order := range orders {
			response[i] = dto.GetOrdersResponse{
				Number:     order.Number,
				Status:     order.Status,
				Accrual:    order.Accrual,
				UploadedAt: order.UploadedAt.Format(time.RFC3339),
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			zap.L().Error("Failed to write response", zap.Error(err))
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	}
}
