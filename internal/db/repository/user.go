package repository

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zavtra-na-rabotu/gophermart/internal/model"
	"go.uber.org/zap"
)

//type UserRepository interface {
//	CreateUser(login string, password string) error
//}

var (
	ErrUniqueConstraintViolation = errors.New("unique constraint violation")
	ErrUserNotFound              = errors.New("user not found")
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (d *UserRepository) CreateUser(login string, password string) error {
	_, err := d.db.Exec(`
		INSERT INTO users (login, password) VALUES ($1, $2)
	`, login, password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrUniqueConstraintViolation
		}
		return err
	}

	return nil
}

func (d *UserRepository) GetUserByLogin(login string) (*model.User, error) {
	row := d.db.QueryRow(`SELECT * FROM users WHERE login = $1`, login)

	var user model.User
	err := row.Scan(&user.Id, &user.Login, &user.Password)
	if err != nil {
		zap.L().Error("Failed to query user by login", zap.String("login", login), zap.Error(err))
		//if errors.Is(err, sql.ErrNoRows) {
		//	return nil, ErrUserNotFound
		//}
		return nil, err
	}

	return &user, nil
}
