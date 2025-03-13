-- +goose Up
create table "user" (
                      id serial primary key,
                      name text not null,
                      email text not null,
                      password text not null,
                      created_at timestamp not null default now(),
                      updated_at timestamp
);

-- +goose Down
drop table "user";
