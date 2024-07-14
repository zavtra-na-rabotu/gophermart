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
	"net/http"
	"time"
)

var (
	ErrNotEnoughBalance = errors.New("not enough balance")
)

type WithdrawalHandler struct {
	withdrawalService *service.WithdrawalService
}

func NewWithdrawalHandler(withdrawalService *service.WithdrawalService) *WithdrawalHandler {
	return &WithdrawalHandler{withdrawalService: withdrawalService}
}

func (h *WithdrawalHandler) GetWithdrawals() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(int)

		withdrawals, err := h.withdrawalService.GetWithdrawals(userID)
		if err != nil {
			if errors.Is(err, repository.ErrNoWithdrawalsFound) {
				http.Error(w, "No withdrawals found", http.StatusNoContent)
				return
			}
			zap.L().Error("Failed to get withdrawals", zap.Error(err))
			http.Error(w, "Failed to get withdrawals", http.StatusInternalServerError)
			return
		}

		response := make([]dto.GetWithdrawalsResponse, len(withdrawals))
		for i, withdrawal := range withdrawals {
			response[i] = dto.GetWithdrawalsResponse{
				Order:       withdrawal.OrderNumber,
				Sum:         withdrawal.Sum,
				ProcessedAt: withdrawal.ProcessedAt.Format(time.RFC3339),
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

func (h *WithdrawalHandler) CreateWithdrawal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			http.Error(w, "Invalid request content type", http.StatusBadRequest)
			return
		}

		userID := r.Context().Value(middleware.UserIDKey).(int)

		var request dto.CreateWithdrawalRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			zap.L().Error("Failed to parse body", zap.Error(err))
			http.Error(w, "Failed to parse body", http.StatusBadRequest)
			return
		}

		if !luhn.Valid(request.Order) {
			http.Error(w, "Bad order number", http.StatusUnprocessableEntity)
			return
		}

		err := h.withdrawalService.CreateWithdrawal(userID, request.Order, request.Sum)
		if err != nil {
			if errors.Is(err, ErrNotEnoughBalance) {
				http.Error(w, "Not enough balance", http.StatusPaymentRequired)
				return
			}
			if errors.Is(err, repository.ErrOrderAlreadyExists) {
				http.Error(w, "Order already exists", http.StatusOK)
				return
			}
			if errors.Is(err, service.ErrOrderCreatedByAnotherUser) {
				http.Error(w, "Order created by another user", http.StatusConflict)
				return
			}
			zap.L().Error("Failed to create withdrawal", zap.Error(err))
			http.Error(w, "Failed to create withdrawal", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
