package repository

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, request *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
}
