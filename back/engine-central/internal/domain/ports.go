package domain

import (
	"context"
	"engine-central/internal/infra/secundary/httpclient/orderbroker/request"
)

type IOrderBroker interface {
	CreateOrder(ctx context.Context, req CreateOrderReq) (Order, error)
	ConfirmOrder(ctx context.Context, id string) error
	UploadFile(ctx context.Context, req request.UploadFileReq) error
}
