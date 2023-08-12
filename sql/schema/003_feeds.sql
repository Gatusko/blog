-- +goose Up
CREATE TABLE feeds(
    id uuid PRIMARY KEY,
    created_at timestamp not null,
    updated_at timestamp not null,
    name text not null,
    url text not null,
    user_id uuid not null,
    CONSTRAINT  FK_FEED_USER FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;