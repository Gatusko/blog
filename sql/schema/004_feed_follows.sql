-- +goose Up
CREATE TABLE feed_follows(
    id uuid PRIMARY KEY ,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    user_id uuid NOT NULL ,
    feed_id uuid NOT NULL,
    CONSTRAINT FK_USER_ID FOREIGN KEY(user_id)
    references users(id),
    CONSTRAINT FK_FEED_ID FOREIGN KEY(feed_id)
    REFERENCES feeds(id)
);
-- +goose Down
DROP TABLE feed_follows;