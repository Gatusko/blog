package data

import (
	"database/sql"
	"github.com/google/uuid"
	"log"
	"time"
)

type Feed struct {
	Id          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Url         string
	User        uuid.UUID
	LastFetched sql.NullTime
}

type FeedJSON struct {
	Id          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Url         string
	User        uuid.UUID
	LastFetched *time.Time `json:"last_fetched,omitempty"`
}

type FeedAndFeedFollow struct {
	Feed       FeedJSON    `json:"feed"`
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
	err = tx.QueryRow(query, args...).Scan(&feed.Id, &feed.CreatedAt, &feed.UpdatedAt, &feed.Name, &feed.Url, &feed.User, &feed.LastFetched)
	if err != nil {
		return feedAndFollow, err
	}

	query = `INSERT INTO feed_follows(id,created_at,updated_at,user_id,feed_id)
			values($1,$2,$3,$4,$5)
			returning *
			`
	log.Print("Here I am")
	argsFeedFollow := []any{uuid.New(), time.Now(), time.Now(), feed.User, feed.Id}
	err = tx.QueryRow(query, argsFeedFollow...).Scan(&feedAndFollow.FeedFollow.Id, &feedAndFollow.FeedFollow.CreatedAt, &feedAndFollow.FeedFollow.UpdatedAt, &feedAndFollow.FeedFollow.UserId, &feedAndFollow.FeedFollow.FeedId)
	if err != nil {
		return feedAndFollow, err
	}
	err = tx.Commit()
	if err != nil {
		return feedAndFollow, err
	}
	feedJson := feedDatabaseToFeedJson(*feed)
	feedAndFollow.Feed = feedJson
	return feedAndFollow, nil
}

func (fm *FeedModel) GetAll() ([]FeedJSON, error) {
	query := `SELECT * FROM feeds`
	rows, err := fm.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	feeds := []FeedJSON{}
	for rows.Next() {
		feed := Feed{}
		err = rows.Scan(&feed.Id, &feed.CreatedAt, &feed.UpdatedAt, &feed.Name, &feed.Url, &feed.User, &feed.LastFetched)
		if err != nil {
			return feeds, err
		}
		feedJSON := feedDatabaseToFeedJson(feed)
		feeds = append(feeds, feedJSON)
	}
	return feeds, nil
}

func (fm *FeedModel) GetNextFeedsToFetch(limit int) ([]FeedJSON, error) {
	query := `SELECT * FROM feeds
			  order by last_fetched_at NULLS FIRST
			  limit $1`
	rows, err := fm.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	feeds := []FeedJSON{}
	for rows.Next() {
		feed := Feed{}
		err = rows.Scan(&feed.Id, &feed.CreatedAt, &feed.UpdatedAt, &feed.Name, &feed.Url, &feed.User, &feed.LastFetched)
		if err != nil {
			return feeds, err
		}
		feedJSON := feedDatabaseToFeedJson(feed)
		feeds = append(feeds, feedJSON)
	}
	return feeds, nil
}

func (fm *FeedModel) MarkFeedFetched(id uuid.UUID) error {
	query := `UPDATE feeds SET 
                 updated_at = $1,
			     last_fetched_at = $1
			     where id=$2`
	result, err := fm.DB.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}
	log.Println("Update Succes:", result)
	return nil
}

func (fm *FeedModel) GetFeedById(id uuid.UUID) (Feed, error) {
	query := `SELECT * FROM feeds where id= $1`
	feed := Feed{}
	err := fm.DB.QueryRow(query, id).Scan(&feed.Id, &feed.CreatedAt, &feed.UpdatedAt, &feed.Name, &feed.Url, &feed.User, &feed.LastFetched)
	if err != nil {
		return feed, err
	}
	return feed, nil
}

func feedDatabaseToFeedJson(feed Feed) FeedJSON {
	if !feed.LastFetched.Valid {
		return FeedJSON{
			Id:        feed.Id,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			Name:      feed.Name,
			Url:       feed.Url,
			User:      feed.User,
		}
	}
	return FeedJSON{
		Id:          feed.Id,
		CreatedAt:   feed.CreatedAt,
		UpdatedAt:   feed.UpdatedAt,
		Name:        feed.Name,
		Url:         feed.Url,
		User:        feed.User,
		LastFetched: &feed.LastFetched.Time,
	}
}
