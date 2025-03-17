package user

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/gen/api/user_v1"
	"github.com/Muvi7z/chat-auth-s/internal/client/db"
	"github.com/Muvi7z/chat-auth-s/internal/repository"
	"github.com/Muvi7z/chat-auth-s/internal/services"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	UserRepository repository.UserRepository
	txManager      db.TxManager
}

func NewService(userRepository repository.UserRepository, txManager db.TxManager) services.UserService {
	return &Service{UserRepository: userRepository, txManager: txManager}
}

func (s *Service) Update(ctx context.Context, reqest *user_v1.UpdateRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *Service) Delete(ctx context.Context, reqest *user_v1.DeleteRequest) (*emptypb.Empty, error) {
	return nil, nil
}
