package user

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/Muvi7z/chat-auth-s/gen/api/user_v1"
	"github.com/Muvi7z/chat-auth-s/internal/repository"
	"github.com/Muvi7z/chat-auth-s/internal/repository/user/converter"
	"github.com/Muvi7z/chat-auth-s/internal/repository/user/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tableName      = "user"
	idColumn       = "id"
	nameColumn     = "name"
	emailColumn    = "email"
	passwordColumn = "password"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, request *user_v1.CreateRequest) (int64, error) {
	builder := squirrel.Insert(tableName).
		Columns("name", "email", "password").
		Values(request.Name, request.Email, request.Password).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}
	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*user_v1.GetResponse, error) {
	builder := squirrel.Select(idColumn, nameColumn, emailColumn, passwordColumn).
		From(tableName).
		Where(squirrel.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var user model.User

	err = r.db.QueryRow(ctx, query, args...).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return converter.ToUserGetResponseFromRepo(&user), err
}
