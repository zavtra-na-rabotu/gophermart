package dto

type CreateWithdrawalRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type GetWithdrawalsResponse struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
