package user

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/gen/api/user_v1"
	"github.com/Muvi7z/chat-auth-s/internal/repository"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

type Server struct {
	user_v1.UnimplementedUserV1Server
	UserRepository repository.UserRepository
}

func (s *Server) Get(context context.Context, request *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	u, err := s.UserRepository.Get(context, request.Id)
	if err != nil {
		//TODO Logging
		return nil, err
	}

	return u, nil
}

func (s *Server) Create(ctx context.Context, request *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	id, err := s.UserRepository.Create(ctx, request)
	if err != nil {
		return nil, err
	}

	log.Printf("inserted user with id: %v", id)

	return &user_v1.CreateResponse{
		Id: id,
	}, nil
}

func (s *Server) Update(ctx context.Context, reqest *user_v1.UpdateRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *Server) Delete(ctx context.Context, reqest *user_v1.DeleteRequest) (*emptypb.Empty, error) {
	return nil, nil
}
