package env

import (
	"fmt"
	"github.com/Muvi7z/chat-auth-s/internal/config"
	"net"
	"os"
)

const grpcHostEnvName = "GRPC_HOST"
const grpcPortEnvName = "GRPC_PORT"

type grpcConfig struct {
	host string
	port string
}

func NewGRPCConfig() (config.GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s is not found", grpcHostEnvName)
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s is not found", grpcPortEnvName)
	}

	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}

func (c *grpcConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
