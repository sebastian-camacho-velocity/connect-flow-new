package env

import (
	"context"
	"engine-central/internal/infra/shared/log"
	"fmt"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type IConfig interface {
	Get(key string) string
}

type config struct {
	values map[string]string
	logger log.ILogger
}

func New(logger log.ILogger) (IConfig, error) {
	_ = godotenv.Load()

	cfg := &Config{}
	missing := []string{}
	values := make(map[string]string)

	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("env")
		if tag == "" {
			continue
		}
		parts := splitTag(tag)
		key := parts[0]
		required := len(parts) > 1 && parts[1] == "required"
		val := os.Getenv(key)
		if val == "" && required {
			missing = append(missing, key)
		}
		values[key] = val
	}

	if len(missing) > 0 {
		logger.Error(context.Background()).
			Strs("missing_env_vars", missing).
			Msg("Faltan variables de entorno obligatorias")
		return nil, fmt.Errorf("faltan variables de entorno obligatorias: %v", missing)
	}
	return &config{values: values, logger: logger}, nil
}

// Get retorna el valor de una variable de entorno cargada
func (c *config) Get(key string) string {
	return c.values[key]
}

// Config solo se usa internamente para reflexión
// No debe ser accedido directamente fuera de este paquete

type Config struct {
	AppEnv             string `env:"APP_ENV,required"`
	ApiPort            string `env:"API_PORT,required"`
	LogLevel           string `env:"LOG_LEVEL,required"`
	JwtSecret          string `env:"JWT_SECRET,required"`
	ViteApiUrl         string `env:"VITE_API_URL,required"`
	PublicUrl          string `env:"PUBLIC_URL,required"`
	DbTimezone         string `env:"DB_TIMEZONE,required"`
	NatsHost           string `env:"NATS_HOST,required"`
	NatsPort           string `env:"NATS_PORT,required"`
	NatsUser           string `env:"NATS_USER,required"`
	NatsPass           string `env:"NATS_PASS,required"`
	DbHost             string `env:"DB_HOST,required"`
	DbUser             string `env:"DB_USER,required"`
	DbPass             string `env:"DB_PASS,required"`
	DbPort             string `env:"DB_PORT,required"`
	DbName             string `env:"DB_NAME,required"`
	DbLogLevel         string `env:"DB_LOG_LEVEL,required"`
	S3Bucket           string `env:"S3_BUCKET,required"`
	S3Region           string `env:"S3_REGION,required"`
	S3Key              string `env:"S3_KEY,required"`
	S3Secret           string `env:"S3_SECRET,required"`
	EmailHost          string `env:"EMAIL_HOST,required"`
	EmailPort          string `env:"EMAIL_PORT,required"`
	EmailUser          string `env:"EMAIL_USER,required"`
	EmailPass          string `env:"EMAIL_PASS,required"`
	EmailFrom          string `env:"EMAIL_FROM,required"`
	InventoryUrl       string `env:"INVENTORY_URL,required"`
	OrderBrokerUrl     string `env:"ORDER_BROKER_URL,required"`
	ApiUrl             string `env:"API_URL,required"`
	ShopifyKey         string `env:"SHOPIFY_KEY,required"`
	ShopifySecret      string `env:"SHOPIFY_SECRET,required"`
	BsaleUrl           string `env:"BSALE_URL,required"`
	BsaleCurrierUrl    string `env:"BSALE_CURRIER_URL,required"`
	ParisUrl           string `env:"PARIS_URL,required"`
	ViteConnectFlowUrl string `env:"VITE_CONNECT_FLOW_URL,required"`
	GrpcPort           string `env:"GRPC_PORT"`
}

func splitTag(tag string) []string {
	// Usamos SplitN para dividir solo en la primera coma
	parts := make([]string, 0, 2)
	for i, c := range tag {
		if c == ',' {
			parts = append(parts, tag[:i], tag[i+1:])
			return parts
		}
	}
	return []string{tag}
}
