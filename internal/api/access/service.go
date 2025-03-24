package access

import (
	"github.com/Muvi7z/chat-auth-s/gen/api/access_v1"
	"github.com/Muvi7z/chat-auth-s/internal/services"
)

type ImplementationAccess struct {
	access_v1.UnimplementedAccessV1Server
	accessService services.AccessService
}

func NewImplementationAccess(accessService services.AccessService) *ImplementationAccess {
	return &ImplementationAccess{
		accessService: accessService,
	}
}
