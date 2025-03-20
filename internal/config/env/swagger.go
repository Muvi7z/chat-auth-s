package env

import (
	"fmt"
	"github.com/Muvi7z/chat-auth-s/internal/config"
	"net"
	"os"
)

const swaggerHostEnvName = "SWAGGER_HOST"
const swaggerPortEnvName = "SWAGGER_PORT"

type swaggerConfig struct {
	port string
	host string
}

func NewSwaggerConfig() (config.HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s is not found", httpHostEnvName)
	}

	port := os.Getenv(httpPortEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s is not found", httpPortEnvName)
	}

	return &swaggerConfig{port, host}, nil
}

func (h *swaggerConfig) Address() string {
	return net.JoinHostPort(h.host, h.port)
}
