package model

import (
	"time"
)

type OrderStatus string

const (
	New        OrderStatus = "NEW"
	Processing OrderStatus = "PROCESSING"
	Invalid    OrderStatus = "INVALID"
	Processed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID         int
	UserID     int
	Number     string
	Status     OrderStatus
	Accrual    float64
	UploadedAt time.Time
}
