package user

import (
	"github.com/Muvi7z/chat-auth-s/gen/api/user_v1"
	"github.com/Muvi7z/chat-auth-s/internal/services"
)

type Implementation struct {
	user_v1.UnimplementedUserV1Server
	userService services.UserService
}

func NewImplementation(userService services.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
