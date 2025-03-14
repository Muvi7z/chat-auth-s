package converter

import (
	"github.com/Muvi7z/chat-auth-s/gen/api/user_v1"
	"github.com/Muvi7z/chat-auth-s/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserGetResponseFromUser(user *model.User) *user_v1.GetResponse {

	return &user_v1.GetResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user_v1.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func ToUserFromGetResponse(response *user_v1.GetResponse) *model.User {
	return &model.User{
		ID:        response.Id,
		Name:      response.Name,
		Email:     response.Email,
		Role:      int32(response.Role),
		CreatedAt: response.CreatedAt.AsTime(),
		UpdatedAt: response.UpdatedAt.AsTime(),
	}
}

func ToUserFromCreate(response *user_v1.CreateRequest) *model.User {
	return &model.User{
		Name:     response.Name,
		Email:    response.Email,
		Role:     int32(response.Role),
		Password: response.Password,
	}
}

func ToUserCreateRequestFromRepo(user *model.User) *user_v1.CreateRequest {
	return &user_v1.CreateRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     user_v1.Role(user.Role),
	}
}
