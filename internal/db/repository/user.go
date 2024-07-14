package repository

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zavtra-na-rabotu/gophermart/internal/model"
	"go.uber.org/zap"
)

var (
	ErrUserAlreadyExists = errors.New("user already exist")
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(tx *sql.Tx, login string, password string) (*model.User, error) {
	row := r.db.QueryRow(`INSERT INTO users (login, password) VALUES ($1, $2) RETURNING *`, login, password)

	var user model.User
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByLogin(login string) (*model.User, error) {
	row := r.db.QueryRow(`SELECT * FROM users WHERE login = $1`, login)

	var user model.User
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		zap.L().Error("Failed to query user by login", zap.String("login", login), zap.Error(err))
		return nil, err
	}

	return &user, nil
}
