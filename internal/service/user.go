package service

import (
	"errors"
	"github.com/zavtra-na-rabotu/gophermart/internal/db/repository"
	"github.com/zavtra-na-rabotu/gophermart/internal/dto"
	"github.com/zavtra-na-rabotu/gophermart/internal/security"
	"go.uber.org/zap"
)

//type UserService interface {
//	RegisterUser(user model.User) error
//}

var (
	ErrIncorrectLoginOrPassword = errors.New("incorrect login or password")
)

type UserService struct {
	userRepository *repository.UserRepository
	jwtGenerator   *security.JwtGenerator
}

func NewUserService(userRepository *repository.UserRepository, jwtGenerator *security.JwtGenerator) *UserService {
	return &UserService{userRepository: userRepository, jwtGenerator: jwtGenerator}
}

func (s *UserService) RegisterUser(request *dto.RegisterUserRequest) (string, error) {
	hash, err := security.HashPassword(request.Password)
	if err != nil {
		zap.L().Error("Failed to hash password", zap.Error(err))
		return "", err
	}

	err = s.userRepository.CreateUser(request.Login, hash)
	if err != nil {
		zap.L().Error("Failed to create user", zap.Error(err))
		return "", err
	}

	return s.jwtGenerator.GenerateJwtToken(request.Login)
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

	return s.jwtGenerator.GenerateJwtToken(user.Login)
}
