package env

import (
	"errors"
	"github.com/Muvi7z/chat-auth-s/internal/config"
	"os"
)

type pgConfig struct {
	dsn string
}

func NewPGConfig() (config.PGConfig, error) {
	dsn := os.Getenv("PG_DNS")
	if len(dsn) == 0 {
		return nil, errors.New("PG_DNS environment variable not set")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

func (c *pgConfig) DSN() string {
	return c.dsn
}
