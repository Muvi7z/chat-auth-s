package app

import (
	"context"
	"github.com/Muvi7z/chat-auth-s/internal/api/user"
	"github.com/Muvi7z/chat-auth-s/internal/client/db"
	"github.com/Muvi7z/chat-auth-s/internal/client/db/pg"
	"github.com/Muvi7z/chat-auth-s/internal/client/db/transaction"
	"github.com/Muvi7z/chat-auth-s/internal/closer"
	"github.com/Muvi7z/chat-auth-s/internal/config"
	"github.com/Muvi7z/chat-auth-s/internal/config/env"
	"github.com/Muvi7z/chat-auth-s/internal/repository"
	user2 "github.com/Muvi7z/chat-auth-s/internal/repository/user"
	"github.com/Muvi7z/chat-auth-s/internal/services"
	user3 "github.com/Muvi7z/chat-auth-s/internal/services/user"
	"log"
)

type serviceProvider struct {
	userService    services.UserService
	userRepository repository.UserRepository
	pgConfig       config.PGConfig
	grpcConfig     config.GRPCConfig
	httpConfig     config.HTTPConfig
	service        *user.Implementation
	dbClient       db.Client
	txManager      db.TxManager
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatal(err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatal(err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatal(err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatal(err)
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHttpConfig()
		if err != nil {
			log.Fatal(err)
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = user2.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) services.UserService {
	if s.userService == nil {
		s.userService = user3.NewService(s.UserRepository(ctx), s.TxManager(ctx))
	}

	return s.userService
}

func (s *serviceProvider) UserServer(ctx context.Context) *user.Implementation {
	if s.service == nil {
		s.service = user.NewImplementation(s.UserService(ctx))
	}

	return s.service
}
