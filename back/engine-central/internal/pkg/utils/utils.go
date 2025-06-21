package utils

import (
	"context"
	"encoding/json"

	"engine-central/internal/infra/shared/log"

	"github.com/labstack/echo/v4"
)

// LoadConfig helper function to load config from json.RawMessage
func LoadConfig[T any](config json.RawMessage) T {
	var c T
	json.Unmarshal(config, &c)
	return c
}

// logWithIntegrationTypeMiddleware is a middleware function for the Echo framework that logs the integration type ID.
// It adds the integration type ID to the request context for logging purposes.
func LogWithIntegrationTypeMiddleware(
	integrationTypeId int,
	log log.ILogger,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := log.With().Int("integration_type_id", integrationTypeId).Logger().WithContext(c.Request().Context())
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

// logWithIntegrationType adds the integration type ID to the logger context.
func LogWithIntegrationType(
	ctx context.Context,
	integrationTypeId int,
	log log.ILogger,
) context.Context {
	return log.With().Int("integration_type_id", integrationTypeId).Logger().WithContext(ctx)
}

// DecodePayload decodes a JSON-encoded payload into a specified type.
// The function takes a byte slice containing the JSON data and returns
// an instance of the specified type along with an error, if any occurred
// during the unmarshalling process.
func DecodePayload[T any](data []byte) (T, error) {
	var v T
	err := json.Unmarshal(data, &v)
	return v, err
}

// EncodePayload encodes the given data of any type into a JSON byte slice.
func EncodePayload[T any](data T) []byte {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return b
}

type OrderUpdatedPayload struct {
	OrderID       string `json:"orderId"`
	IntegrationID int    `json:"integrationId"`
}
