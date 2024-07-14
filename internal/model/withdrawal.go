package model

import "time"

type Withdrawal struct {
	ID          int
	UserID      int
	OrderNumber string
	Sum         float64
	ProcessedAt time.Time
}
