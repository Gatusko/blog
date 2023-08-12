package data

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Feed struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Url       string
	User      uuid.UUID
}

type FeedAndFeedFollow struct {
	Feed       Feed        `json:"feed"`
	FeedFollow FeedFollows `json:"feed_follow"`
}

type FeedModel struct {
	DB *sql.DB
}

func NewFeed(name string, url string, userId uuid.UUID) (*Feed, error) {
	return &Feed{
		Id:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		User:      userId,
	}, nil
}

func (fm *FeedModel) Insert(feed *Feed) (FeedAndFeedFollow, error) {
	tx, err := fm.DB.Begin()
	feedAndFollow := FeedAndFeedFollow{}
	if err != nil {
		return feedAndFollow, err
	}
	defer tx.Rollback()
	query := `INSERT INTO feeds(id,created_at,updated_at,name,url,user_id)
			  values ($1,$2,$3,$4,$5,$6)
			  returning *
			  `
	args := []any{feed.Id, feed.CreatedAt, feed.UpdatedAt, feed.Name, feed.Url, feed.User}
	err = tx.QueryRow(query, args...).Scan(&feedAndFollow.Feed.Id, &feedAndFollow.Feed.CreatedAt, &feedAndFollow.FeedFollow.UpdatedAt, &feedAndFollow.Feed.Name, &feedAndFollow.Feed.Url, &feedAndFollow.Feed.User)
	if err != nil {
		return feedAndFollow, err
	}

	query = `INSERT INTO feed_follows(id,created_at,updated_at,user_id,feed_id)
			values($1,$2,$3,$4,$5)
			returning *
			`
	argsFeedFollow := []any{uuid.New(), time.Now(), time.Now(), feedAndFollow.Feed.User, feedAndFollow.Feed.Id}
	err = tx.QueryRow(query, argsFeedFollow...).Scan(&feedAndFollow.FeedFollow.Id, &feedAndFollow.FeedFollow.CreatedAt, &feedAndFollow.FeedFollow.UpdatedAt, &feedAndFollow.FeedFollow.UserId, &feedAndFollow.FeedFollow.FeedId)
	if err != nil {
		return feedAndFollow, err
	}
	err = tx.Commit()
	if err != nil {
		return feedAndFollow, err
	}
	return feedAndFollow, nil
}

func (fm *FeedModel) GetAll() ([]Feed, error) {
	query := `SELECT * FROM feeds`
	rows, err := fm.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	feeds := []Feed{}
	for rows.Next() {
		feed := Feed{}
		err = rows.Scan(&feed.Id, &feed.CreatedAt, &feed.UpdatedAt, &feed.Name, &feed.Url, &feed.User)
		if err != nil {
			return feeds, err
		}
		feeds = append(feeds, feed)
	}
	return feeds, nil
}
