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
	Number     string
	UserID     int
	Status     OrderStatus
	Accrual    float64
	UploadedAt time.Time
}
