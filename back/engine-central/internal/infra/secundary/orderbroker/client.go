package orderbroker

import (
	"bytes"
	"context"
	"encoding/json"
	"engine-central/internal/domain/dtos"
	"engine-central/internal/infra/secundary/orderbroker/mapperorderbroker"
	"engine-central/internal/infra/secundary/orderbroker/request"
	"engine-central/internal/infra/secundary/orderbroker/response"
	"engine-central/internal/infra/shared/env"
	"engine-central/internal/infra/shared/httpclient"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type OrderBroker interface {
	CreateOrder(ctx context.Context, req dtos.CreateOrderReq) (response.CreateOrderRes, error)
	ConfirmOrder(ctx context.Context, id string) error
	UploadFile(ctx context.Context, req request.UploadFileReq) error
}

type Client struct {
	http    *http.Client
	baseURL string
}

func NewClient(env env.IConfig) *Client {
	httpClient := httpclient.NewHTTPClient(httpclient.HTTPClientConfig{
		Timeout:            15 * time.Second,
		MaxIdleConns:       100,
		IdleConnTimeout:    90 * time.Second,
		DisableCompression: false,
	})
	return &Client{
		http:    httpClient,
		baseURL: env.Get("ORDER_BROKER_URL"),
	}
}

func (c *Client) CreateOrder(ctx context.Context, dtos dtos.CreateOrderReq) (response.CreateOrderRes, error) {
	var res response.CreateOrderRes
	req := mapperorderbroker.ToOrderBrokerRequest(dtos)
	err := c.postJSON(ctx, "/create", req, &res)
	return res, err
}

func (c *Client) ConfirmOrder(ctx context.Context, id string) error {
	body := map[string]string{"order_id": id}
	return c.postJSON(ctx, "/confirm", body, nil)
}

func (c *Client) UploadFile(ctx context.Context, req request.UploadFileReq) error {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	if err := addFile(writer, "file", req.FileName, req.FileReader); err != nil {
		return err
	}

	_ = writer.WriteField("order_id", req.OrderId)
	_ = writer.WriteField("note", req.Note)

	if err := writer.Close(); err != nil {
		return err
	}

	requestURL := c.baseURL + "/upload-file"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, buf)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	return c.doRequest(httpReq, nil)
}

// postJSON centraliza las peticiones POST JSON
func (c *Client) postJSON(ctx context.Context, path string, body any, res any) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return fmt.Errorf("encoding json: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, buf)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	return c.doRequest(httpReq, res)
}

// doRequest ejecuta una petición y decodifica la respuesta si se espera resultado
func (c *Client) doRequest(req *http.Request, res any) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("order broker request failed [%d]: %s - %s", resp.StatusCode, resp.Status, string(body))
	}

	if res != nil {
		if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}

// addFile ayuda a armar el multipart
func addFile(writer *multipart.Writer, fieldName, filename string, reader io.Reader) error {
	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		return fmt.Errorf("creating form file: %w", err)
	}
	if _, err := io.Copy(part, reader); err != nil {
		return fmt.Errorf("copying file: %w", err)
	}
	return nil
}
