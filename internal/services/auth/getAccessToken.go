package auth

import (
	"context"
	"errors"
	"github.com/Muvi7z/chat-auth-s/internal/model"
	"github.com/Muvi7z/chat-auth-s/internal/utils"
	"time"
)

func (s *Service) GetAccessToken(ctx context.Context, tokenRefresh, secretKey, accessTokenSecretKey string, accessTokenExp int32) (string, error) {

	claims, err := utils.VerifyToken(tokenRefresh, []byte(secretKey))
	if err != nil {
		return "", errors.New("invalid refresh token")
	}
	//TODO Можем слазать в базу или в кэш за доп данными пользователя
	accessToken, err := utils.GenerateToken(model.UserInfo{
		Username: claims.Username,
		Role:     0,
	}, []byte(accessTokenSecretKey), time.Duration(accessTokenExp)*time.Minute)

	if err != nil {
		return "", err
	}

	return accessToken, nil
}
