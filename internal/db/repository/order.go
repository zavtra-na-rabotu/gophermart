package repository

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zavtra-na-rabotu/gophermart/internal/model"
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

func (r *OrderRepository) UpdateOrderByNumber(tx *sql.Tx, accrual float64, status string, number string) (*model.Order, error) {
	row := tx.QueryRow(`UPDATE orders SET accrual=$1, status=$2 WHERE number=$3 RETURNING id, number, status, accrual, user_id, uploaded_at`, accrual, status, number)

	var order model.Order
	err := row.Scan(&order.ID, &order.Number, &order.Status, &order.Accrual, &order.UserID, &order.UploadedAt)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) GetAllNotTerminated() ([]model.Order, error) {
	rows, err := r.db.Query(`SELECT id, number, status, user_id, accrual, uploaded_at FROM orders WHERE status not in ('INVALID','PROCESSED')`)
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

	var orders []model.Order
	for rows.Next() {
		var order model.Order

		err = rows.Scan(&order.ID, &order.Number, &order.Status, &order.UserID, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) CreateOrder(orderNumber string, userID int) error {
	_, err := r.db.Exec(`INSERT INTO orders (number, user_id, status) VALUES ($1, $2, $3)`, orderNumber, userID, model.New)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrOrderAlreadyExists
		}
		return err
	}

	return nil
}

func (r *OrderRepository) CreateOrderInTransaction(tx *sql.Tx, orderNumber string, userID int) error {
	_, err := tx.Exec(`INSERT INTO orders (number, user_id, status) VALUES ($1, $2, $3)`, orderNumber, userID, model.New)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrOrderAlreadyExists
		}
		return err
	}

	return nil
}

func (r *OrderRepository) GetOrder(orderNumber string) (*model.Order, error) {
	row := r.db.QueryRow(`SELECT id, number, status, user_id, accrual, uploaded_at FROM orders WHERE number = $1`, orderNumber)

	var order model.Order

	err := row.Scan(&order.ID, &order.Number, &order.Status, &order.UserID, &order.Accrual, &order.UploadedAt)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) GetOrders(userID int) ([]model.Order, error) {
	rows, err := r.db.Query(`SELECT id, number, status, user_id, accrual, uploaded_at FROM orders WHERE user_id=$1 order by uploaded_at`, userID)
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

	var orders []model.Order
	for rows.Next() {
		var order model.Order

		err = rows.Scan(&order.ID, &order.Number, &order.Status, &order.UserID, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
