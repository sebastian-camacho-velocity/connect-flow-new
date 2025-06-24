package usecaseorders

import (
	"context"
	"engine-central/internal/domain"
	"engine-central/internal/infra/shared/log"
)

type IUseCaseOrders interface {
	CreateOrder(ctx context.Context, req domain.CreateOrder) (domain.Order, error)
}

type OrderUseCase struct {
	orderBroker domain.IOrderBroker
	log         log.ILogger
}

func NewOrderUseCase(orderBroker domain.IOrderBroker, log log.ILogger) IUseCaseOrders {
	return &OrderUseCase{
		orderBroker: orderBroker,
		log:         log,
	}
}
