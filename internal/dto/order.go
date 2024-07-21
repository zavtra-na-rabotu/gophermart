package dto

import (
	"github.com/zavtra-na-rabotu/gophermart/internal/model"
)

type GetOrdersResponse struct {
	Number     string            `json:"number"`
	Status     model.OrderStatus `json:"status"`
	Accrual    float64           `json:"accrual,omitempty"`
	UploadedAt string            `json:"uploaded_at"`
}
