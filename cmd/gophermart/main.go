package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zavtra-na-rabotu/gophermart/internal/configuration"
	"github.com/zavtra-na-rabotu/gophermart/internal/db"
	"github.com/zavtra-na-rabotu/gophermart/internal/db/repository"
	"github.com/zavtra-na-rabotu/gophermart/internal/handler"
	"github.com/zavtra-na-rabotu/gophermart/internal/logger"
	"github.com/zavtra-na-rabotu/gophermart/internal/security"
	"github.com/zavtra-na-rabotu/gophermart/internal/service"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	config := configuration.Configure()
	logger.InitLogger()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	// Init database
	dbConnection, err := db.NewDBStorage(config.DatabaseUri)
	if err != nil {
		zap.S().Fatal("Failed to connect to database", zap.Error(err))
	}

	// Run migrations
	err = db.RunMigrations(dbConnection)
	if err != nil {
		zap.S().Fatal("Failed to run migrations", zap.Error(err))
	}

	// Build dependencies for UserHandler
	userRepository := repository.NewUserRepository(dbConnection)
	jwtGenerator := security.NewJwtGenerator(config.JwtSecret, config.JwtLifetimeHours)
	userService := service.NewUserService(userRepository, jwtGenerator)
	userHandler := handler.NewUserHandler(userService)

	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", userHandler.RegisterUser())
		r.Post("/login", userHandler.LoginUser())
		r.Post("/orders", handler.UploadOrder())
		r.Get("/orders", handler.GetOrders())
		r.Get("/balance", handler.GetBalance())
		r.Post("/balance/withdraw", handler.WithdrawBalance())
		r.Get("/withdrawals", handler.GetWithdrawals())
	})

	//TODO: Не забыть обработку сигналов
	err = http.ListenAndServe(config.RunAddress, router)
	if err != nil {
		zap.L().Fatal("Failed to start server", zap.Error(err))
	}
}
