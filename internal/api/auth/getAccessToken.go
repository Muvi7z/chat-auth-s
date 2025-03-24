package auth

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/gen/api/auth_v1"
)

func (a *ImplementationAuth) GetAccessToken(ctx context.Context, request *auth_v1.GetAccessTokenRequest) (*auth_v1.GetAccessTokenResponse, error) {
	accessToken, err := a.authService.GetAccessToken(ctx, request.GetRefreshToken(), a.refreshSecretKey, a.accessSecretKey, a.accessTokenExp)
	if err != nil {
		return nil, err
	}

	return &auth_v1.GetAccessTokenResponse{AccessToken: accessToken}, nil
}
