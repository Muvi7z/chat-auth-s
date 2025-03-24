package auth

import (
	"context"
	"errors"
	"github.com/Muvi7z/chat-auth-s/internal/model"
	"github.com/Muvi7z/chat-auth-s/internal/utils"
	"time"
)

func (s *Service) Login(ctx context.Context, user *model.UserInfo, secret string, duration int32) (string, error) {
	//TODO Получение данных из базы или кэша

	refreshToken, err := utils.GenerateToken(*user, []byte(secret), time.Duration(duration)*time.Minute)
	if err != nil {
		return "", errors.New("failed to generate refresh token")
	}

	return refreshToken, nil
}
