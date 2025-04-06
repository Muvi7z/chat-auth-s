package user

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/gen/api/user_v1"
	converter2 "github.com/Muvi7z/chat-auth-s/internal/converter"
	"github.com/Muvi7z/chat-auth-s/internal/logger"
	"github.com/Muvi7z/chat-auth-s/internal/sys"
	"github.com/Muvi7z/chat-auth-s/internal/sys/codes"
	"github.com/Muvi7z/chat-auth-s/internal/sys/validate"
	validate2 "github.com/Muvi7z/chat-auth-s/validate"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func (s *Implementation) Get(ctx context.Context, req *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	err := validate.Validate(
		ctx,
		validate2.ValidateId(req.GetId()),
		validate2.OtherValidateID(req.GetId()),
	)
	if err != nil {
		return nil, err
	}

	if req.GetId() > 100 {
		return nil, sys.NewCommonError("id must be less than 100", codes.ResourceExhausted)
	}

	logger.Info("Get user", zap.Int64("id", req.GetId()))

	span, ctx := opentracing.StartSpanFromContext(ctx, "GetUser")
	defer span.Finish()

	span.SetTag("id", req.GetId())

	u, err := s.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return converter2.ToUserGetResponseFromUser(u), nil
}
