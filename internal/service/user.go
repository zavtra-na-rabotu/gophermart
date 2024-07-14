package service

import (
	"database/sql"
	"errors"
	"github.com/zavtra-na-rabotu/gophermart/internal/db"
	"github.com/zavtra-na-rabotu/gophermart/internal/db/repository"
	"github.com/zavtra-na-rabotu/gophermart/internal/dto"
	"github.com/zavtra-na-rabotu/gophermart/internal/model"
	"github.com/zavtra-na-rabotu/gophermart/internal/security"
	"go.uber.org/zap"
)

var (
	ErrIncorrectLoginOrPassword = errors.New("incorrect login or password")
)

type UserService struct {
	transactionManager *db.TransactionManager
	userRepository     *repository.UserRepository
	balanceRepository  *repository.BalanceRepository
	jwtGenerator       *security.JwtService
}

func NewUserService(
	transactionManager *db.TransactionManager,
	userRepository *repository.UserRepository,
	balanceRepository *repository.BalanceRepository,
	jwtGenerator *security.JwtService,
) *UserService {
	return &UserService{
		transactionManager: transactionManager,
		userRepository:     userRepository,
		balanceRepository:  balanceRepository,
		jwtGenerator:       jwtGenerator,
	}
}

func (s *UserService) RegisterUser(request *dto.RegisterUserRequest) (string, error) {
	hash, err := security.HashPassword(request.Password)
	if err != nil {
		zap.L().Error("Failed to hash password", zap.Error(err))
		return "", err
	}

	user, err := s.transactionManager.RunInTransaction(func(tx *sql.Tx) (any, error) {
		user, err := s.userRepository.CreateUser(tx, request.Login, hash)
		if err != nil {
			zap.L().Error("Failed to create user", zap.Error(err))
			return nil, err
		}

		_, err = s.balanceRepository.CreateBalance(tx, user.ID)
		if err != nil {
			zap.L().Error("Failed to create balance", zap.Error(err))
			return nil, err
		}

		return user, nil
	})
	if err != nil {
		return "", nil
	}

	return s.jwtGenerator.GenerateJwtToken(user.(*model.User).ID)
}

func (s *UserService) LoginUser(request *dto.LoginUserRequest) (string, error) {
	user, err := s.userRepository.GetUserByLogin(request.Login)
	if err != nil {
		zap.L().Error("User not found", zap.String("login", request.Login), zap.Error(err))
		return "", ErrIncorrectLoginOrPassword
	}

	if !security.CheckPassword(user.Password, request.Password) {
		zap.L().Error("Invalid password")
		return "", ErrIncorrectLoginOrPassword
	}

	return s.jwtGenerator.GenerateJwtToken(user.ID)
}
