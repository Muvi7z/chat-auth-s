package interceptor

import (
	"context"
	"google.golang.org/grpc"
)

type validator interface {
	Validate() error
}

func ValidateInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if val, ok := req.(validator); ok {
		err = val.Validate()
		if err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}
