package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/zavtra-na-rabotu/gophermart/internal/configuration"
	"github.com/zavtra-na-rabotu/gophermart/internal/db"
	"github.com/zavtra-na-rabotu/gophermart/internal/db/repository"
	"github.com/zavtra-na-rabotu/gophermart/internal/handler"
	"github.com/zavtra-na-rabotu/gophermart/internal/integration"
	"github.com/zavtra-na-rabotu/gophermart/internal/job"
	"github.com/zavtra-na-rabotu/gophermart/internal/logger"
	"github.com/zavtra-na-rabotu/gophermart/internal/middleware"
	"github.com/zavtra-na-rabotu/gophermart/internal/security"
	"github.com/zavtra-na-rabotu/gophermart/internal/service"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	config := configuration.Configure()
	logger.InitLogger()

	router := chi.NewRouter()

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

	transactionManager := db.NewTransactionManager(dbConnection)

	// Build repositories
	orderRepository := repository.NewOrderRepository(dbConnection)
	balanceRepository := repository.NewBalanceRepository(dbConnection)
	userRepository := repository.NewUserRepository(dbConnection)
	withdrawalRepository := repository.NewWithdrawalRepository(dbConnection)

	// Build services
	orderService := service.NewOrderService(transactionManager, orderRepository, balanceRepository)
	balanceService := service.NewBalanceService(balanceRepository)
	jwtService := security.NewJwtService([]byte(config.JwtSecret), config.JwtLifetimeHours)
	userService := service.NewUserService(transactionManager, userRepository, balanceRepository, jwtService)
	withdrawalService := service.NewWithdrawalService(transactionManager, withdrawalRepository, orderRepository, balanceRepository)

	// Build handlers
	orderHandler := handler.NewOrderHandler(orderService)
	balanceHandler := handler.NewBalanceHandler(balanceService)
	userHandler := handler.NewUserHandler(userService)
	withdrawalHandler := handler.NewWithdrawalHandler(withdrawalService)

	accrualClient := integration.NewAccrualClient(config.AccrualSystemAddress)
	accrualJob := job.NewAccrualJob(accrualClient, orderService)

	router.Route("/api/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/register", userHandler.RegisterUser())
			r.Post("/login", userHandler.LoginUser())
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthorizationMiddleware(jwtService))
			r.Post("/orders", orderHandler.CreateOrder())
			r.Get("/orders", orderHandler.GetOrders())
			r.Get("/balance", balanceHandler.GetBalance())
			r.Post("/balance/withdraw", withdrawalHandler.CreateWithdrawal())
			r.Get("/withdrawals", withdrawalHandler.GetWithdrawals())
		})
	})

	go func() {
		ticker := time.NewTicker(1000 * time.Millisecond)
		for range ticker.C {
			accrualJob.Start()
		}
	}()

	//TODO: Не забыть обработку сигналов
	err = http.ListenAndServe(config.RunAddress, router)
	if err != nil {
		zap.L().Fatal("Failed to start server", zap.Error(err))
	}
}
