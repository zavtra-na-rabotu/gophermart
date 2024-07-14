package handler

import (
	"encoding/json"
	"github.com/zavtra-na-rabotu/gophermart/internal/dto"
	"github.com/zavtra-na-rabotu/gophermart/internal/middleware"
	"github.com/zavtra-na-rabotu/gophermart/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type BalanceHandler struct {
	balanceService *service.BalanceService
}

func NewBalanceHandler(balanceService *service.BalanceService) *BalanceHandler {
	return &BalanceHandler{balanceService: balanceService}
}

func (h *BalanceHandler) GetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(int)

		balance, err := h.balanceService.GetBalance(userID)
		if err != nil {
			zap.L().Error("Failed to get balance", zap.Error(err))
			http.Error(w, "Failed to get balance", http.StatusInternalServerError)
			return
		}

		var response = dto.GetBalanceResponse{
			Current:   balance.Current,
			Withdrawn: balance.Withdrawn,
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
