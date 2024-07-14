package job

import (
	"github.com/zavtra-na-rabotu/gophermart/internal/integration"
	"github.com/zavtra-na-rabotu/gophermart/internal/service"
	"go.uber.org/zap"
)

type AccrualJob struct {
	accrualClient *integration.AccrualClient
	orderService  *service.OrderService
}

func NewAccrualJob(accrualClient *integration.AccrualClient, orderService *service.OrderService) *AccrualJob {
	return &AccrualJob{accrualClient: accrualClient, orderService: orderService}
}

func (j *AccrualJob) Start() {
	orders, err := j.orderService.GetAllNotTerminated()
	if err != nil {
		zap.L().Error("Cannot get orders to process", zap.Error(err))
		return
	}

	for _, order := range orders {
		accrualResponse, err := j.accrualClient.ProcessOrder(order.Number)
		if err != nil {
			zap.L().Error("Cannot process order", zap.Error(err))
			continue
		}

		if accrualResponse.Status == "REGISTERED" {
			continue
		}

		err = j.orderService.UpdateOrder(accrualResponse.Order, accrualResponse.Accrual, accrualResponse.Status)
		if err != nil {
			zap.L().Error("Cannot update order", zap.Error(err))
			continue
		}

		zap.L().Info(
			"Order processed",
			zap.String("order", accrualResponse.Order),
			zap.String("status", accrualResponse.Status),
			zap.Float64("accrual", accrualResponse.Accrual),
		)
	}
}
