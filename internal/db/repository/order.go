package repository

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zavtra-na-rabotu/gophermart/internal/model"
	"go.uber.org/zap"
	"time"
)

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrNoOrdersFound      = errors.New("no orders found")
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (d *OrderRepository) CreateOrder(orderNumber string, userID int) error {
	_, err := d.db.Exec(`INSERT INTO orders (number, user_id, status) VALUES ($1, $2, $3)`, orderNumber, userID, model.New)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrOrderAlreadyExists
		}
		return err
	}

	return nil
}

func (d *OrderRepository) GetOrder(number string) (*model.Order, error) {
	row := d.db.QueryRow(`SELECT id, number, status, user_id, accrual, uploaded_at FROM orders WHERE number = $1`, number)

	var order model.Order
	err := row.Scan(&order.ID, &order.Number, &order.Status, &order.UserID, &order.Accrual, &order.UploadedAt)
	if err != nil {
		zap.L().Error("Failed to get order by number", zap.String("number", number), zap.Error(err))
		return nil, err
	}

	return &order, nil
}

func (d *OrderRepository) GetOrders(userID int) ([]model.Order, error) {
	var orders []model.Order

	rows, err := d.db.Query(`SELECT id, number, status, user_id, accrual, uploaded_at FROM orders WHERE user_id=$1 order by uploaded_at`, userID)
	if err != nil {
		return nil, err
	}

	err = rows.Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoOrdersFound
		}
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var ID, userID int
		var number string
		var status model.OrderStatus
		var accrual sql.NullFloat64
		var uploadedAt time.Time

		err = rows.Scan(&ID, &number, &status, &userID, &accrual, &uploadedAt)
		if err != nil {
			return nil, err
		}
		order := model.Order{
			Number:     number,
			Status:     status,
			Accrual:    accrual.Float64,
			UserID:     userID,
			UploadedAt: uploadedAt,
		}
		orders = append(orders, order)
	}

	return orders, nil
}
