package integration

import (
	"github.com/go-resty/resty/v2"
	"github.com/zavtra-na-rabotu/gophermart/internal/dto"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

const defaultTooManyRequestsWaitTime int = 60

type AccrualClient struct {
	client *resty.Client
}

func NewAccrualClient(url string) *AccrualClient {
	return &AccrualClient{
		client: resty.New().SetBaseURL(url),
	}
}

func (c *AccrualClient) ProcessOrder(orderNumber string) (*dto.AccrualOrderResponse, error) {
	response, err := c.client.R().
		SetResult(&dto.AccrualOrderResponse{}).
		Get("/api/orders/" + orderNumber)

	if err != nil {
		zap.L().Error("Failed to process order", zap.String("orderNumber", orderNumber), zap.Error(err))
		return nil, err
	}

	if response.StatusCode() == http.StatusNoContent {
		zap.L().Error("Order not registered", zap.String("orderNumber", orderNumber))
		return nil, err
	}

	if response.StatusCode() == http.StatusTooManyRequests {
		retryAfterHeader := response.Header().Get("Retry-After")

		timeToWait := defaultTooManyRequestsWaitTime

		if retryAfterHeader != "" {
			retryAfter, err := strconv.Atoi(retryAfterHeader)
			if err == nil {
				timeToWait = retryAfter
			} else {
				zap.L().Error("Invalid Retry-After header", zap.String("orderNumber", orderNumber), zap.Int("Retry-After", retryAfter), zap.Error(err))
			}
		}

		zap.L().Info("Too many requests, waiting before retry", zap.Int("Retry-After (seconds)", timeToWait))
		time.Sleep(time.Duration(timeToWait) * time.Second)
		return c.ProcessOrder(orderNumber)
	}

	if response.StatusCode() == http.StatusInternalServerError {
		zap.L().Error("Failed to process order", zap.String("orderNumber", orderNumber))
		return nil, err
	}

	return response.Result().(*dto.AccrualOrderResponse), nil
}
