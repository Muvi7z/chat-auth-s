-- +goose Up
create table "users" (
                      id serial primary key,
                      name text not null,
                      email text not null,
                      password text not null,
                      created_at timestamp default now(),
                      updated_at timestamp
);

-- +goose Down
drop table "users";
