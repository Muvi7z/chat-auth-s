package main

import (
	"context"
	"flag"
	"github.com/Muvi7z/chat-auth-s/gen/api/user_v1"
	user3 "github.com/Muvi7z/chat-auth-s/internal/api/user"
	"github.com/Muvi7z/chat-auth-s/internal/config"
	"github.com/Muvi7z/chat-auth-s/internal/config/env"
	user2 "github.com/Muvi7z/chat-auth-s/internal/repository/user"
	"github.com/Muvi7z/chat-auth-s/internal/services/user"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatal(err)
	}

	pgConfig, err := env.NewPGConfig()

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	userRepo := user2.NewRepository(pool)
	userService := user.NewService(userRepo)

	s := grpc.NewServer()
	reflection.Register(s)

	user_v1.RegisterUserV1Server(s, user3.NewService(userService))

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
