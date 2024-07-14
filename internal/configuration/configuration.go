package configuration

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"github.com/zavtra-na-rabotu/gophermart/internal/utils/stringutils"
	"go.uber.org/zap"
	"os"
)

type Configuration struct {
	RunAddress           string
	DatabaseUri          string
	AccrualSystemAddress string
	JwtSecret            string
	JwtLifetimeHours     int
}

type envs struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseUri          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JwtSecret            string `env:"JWT_SECRET"`
	JwtLifetimeHours     int    `env:"JWT_LIFETIME_HOURS"`
}

func Configure() *Configuration {
	var config Configuration

	flag.StringVar(&config.RunAddress, "a", "localhost:8080", "Адрес и порт запуска сервиса")
	flag.StringVar(&config.DatabaseUri, "d", "", "Адрес подключения к базе данных")
	flag.StringVar(&config.AccrualSystemAddress, "r", "", "Адрес системы расчёта начислений")
	flag.StringVar(&config.JwtSecret, "j", "secret", "JWT секрет")
	flag.IntVar(&config.JwtLifetimeHours, "l", 24, "Время жизни JWT токена в часах")
	flag.Parse()

	envVariables := envs{}
	err := env.Parse(&envVariables)
	if err != nil {
		zap.L().Error("Failed to parse environment variables", zap.Error(err))
	}

	_, exists := os.LookupEnv("RUN_ADDRESS")
	if exists && !stringutils.IsEmpty(envVariables.RunAddress) {
		config.RunAddress = envVariables.RunAddress
	}

	_, exists = os.LookupEnv("DATABASE_URI")
	if exists {
		config.DatabaseUri = envVariables.DatabaseUri
	}

	_, exists = os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if exists {
		config.AccrualSystemAddress = envVariables.AccrualSystemAddress
	}

	_, exists = os.LookupEnv("JWT_SECRET")
	if exists {
		config.JwtSecret = envVariables.JwtSecret
	}

	_, exists = os.LookupEnv("JWT_LIFETIME_HOURS")
	if exists {
		config.JwtLifetimeHours = envVariables.JwtLifetimeHours
	}

	return &config
}
