-- +goose Up
CREATE TABLE posts(
    id uuid PRIMARY KEY,
    created_at timestamp not null,
    updated_at timestamp not null,
    title text not null,
    url text not null,
    description text not null,
    published_at timestamp not null,
    feed_id uuid not null,
    CONSTRAINT FK_FEED_ID FOREIGN KEY (feed_id)
    references feeds(id)
);


-- +goose Down
DROP TABLE posts;