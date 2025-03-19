package user

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/gen/api/user_v1"
	converter2 "github.com/Muvi7z/chat-auth-s/internal/converter"
	"log"
)

func (s *Implementation) Create(ctx context.Context, req *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	id, err := s.userService.Create(ctx, converter2.ToUserFromCreate(req))
	if err != nil {
		log.Printf("error to create user: %s", err.Error())
		return nil, err
	}

	log.Printf("Create User ID: %d", id)

	return &user_v1.CreateResponse{
		Id: id,
	}, nil
}
