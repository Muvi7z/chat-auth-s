package user

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/gen/api/user_v1"
	converter2 "github.com/Muvi7z/chat-auth-s/internal/converter"
	"github.com/Muvi7z/chat-auth-s/internal/logger"
	"go.uber.org/zap"
)

func (s *Implementation) Get(ctx context.Context, req *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	logger.Info("Get user", zap.Int64("id", req.GetId()))

	u, err := s.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return converter2.ToUserGetResponseFromUser(u), nil
}
