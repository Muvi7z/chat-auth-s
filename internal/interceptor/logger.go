package interceptor

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	now := time.Now()

	req, err := handler(ctx, req)
	if err != nil {
		logger.Error(err.Error(), zap.String("method", info.FullMethod), zap.Any("request", req))
	}

	logger.Info("request", zap.Any("req", req), zap.Duration("duration", time.Since(now)))

	return req, err
}
