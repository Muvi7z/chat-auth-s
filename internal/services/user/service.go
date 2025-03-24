package user

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/internal/client/db"
	"github.com/Muvi7z/chat-auth-s/internal/model"
	"github.com/Muvi7z/chat-auth-s/internal/repository"
	"github.com/Muvi7z/chat-auth-s/internal/services"
)

type Service struct {
	UserRepository repository.UserRepository
	txManager      db.TxManager
}

type serv struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository, txManager db.TxManager) services.UserService {
	return &Service{UserRepository: userRepository, txManager: txManager}
}

func NewMockService(deps ...interface{}) services.UserService {
	srv := serv{}

	for _, dep := range deps {
		switch s := dep.(type) {
		case repository.UserRepository:
			srv.userRepository = s

		}
	}

	return &srv
}

func (s serv) Get(context context.Context, id int64) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s serv) Create(ctx context.Context, user *model.User) (int64, error) {
	//TODO implement me
	panic("implement me")
}
