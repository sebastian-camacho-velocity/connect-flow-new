package domain

import "time"

type CreateOrder struct {
	OrderNumber           string
	ExternalOrderID       string
	BusinessID            string
	IntegrationTypeID     int
	ExternalIntegrationID int
	OrderTypeID           uint
	TotalShipment         float64
	TrackingNumber        string
	PaymentMethodID       int
	PaymentTypeRaw        string
	PaymentTypeID         int
	CountryID             int
	CountryNameRaw        string
	WarehouseID           int
	ExtraData             any
	Customer              CreateCustomerOrder
	Shipping              CreateShippingOrder
	OriginShipping        CreateShippingOrder
	Products              []ProductOrderReq
	Notes                 []string
	CodTotal              float64
	DeliveryDate          *time.Time
	Coupon                string
	Discount              float64
	Total                 float64
	Boxes                 int
	IsLastMile            bool
	Invoiceable           bool
	OrderStatusID         int
}

type CreateCustomerOrder struct {
	FullName          string
	MobilePhoneNumber string
	DocumentTypeId    int
	Dni               string
	Email             string
}

type CreateShippingOrder struct {
	Country           string
	State             string
	City              string
	Address           string
	AddressLine       string
	MobilePhoneNumber string
	FullName          string
	Zip               string
	Lat               *float64
	Lng               *float64
	CityDaneId        *int
}

type ProductOrderReq struct {
	ProductID         *string
	Sku               *string
	ExternalId        *string
	Ean               *string
	Name              string
	Notes             *string
	Large             *float64
	Width             *float64
	Weight            *float64
	Height            *float64
	MeasurementUnitId int64
	Description       string
	Quantity          int
	Price             float64
	Discount          float64
	Tax               *float64
	Items             []ProductOrderReq
	IsCustomKit       bool
	Active            bool
}

type UpdateOrderStatusReq struct {
	OrderID       string
	OrderStatusID int
	Notes         string
	ExtraData     any
}

type Order struct {
	OrderId string
}

type CreateOrderReq struct {
	BusinessID                 string
	ExternalOrderID            *string
	ExternalIntegrationID      *int
	OrderNumber                *string
	IntegrationTypeId          *int
	OrderTypeID                uint
	OrderType                  string
	TotalShipment              float64
	TrackingNumber             *string
	PaymentMethodId            int
	IsPaid                     *bool
	CountryId                  *int
	CityDaneId                 *int
	WarehouseId                *int
	ExtraData                  any
	Customer                   CustomerOrderReq
	Shipping                   ShippingOrderReq
	OriginShipping             *ShippingOrderReq
	Products                   []ProductOrderReq
	Notes                      []string
	CodTotal                   *float64
	DeliveryDate               *time.Time
	Coupon                     *string
	SiigoInvoiceId             *string
	DeliveryProviderTypeZoneId *int
	Discount                   float64
	PaymentType                string
	PaymentTypeId              *int
	Total                      *float64
	Boxes                      int
	OrderStatusID              int
}

type CustomerOrderReq struct {
	FullName          string
	MobilePhoneNumber string
	DocumentTypeId    int
	Dni               string
	Email             string
}

type ShippingOrderReq struct {
	Country           string
	State             string
	City              string
	Address           string
	AddressLine       string
	MobilePhoneNumber string
	FullName          string
	Zip               string
	Lat               *float64
	Lng               *float64
	CityDaneId        *int
}

type CancelOrderReq struct {
	OrderId string
	Reason  string
	UserId  *int
}
