package mapper

import (
	"engine-central/internal/domain"
	"engine-central/internal/infra/primary/grpc/order/proto"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func MapCreateOrderRequestToDTO(req *proto.CreateOrderRequest) domain.CreateOrder {
	return domain.CreateOrder{
		OrderNumber:           req.OrderNumber,
		ExternalOrderID:       req.ExternalOrderId,
		BusinessID:            req.BusinessId,
		IntegrationTypeID:     int(req.IntegrationTypeId),
		ExternalIntegrationID: int(req.ExternalIntegrationId),
		OrderTypeID:           uint(req.OrderTypeId),
		TotalShipment:         req.TotalShipment,
		TrackingNumber:        req.TrackingNumber,
		PaymentMethodID:       int(req.PaymentMethodId),
		PaymentTypeRaw:        req.PaymentTypeRaw,
		PaymentTypeID:         int(req.PaymentTypeId),
		CountryID:             int(req.CountryId),
		CountryNameRaw:        req.CountryNameRaw,
		WarehouseID:           int(req.WarehouseId),
		ExtraData:             req.ExtraData.AsMap(),
		Customer:              mapCustomerOrder(req.Customer),
		Shipping:              mapShippingOrder(req.Shipping),
		OriginShipping:        mapShippingOrder(req.OriginShipping),
		Products:              mapProducts(req.Items),
		Notes:                 req.Notes,
		CodTotal:              req.CodTotal,
		DeliveryDate:          protoTimestampToTimePtr(req.DeliveryDate),
		Coupon:                req.Coupon,
		Discount:              req.Discount,
		Total:                 req.Total,
		Boxes:                 int(req.Boxes),
		IsLastMile:            req.IsLastMile,
		Invoiceable:           req.Invoiceable,
		OrderStatusID:         int(req.OrderStatusId),
	}
}

func mapCustomerOrder(c *proto.CreateCustomerOrder) domain.CreateCustomerOrder {
	if c == nil {
		return domain.CreateCustomerOrder{}
	}
	return domain.CreateCustomerOrder{
		FullName:          c.FullName,
		MobilePhoneNumber: c.MobilePhoneNumber,
		DocumentTypeId:    int(c.DocumentTypeId),
		Dni:               c.Dni,
		Email:             c.Email,
	}
}

func mapShippingOrder(s *proto.CreateShippingOrder) domain.CreateShippingOrder {
	if s == nil {
		return domain.CreateShippingOrder{}
	}
	return domain.CreateShippingOrder{
		Country:           s.Country,
		State:             s.State,
		City:              s.City,
		Address:           s.Address,
		AddressLine:       s.AddressLine,
		MobilePhoneNumber: s.MobilePhoneNumber,
		FullName:          s.FullName,
		Zip:               s.Zip,
		Lat:               doubleWrapperToPtr(s.Lat),
		Lng:               doubleWrapperToPtr(s.Lng),
		CityDaneId:        int32WrapperToPtr(s.CityDaneId),
	}
}

func mapProducts(items []*proto.ProductOrderReq) []domain.ProductOrderReq {
	products := make([]domain.ProductOrderReq, 0, len(items))
	for _, p := range items {
		products = append(products, mapProduct(p))
	}
	return products
}

func mapProduct(p *proto.ProductOrderReq) domain.ProductOrderReq {
	if p == nil {
		return domain.ProductOrderReq{}
	}
	return domain.ProductOrderReq{
		ProductID:         stringWrapperToPtr(p.ProductId),
		Sku:               stringWrapperToPtr(p.Sku),
		ExternalId:        stringWrapperToPtr(p.ExternalId),
		Ean:               stringWrapperToPtr(p.Ean),
		Name:              p.Name,
		Notes:             stringWrapperToPtr(p.Notes),
		Large:             doubleWrapperToPtr(p.Large),
		Width:             doubleWrapperToPtr(p.Width),
		Weight:            doubleWrapperToPtr(p.Weight),
		Height:            doubleWrapperToPtr(p.Height),
		MeasurementUnitId: p.MeasurementUnitId,
		Description:       p.Description,
		Quantity:          int(p.Quantity),
		Price:             p.Price,
		Discount:          p.Discount,
		Tax:               doubleWrapperToPtr(p.Tax),
		Items:             mapProducts(p.Items),
		IsCustomKit:       p.IsCustomKit,
		Active:            p.Active,
	}
}

// Utilidades para punteros y fechas
func strToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func floatToPtr(f float64) *float64 {
	return &f
}

func intToPtr(i int32) *int {
	v := int(i)
	return &v
}

func protoTimestampToTimePtr(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

func int32WrapperToPtr(w *wrapperspb.Int32Value) *int {
	if w == nil {
		return nil
	}
	v := int(w.Value)
	return &v
}

func doubleWrapperToPtr(w *wrapperspb.DoubleValue) *float64 {
	if w == nil {
		return nil
	}
	v := w.Value
	return &v
}

func stringWrapperToPtr(w *wrapperspb.StringValue) *string {
	if w == nil {
		return nil
	}
	v := w.Value
	return &v
}
