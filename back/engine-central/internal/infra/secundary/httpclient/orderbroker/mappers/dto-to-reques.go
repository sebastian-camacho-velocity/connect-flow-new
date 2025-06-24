package mappers

import (
	"engine-central/internal/domain"
	"engine-central/internal/infra/secundary/httpclient/orderbroker/request"
	"engine-central/internal/infra/secundary/httpclient/orderbroker/response"
)

func ToOrderBrokerRequest(in domain.CreateOrderReq) request.CreateOrderReq {
	return request.CreateOrderReq{
		BusinessID:                 in.BusinessID,
		ExternalOrderID:            in.ExternalOrderID,
		ExternalIntegrationID:      in.ExternalIntegrationID,
		OrderNumber:                in.OrderNumber,
		IntegrationTypeId:          in.IntegrationTypeId,
		OrderTypeID:                in.OrderTypeID,
		OrderType:                  in.OrderType,
		TotalShipment:              in.TotalShipment,
		TrackingNumber:             in.TrackingNumber,
		PaymentMethodId:            in.PaymentMethodId,
		IsPaid:                     in.IsPaid,
		CountryId:                  in.CountryId,
		CityDaneId:                 in.CityDaneId,
		WarehouseId:                in.WarehouseId,
		ExtraData:                  in.ExtraData,
		Customer:                   toOrderBrokerCustomer(in.Customer),
		Shipping:                   toOrderBrokerShipping(in.Shipping),
		OriginShipping:             toOrderBrokerShippingPtr(in.OriginShipping),
		Products:                   toOrderBrokerProducts(in.Products),
		Notes:                      in.Notes,
		CodTotal:                   in.CodTotal,
		DeliveryDate:               in.DeliveryDate,
		Coupon:                     in.Coupon,
		SiigoInvoiceId:             in.SiigoInvoiceId,
		DeliveryProviderTypeZoneId: in.DeliveryProviderTypeZoneId,
		Discount:                   in.Discount,
		PaymentType:                in.PaymentType,
		PaymentTypeId:              in.PaymentTypeId,
		Total:                      in.Total,
		Boxes:                      in.Boxes,
		OrderStatusID:              in.OrderStatusID,
	}
}

func FromOrderBrokerRequest(in request.CreateOrderReq) domain.CreateOrder {
	return domain.CreateOrder{
		BusinessID:            in.BusinessID,
		ExternalOrderID:       derefStr(in.ExternalOrderID),
		ExternalIntegrationID: derefInt(in.ExternalIntegrationID),
		OrderNumber:           derefStr(in.OrderNumber),
		IntegrationTypeID:     derefInt(in.IntegrationTypeId),
		OrderTypeID:           in.OrderTypeID,
		TotalShipment:         in.TotalShipment,
		TrackingNumber:        derefStr(in.TrackingNumber),
		PaymentMethodID:       in.PaymentMethodId,
		CountryID:             derefInt(in.CountryId),
		WarehouseID:           derefInt(in.WarehouseId),
		ExtraData:             in.ExtraData,
		Customer:              customerOrderReqToCreate(fromOrderBrokerCustomer(in.Customer)),
		Shipping:              shippingOrderReqToCreate(fromOrderBrokerShipping(in.Shipping)),
		OriginShipping:        shippingOrderReqToCreate(fromOrderBrokerShippingPtr(in.OriginShipping)),
		Products:              fromOrderBrokerProducts(in.Products),
		Notes:                 in.Notes,
		CodTotal:              derefFloat64(in.CodTotal),
		DeliveryDate:          in.DeliveryDate,
		Coupon:                derefStr(in.Coupon),
		Discount:              in.Discount,
		Total:                 derefFloat64(in.Total),
		Boxes:                 in.Boxes,
		OrderStatusID:         in.OrderStatusID,
	}
}

func FromOrderBrokerResponse(res response.CreateOrderRes) domain.Order {
	return domain.Order{
		OrderId: res.OrderId,
	}
}

// Helpers para punteros y conversiones
func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
func intPtrOrNil(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}
func float64PtrOrNil(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}
func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
func derefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}
func derefFloat64(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func toOrderBrokerCustomer(in domain.CustomerOrderReq) request.CustomerOrderReq {
	return request.CustomerOrderReq{
		FullName:          in.FullName,
		MobilePhoneNumber: in.MobilePhoneNumber,
		DocumentTypeId:    in.DocumentTypeId,
		Dni:               in.Dni,
		Email:             in.Email,
	}
}
func fromOrderBrokerCustomer(in request.CustomerOrderReq) domain.CustomerOrderReq {
	return domain.CustomerOrderReq{
		FullName:          in.FullName,
		MobilePhoneNumber: in.MobilePhoneNumber,
		DocumentTypeId:    in.DocumentTypeId,
		Dni:               in.Dni,
		Email:             in.Email,
	}
}
func toOrderBrokerShipping(in domain.ShippingOrderReq) request.ShippingOrderReq {
	return request.ShippingOrderReq{
		Country:           in.Country,
		State:             in.State,
		City:              in.City,
		Address:           in.Address,
		AddressLine:       in.AddressLine,
		MobilePhoneNumber: in.MobilePhoneNumber,
		FullName:          in.FullName,
		Zip:               in.Zip,
		Lat:               in.Lat,
		Lng:               in.Lng,
		CityDaneId:        in.CityDaneId,
	}
}
func fromOrderBrokerShipping(in request.ShippingOrderReq) domain.ShippingOrderReq {
	return domain.ShippingOrderReq{
		Country:           in.Country,
		State:             in.State,
		City:              in.City,
		Address:           in.Address,
		AddressLine:       in.AddressLine,
		MobilePhoneNumber: in.MobilePhoneNumber,
		FullName:          in.FullName,
		Zip:               in.Zip,
		Lat:               in.Lat,
		Lng:               in.Lng,
		CityDaneId:        in.CityDaneId,
	}
}
func toOrderBrokerShippingPtr(in *domain.ShippingOrderReq) *request.ShippingOrderReq {
	if in == nil {
		return nil
	}
	v := toOrderBrokerShipping(*in)
	return &v
}
func fromOrderBrokerShippingPtr(in *request.ShippingOrderReq) domain.ShippingOrderReq {
	if in == nil {
		return domain.ShippingOrderReq{}
	}
	return fromOrderBrokerShipping(*in)
}
func toOrderBrokerProducts(in []domain.ProductOrderReq) []request.ProductOrderReq {
	out := make([]request.ProductOrderReq, len(in))
	for i, p := range in {
		out[i] = request.ProductOrderReq{
			ProductID:         p.ProductID,
			Sku:               p.Sku,
			ExternalId:        p.ExternalId,
			Name:              p.Name,
			Notes:             p.Notes,
			Large:             p.Large,
			Width:             p.Width,
			Height:            p.Height,
			Weight:            p.Weight,
			MeasurementUnitId: p.MeasurementUnitId,
			Description:       p.Description,
			Quantity:          p.Quantity,
			Price:             p.Price,
			Discount:          p.Discount,
			Tax:               p.Tax,
			Items:             toOrderBrokerProducts(p.Items),
			IsCustomKit:       p.IsCustomKit,
			Active:            p.Active,
		}
	}
	return out
}
func fromOrderBrokerProducts(in []request.ProductOrderReq) []domain.ProductOrderReq {
	out := make([]domain.ProductOrderReq, len(in))
	for i, p := range in {
		out[i] = domain.ProductOrderReq{
			ProductID:         p.ProductID,
			Sku:               p.Sku,
			ExternalId:        p.ExternalId,
			Name:              p.Name,
			Notes:             p.Notes,
			Large:             p.Large,
			Width:             p.Width,
			Height:            p.Height,
			Weight:            p.Weight,
			MeasurementUnitId: p.MeasurementUnitId,
			Description:       p.Description,
			Quantity:          p.Quantity,
			Price:             p.Price,
			Discount:          p.Discount,
			Tax:               p.Tax,
			Items:             fromOrderBrokerProducts(p.Items),
			IsCustomKit:       p.IsCustomKit,
			Active:            p.Active,
		}
	}
	return out
}

// Función auxiliar para convertir dtos.ShippingOrderReq a dtos.CreateShippingOrder
func shippingOrderReqToCreate(in domain.ShippingOrderReq) domain.CreateShippingOrder {
	return domain.CreateShippingOrder{
		Country:           in.Country,
		State:             in.State,
		City:              in.City,
		Address:           in.Address,
		AddressLine:       in.AddressLine,
		MobilePhoneNumber: in.MobilePhoneNumber,
		FullName:          in.FullName,
		Zip:               in.Zip,
		Lat:               in.Lat,
		Lng:               in.Lng,
		CityDaneId:        in.CityDaneId,
	}
}

// Función auxiliar para convertir dtos.CustomerOrderReq a dtos.CreateCustomerOrder
func customerOrderReqToCreate(in domain.CustomerOrderReq) domain.CreateCustomerOrder {
	return domain.CreateCustomerOrder{
		FullName:          in.FullName,
		MobilePhoneNumber: in.MobilePhoneNumber,
		DocumentTypeId:    in.DocumentTypeId,
		Dni:               in.Dni,
		Email:             in.Email,
	}
}
