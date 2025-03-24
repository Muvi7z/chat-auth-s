package services

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/internal/model"
)

type UserService interface {
	Get(context context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.User) (int64, error)
}

type AuthService interface {
	Login(ctx context.Context, user *model.UserInfo, secret string, duration int32) (string, error)
	GetRefreshToken(ctx context.Context, token string, secret string, duration int32) (string, error)
	GetAccessToken(ctx context.Context, tokenRefresh, refreshSecretKey, accessTokenSecretKey string, accessTokenExp int32) (string, error)
}

type AccessService interface {
	Check(ctx context.Context, accessToken, accessTokenSecretKey, endpoint string) error
}
