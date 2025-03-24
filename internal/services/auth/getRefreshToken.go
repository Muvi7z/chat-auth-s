package auth

import (
	"context"
	"errors"
	"github.com/Muvi7z/chat-auth-s/internal/model"
	"github.com/Muvi7z/chat-auth-s/internal/utils"
	"time"
)

func (s *Service) GetRefreshToken(ctx context.Context, token string, secret string, duration int32) (string, error) {
	refreshTokenSecretByte := []byte(secret)

	claims, err := utils.VerifyToken(token, refreshTokenSecretByte)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	//TODO Можем слазать в базу или в кэш за доп данными пользователя

	refreshToken, err := utils.GenerateToken(model.UserInfo{
		Username: claims.Username,
		Role:     claims.Role,
	}, refreshTokenSecretByte, time.Duration(duration)*time.Minute)

	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
