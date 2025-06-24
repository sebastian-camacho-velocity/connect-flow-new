package orderbroker

import (
	"context"
	"engine-central/internal/domain"
	"engine-central/internal/infra/secundary/httpclient/orderbroker/mappers"
	"engine-central/internal/infra/secundary/httpclient/orderbroker/response"
)

func (c *Client) CreateOrder(ctx context.Context, dtos domain.CreateOrderReq) (domain.Order, error) {
	var res response.CreateOrderRes
	req := mappers.ToOrderBrokerRequest(dtos)

	err := c.postJSON(ctx, "/create", req, &res)
	if err != nil {
		return domain.Order{}, err
	}

	order := mappers.FromOrderBrokerResponse(res)
	return order, nil
}
