package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/zavtra-na-rabotu/gophermart/internal/configuration"
	"github.com/zavtra-na-rabotu/gophermart/internal/db"
	"github.com/zavtra-na-rabotu/gophermart/internal/db/repository"
	"github.com/zavtra-na-rabotu/gophermart/internal/handler"
	"github.com/zavtra-na-rabotu/gophermart/internal/logger"
	"github.com/zavtra-na-rabotu/gophermart/internal/middleware"
	"github.com/zavtra-na-rabotu/gophermart/internal/security"
	"github.com/zavtra-na-rabotu/gophermart/internal/service"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	config := configuration.Configure()
	logger.InitLogger()

	router := chi.NewRouter()
	//router.Use(middleware.RequestID)
	//router.Use(middleware.Logger)

	// Init database
	dbConnection, err := db.NewDBStorage(config.DatabaseURI)
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
	jwtService := security.NewJwtService([]byte(config.JwtSecret), config.JwtLifetimeHours)
	userService := service.NewUserService(userRepository, jwtService)
	userHandler := handler.NewUserHandler(userService)

	// Build dependencies for OrderHandler
	orderRepository := repository.NewOrderRepository(dbConnection)
	orderService := service.NewOrderService(orderRepository)
	orderHandler := handler.NewOrderHandler(orderService)

	router.Route("/api/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/register", userHandler.RegisterUser())
			r.Post("/login", userHandler.LoginUser())
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthorizationMiddleware(jwtService))
			r.Post("/orders", orderHandler.CreateOrder())
			r.Get("/orders", orderHandler.GetOrders())
			r.Get("/balance", handler.GetBalance())
			r.Post("/balance/withdraw", handler.WithdrawBalance())
			r.Get("/withdrawals", handler.GetWithdrawals())
		})
	})

	//TODO: Не забыть обработку сигналов
	err = http.ListenAndServe(config.RunAddress, router)
	if err != nil {
		zap.L().Fatal("Failed to start server", zap.Error(err))
	}
}
