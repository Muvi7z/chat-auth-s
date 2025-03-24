package auth

import "github.com/Muvi7z/chat-auth-s/internal/services"

type Service struct {
}

func NewService() services.AuthService {
	return &Service{}
}
