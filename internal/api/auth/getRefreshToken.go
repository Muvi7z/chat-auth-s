package auth

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/gen/api/auth_v1"
)

func (a *ImplementationAuth) GetRefreshToken(ctx context.Context, req *auth_v1.GetRefreshTokenRequest) (*auth_v1.GetRefreshTokenResponse, error) {
	refreshToken, err := a.authService.GetRefreshToken(ctx, req.GetRefreshToken(), a.refreshSecretKey, a.refreshTokenExp)
	if err != nil {
		return nil, err
	}

	return &auth_v1.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
}
