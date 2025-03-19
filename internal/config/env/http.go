package env

import (
	"fmt"
	"github.com/Muvi7z/chat-auth-s/internal/config"
	"net"
	"os"
)

const httpHostEnvName = "HTTP_HOST"
const httpPortEnvName = "HTTP_PORT"

type httpConfig struct {
	port string
	host string
}

func NewHttpConfig() (config.HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s is not found", httpHostEnvName)
	}

	port := os.Getenv(httpPortEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s is not found", httpPortEnvName)
	}

	return &httpConfig{port, host}, nil
}

func (h *httpConfig) Address() string {
	return net.JoinHostPort(h.host, h.port)
}
