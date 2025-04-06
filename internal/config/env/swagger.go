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

func NewSwaggerConfig() (config.SwaggerConfig, error) {
	host := os.Getenv(swaggerHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s is not found", swaggerHostEnvName)
	}

	port := os.Getenv(swaggerPortEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s is not found", swaggerPortEnvName)
	}

	return &swaggerConfig{port, host}, nil
}

func (s *swaggerConfig) Address() string {
	return net.JoinHostPort(s.host, s.port)
}
