package access

import (
	"context"
	"errors"
	"github.com/Muvi7z/chat-auth-s/gen/api/access_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

const (
	authPrefix = "Bearer "
)

func (a *ImplementationAccess) Check(ctx context.Context, req *access_v1.CheckRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return nil, errors.New("authorization header is not provided")
	}

	if !strings.HasSuffix(authHeader[0], authPrefix) {
		return nil, errors.New("invalid authorization header")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	//TODO клбюч
	err := a.accessService.Check(ctx, accessToken, "", req.GetEndpointAddress())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
