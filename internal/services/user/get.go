package user

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/internal/model"
)

func (s *Service) Get(context context.Context, id int64) (*model.User, error) {
	u, err := s.UserRepository.Get(context, id)
	if err != nil {
		//TODO Logging
		return nil, err
	}

	return u, nil
}
