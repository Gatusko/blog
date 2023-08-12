-- +goose Up
alter table users
add column apikey varchar(64)  not null DEFAULT encode(sha256(random()::text::bytea), 'hex');


-- +goose Down
alter table users
drop column apikey;