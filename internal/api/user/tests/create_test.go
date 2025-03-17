package tests

import (
	"context"
	"fmt"
	"github.com/Muvi7z/chat-auth-s/gen/api/user_v1"
	"github.com/Muvi7z/chat-auth-s/internal/api/user"
	"github.com/Muvi7z/chat-auth-s/internal/model"
	"github.com/Muvi7z/chat-auth-s/internal/services"
	mocks2 "github.com/Muvi7z/chat-auth-s/internal/services/mocks"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_Create(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) services.UserService

	type args struct {
		ctx context.Context
		req *user_v1.CreateRequest
	}
	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = "TName"
		email    = "TEmail"
		password = "TPassword"
		role     = 0

		req = &user_v1.CreateRequest{
			Name:     name,
			Email:    email,
			Password: password,
			Role:     user_v1.Role(role),
		}
		res = &user_v1.CreateResponse{
			Id: id,
		}
		convertedUser = &model.User{
			Name:     name,
			Email:    email,
			Password: password,
			Role:     int32(role),
		}
		serviceErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *user_v1.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) services.UserService {
				mock := mocks2.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, convertedUser).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) services.UserService {
				mock := mocks2.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, convertedUser).Return(0, serviceErr)
				return mock
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userServiceMock := tt.userServiceMock(mc)
			api := user.NewImplementation(userServiceMock)

			resCreate, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resCreate)
		})
	}
}
