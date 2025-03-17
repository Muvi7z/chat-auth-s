package user

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/gen/api/user_v1"
	converter2 "github.com/Muvi7z/chat-auth-s/internal/converter"
)

func (s *Implementation) Get(ctx context.Context, req *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	u, err := s.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return converter2.ToUserGetResponseFromUser(u), nil
}
