package user

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/Muvi7z/chat-auth-s/internal/client/db"
	"github.com/Muvi7z/chat-auth-s/internal/model"
	"github.com/Muvi7z/chat-auth-s/internal/repository"
	"github.com/Muvi7z/chat-auth-s/internal/repository/user/converter"
	modalRepo "github.com/Muvi7z/chat-auth-s/internal/repository/user/model"
)

const (
	tableName      = "users"
	idColumn       = "id"
	nameColumn     = "name"
	emailColumn    = "email"
	passwordColumn = "password"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, request *model.User) (int64, error) {
	builder := squirrel.Insert(tableName).
		Columns(nameColumn, emailColumn, passwordColumn).
		Values(request.Name, request.Email, request.Password).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING \"id\"")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}
	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	builder := squirrel.Select(idColumn, nameColumn, emailColumn, passwordColumn).
		From(tableName).
		Where(squirrel.Eq{idColumn: id}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var user modalRepo.User

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), err
}
