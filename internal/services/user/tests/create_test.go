package tests

import (
	"context"
	"fmt"
	"github.com/Muvi7z/chat-auth-s/internal/model"
	"github.com/Muvi7z/chat-auth-s/internal/repository"
	repoMocks "github.com/Muvi7z/chat-auth-s/internal/repository/mocks"
	"github.com/Muvi7z/chat-auth-s/internal/services/user"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_Create(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		req *model.User
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = "TName"
		email    = "TEmail"
		password = "TPassword"
		role     = 0

		req = &model.User{
			Name:     name,
			Email:    email,
			Password: password,
			Role:     int32(role),
		}
		res           = id
		convertedUser = &model.User{
			Name:     name,
			Email:    email,
			Password: password,
			Role:     int32(role),
		}
		serviceErr = fmt.Errorf("service error")
	)
	t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
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
			want: 0,
			err:  serviceErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, convertedUser).Return(0, serviceErr)
				return mock
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepositoryMock := tt.userRepositoryMock(mc)
			api := user.NewService(userRepositoryMock, nil)

			resCreate, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resCreate)
		})
	}
}
