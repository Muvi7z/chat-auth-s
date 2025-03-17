package user

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/internal/model"
	"log"
)

func (s *Service) Create(ctx context.Context, user *model.User) (int64, error) {
	var id int64
	id, err := s.UserRepository.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	//err = s.txManager.ReadCommited(ctx, func(ctx context.Context) error { //использование тразикции
	//	var errTx error
	//
	//	id, errTx = s.UserRepository.Create(ctx, user)
	//	if errTx != nil {
	//		return errTx
	//	}
	//
	//	return nil
	//})
	//if err != nil {
	//	return 0, err
	//}

	log.Printf("inserted user with id: %v", id)

	return id, nil
}
