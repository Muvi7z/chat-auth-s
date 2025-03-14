package converter

import (
	"github.com/Muvi7z/chat-auth-s/gen/api/user_v1"
	"github.com/Muvi7z/chat-auth-s/internal/model"
	modalRepo "github.com/Muvi7z/chat-auth-s/internal/repository/user/model"
)

func ToUserFromRepo(user *modalRepo.User) *model.User {

	return &model.User{
		ID:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      0,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserCreateRequestFromRepo(user *modalRepo.User) *user_v1.CreateRequest {
	return &user_v1.CreateRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     0,
	}
}
