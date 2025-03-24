package auth

import (
	"github.com/Muvi7z/chat-auth-s/gen/api/auth_v1"
	"github.com/Muvi7z/chat-auth-s/internal/services"
)

type ImplementationAuth struct {
	auth_v1.UnimplementedAuthV1Server
	authService      services.AuthService
	refreshTokenExp  int32
	accessTokenExp   int32
	refreshSecretKey string
	accessSecretKey  string
}

func NewImplementationAuth(authService services.AuthService, duration int32, secretKey string) *ImplementationAuth {
	return &ImplementationAuth{
		authService:      authService,
		refreshTokenExp:  duration,
		refreshSecretKey: secretKey,
	}
}
