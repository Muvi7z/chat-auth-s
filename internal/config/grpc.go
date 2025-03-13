package config

import "net"

type grpcConfig struct {
	host string
	port string
}

func NewGRPCConfig() (GRPCConfig, error) {
	return nil, nil
}

func (c *grpcConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
