package repository

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/gen/api/user_v1"
)

type UserRepository interface {
	Create(ctx context.Context, request *user_v1.CreateRequest) (int64, error)
	Get(ctx context.Context, id int64) (*user_v1.GetResponse, error)
}
