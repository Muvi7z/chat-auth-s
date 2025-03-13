package converter

import (
	"github.com/Muvi7z/chat-auth-s/gen/api/user_v1"
	"github.com/Muvi7z/chat-auth-s/internal/repository/user/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserGetResponseFromRepo(user *model.User) *user_v1.GetResponse {

	return &user_v1.GetResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      0,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func ToUserCreateRequestFromRepo(user *model.User) *user_v1.CreateRequest {
	return &user_v1.CreateRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     0,
	}
}
