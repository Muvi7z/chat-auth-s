package user

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/gen/api/user_v1"
	"github.com/Muvi7z/chat-auth-s/internal/repository"
	"github.com/Muvi7z/chat-auth-s/internal/services"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	user_v1.UnimplementedUserV1Server
	UserRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) services.UserService {
	return &Service{UserRepository: userRepository}
}

func (s *Service) Update(ctx context.Context, reqest *user_v1.UpdateRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *Service) Delete(ctx context.Context, reqest *user_v1.DeleteRequest) (*emptypb.Empty, error) {
	return nil, nil
}
