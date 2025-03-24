package access

import (
	"context"
	"errors"
	"github.com/Muvi7z/chat-auth-s/internal/utils"
)

func (s *Service) Check(ctx context.Context, accessToken, accessTokenSecretKey, endpoint string) error {
	claims, err := utils.VerifyToken(accessToken, []byte(accessTokenSecretKey))
	if err != nil {
		return errors.New("access token is invalid")
	}

	role, ok := s.accessibleRoles[endpoint]
	if !ok {
		return nil
	}

	if role == claims.Role {
		return nil
	}

	return errors.New("access denied")
}

func (s *Service) getAccessibleRoles(ctx context.Context) (map[string]int32, error) {
	if s.accessibleRoles == nil {

	}

	return s.accessibleRoles, nil
}
