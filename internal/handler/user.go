package handler

import (
	"encoding/json"
	"errors"
	"github.com/zavtra-na-rabotu/gophermart/internal/db/repository"
	"github.com/zavtra-na-rabotu/gophermart/internal/dto"
	"github.com/zavtra-na-rabotu/gophermart/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			http.Error(w, "Invalid request content type", http.StatusBadRequest)
			return
		}

		var request dto.RegisterUserRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			zap.L().Error("Failed to read body", zap.Error(err))
			http.Error(w, "Failed to parse body", http.StatusBadRequest)
			return
		}

		//var user = model.User{Login: request.Login, Password: request.Password}
		token, err := h.userService.RegisterUser(&request)
		if err != nil {
			if errors.Is(err, repository.ErrUserAlreadyExists) {
				http.Error(w, "User already exists", http.StatusConflict)
				return
			}
			zap.L().Error("Failed to register user", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Authorization", "Bearer "+token)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *UserHandler) LoginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			http.Error(w, "Invalid request content type", http.StatusBadRequest)
			return
		}

		var request dto.LoginUserRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			zap.L().Error("Failed to read body", zap.Error(err))
			http.Error(w, "Failed to parse body", http.StatusBadRequest)
			return
		}

		token, err := h.userService.LoginUser(&request)
		if err != nil {
			zap.L().Error("Failed to login user", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Authorization", "Bearer "+token)
		w.WriteHeader(http.StatusOK)
	}
}

func GetOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func GetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func WithdrawBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func GetWithdrawals() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
