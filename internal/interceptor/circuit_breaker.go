package interceptor

import (
	"context"
	"errors"
	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CircuitBreakerInterceptor struct {
	cb *gobreaker.CircuitBreaker[any]
}

func NewCircuitBreakerInterceptor(cb *gobreaker.CircuitBreaker[any]) *CircuitBreakerInterceptor {
	return &CircuitBreakerInterceptor{cb: cb}
}

func (c *CircuitBreakerInterceptor) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, err := c.cb.Execute(func() (interface{}, error) {
		return handler(ctx, req)
	})

	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			return nil, status.Error(codes.Unavailable, "server unavailable")
		}

		return nil, err
	}

	return res, nil
}
