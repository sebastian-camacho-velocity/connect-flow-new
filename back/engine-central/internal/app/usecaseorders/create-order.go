package usecaseorders

import (
	"context"
	"engine-central/internal/domain"
	"engine-central/internal/domain/dtos"
)

func (s *OrderUseCase) CreateOrder(ctx context.Context, req dtos.CreateOrder) (dtos.Order, error) {

	orderBrokerReq := dtos.CreateOrderReq{}
	if req.OrderNumber != "" {
		orderBrokerReq.OrderNumber = &req.OrderNumber
	}
	if req.ExternalOrderID != "" {
		orderBrokerReq.ExternalOrderID = &req.ExternalOrderID
	}

	orderBrokerReq.BusinessID = req.BusinessID
	orderBrokerReq.IntegrationTypeId = &req.IntegrationTypeID
	orderBrokerReq.ExternalIntegrationID = &req.ExternalIntegrationID
	orderBrokerReq.TotalShipment = req.TotalShipment
	orderBrokerReq.OrderTypeID = req.OrderTypeID

	if req.TrackingNumber != "" {
		orderBrokerReq.TrackingNumber = &req.TrackingNumber
	}
	orderBrokerReq.PaymentMethodId = req.PaymentMethodID // TODO: auto calculate in base of cod_total
	orderBrokerReq.PaymentType = req.PaymentTypeRaw
	if req.PaymentTypeID != 0 {
		orderBrokerReq.PaymentTypeId = &req.PaymentTypeID
	}
	if req.CountryNameRaw != "" { // TODO: build a normalizer for countries
		countryID := domain.NormalizeCountry(req.CountryNameRaw)
		orderBrokerReq.CountryId = &countryID
	} else if req.CountryID != 0 {
		orderBrokerReq.CountryId = &req.CountryID
	} else {
		defaultCountry := 1
		orderBrokerReq.CountryId = &defaultCountry
	}

	if req.WarehouseID != 0 {
		orderBrokerReq.WarehouseId = &req.WarehouseID
	}
	if req.ExtraData != nil {
		orderBrokerReq.ExtraData = req.ExtraData
	}

	orderBrokerReq.Customer = dtos.CustomerOrderReq{
		FullName:          req.Customer.FullName,
		MobilePhoneNumber: req.Customer.MobilePhoneNumber,
		DocumentTypeId:    req.Customer.DocumentTypeId,
		Dni:               req.Customer.Dni,
		Email:             req.Customer.Email, // TODO: check if is necessart create generic info
	}
	orderBrokerReq.Shipping = dtos.ShippingOrderReq{
		Country:           req.Shipping.Country,
		State:             req.Shipping.State,
		City:              req.Shipping.City,
		Address:           req.Shipping.Address,
		AddressLine:       req.Shipping.AddressLine,
		MobilePhoneNumber: req.Shipping.MobilePhoneNumber,
		FullName:          req.Shipping.FullName,
		Lat:               req.Shipping.Lat,
		Lng:               req.Shipping.Lng,
		Zip:               req.Shipping.Zip,
		CityDaneId:        req.Shipping.CityDaneId,
	}
	if req.OriginShipping != (dtos.CreateShippingOrder{}) {
		orderBrokerReq.OriginShipping = &dtos.ShippingOrderReq{
			Country:           req.OriginShipping.Country,
			State:             req.OriginShipping.State,
			City:              req.OriginShipping.City,
			Address:           req.OriginShipping.Address,
			AddressLine:       req.OriginShipping.AddressLine,
			MobilePhoneNumber: req.OriginShipping.MobilePhoneNumber,
			FullName:          req.OriginShipping.FullName,
			Lat:               req.OriginShipping.Lat,
			Lng:               req.OriginShipping.Lng,
			Zip:               req.OriginShipping.Zip,
			CityDaneId:        req.OriginShipping.CityDaneId,
		}
	}

	var products []dtos.ProductOrderReq
	for _, p := range req.Products {
		products = append(products, dtos.ProductOrderReq{
			ProductID:  p.ProductID,
			Sku:        p.Sku,
			ExternalId: p.ExternalId,
			// Ean:               p.Ean,  // TODO
			Name:              p.Name,
			Notes:             p.Notes,
			Large:             p.Large,
			Width:             p.Width,
			Weight:            p.Weight,
			Height:            p.Height,
			MeasurementUnitId: p.MeasurementUnitId,
			Description:       p.Description,
			Quantity:          p.Quantity,
			Price:             p.Price,
			Discount:          p.Discount,
			Tax:               p.Tax,
			Items:             nil, // TODO
			IsCustomKit:       p.IsCustomKit,
			Active:            p.Active,
		})
	}
	orderBrokerReq.OrderStatusID = req.OrderStatusID
	orderBrokerReq.Products = products
	orderBrokerReq.Notes = req.Notes
	if req.CodTotal != 0 {
		orderBrokerReq.CodTotal = &req.CodTotal
	}
	orderBrokerReq.DeliveryDate = req.DeliveryDate
	if req.Coupon != "" {
		orderBrokerReq.Coupon = &req.Coupon
	}
	orderBrokerReq.Discount = req.Discount
	if req.Total != 0 {
		orderBrokerReq.Total = &req.Total
	}
	orderBrokerReq.Boxes = req.Boxes

	_, err := s.orderBroker.CreateOrder(ctx, orderBrokerReq)
	if err != nil {
		s.log.Error(ctx).Err(err).Msg("failed to create order")
		return dtos.Order{}, err
	}

	return dtos.Order{}, err
}
