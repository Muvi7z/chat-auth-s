package services

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/internal/model"
)

type UserService interface {
	Get(context context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.User) (int64, error)
}
