package usecaseorders

import (
	"engine-central/internal/infra/secundary/orderbroker"
	"engine-central/internal/infra/shared/log"
)

type OrderUseCase struct {
	orderBroker orderbroker.OrderBroker
	log         log.ILogger
}

func NewOrderUseCase(ob orderbroker.OrderBroker, log log.ILogger) *OrderUseCase {
	return &OrderUseCase{
		orderBroker: ob,
		log:         log,
	}
}
