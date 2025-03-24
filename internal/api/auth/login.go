package auth

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/gen/api/auth_v1"
	"github.com/Muvi7z/chat-auth-s/internal/model"
)

func (a *ImplementationAuth) Login(ctx context.Context, login *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {
	refreshToken, err := a.authService.Login(ctx, &model.UserInfo{
		Username: login.Username,
		Password: login.Password,
	}, a.refreshSecretKey, a.refreshTokenExp)
	if err != nil {
		return nil, err
	}

	return &auth_v1.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
