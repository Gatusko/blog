package data

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type FeedFollows struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    uuid.UUID `json:"user_id"`
	FeedId    uuid.UUID `json:"feed_id"`
}

type FeedFollowModel struct {
	DB *sql.DB
}

func NewFeedFollows(userId uuid.UUID, feedId uuid.UUID) (*FeedFollows, error) {
	return &FeedFollows{
		uuid.New(),
		time.Now(),
		time.Now(),
		userId,
		feedId,
	}, nil
}

func (feedFollowModel *FeedFollowModel) CreateFeedFollow(feedFollow FeedFollows) error {
	query := `INSERT INTO feed_follows(id,created_at,updated_at,user_id,feed_id)
			values($1,$2,$3,$4,$5)
			returning *
			`
	args := []any{feedFollow.FeedId, feedFollow.CreatedAt, feedFollow.UpdatedAt, feedFollow.UserId, feedFollow.FeedId}
	return feedFollowModel.DB.QueryRow(query, args...).Scan(&feedFollow.FeedId, &feedFollow.CreatedAt, &feedFollow.UpdatedAt, &feedFollow.UserId, &feedFollow.FeedId)
}

func (feedFollowModel *FeedFollowModel) GetAllFeedFollows(userId uuid.UUID) ([]FeedFollows, error) {
	query := `SELECT * FROM feed_follows
			WHERE user_id=$1	
		`
	rows, err := feedFollowModel.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	feedFollows := []FeedFollows{}
	for rows.Next() {
		feedFollow := FeedFollows{}
		err = rows.Scan(&feedFollow.Id, &feedFollow.CreatedAt, &feedFollow.UpdatedAt, &feedFollow.UserId, &feedFollow.FeedId)
		if err != nil {
			return feedFollows, err
		}
		feedFollows = append(feedFollows, feedFollow)
	}
	return feedFollows, nil
}

func (feedFollowModel *FeedFollowModel) DeleteFeedFollow(id uuid.UUID) error {
	query := `DELETE FROM feed_follows
			  where id = $1`

	result, err := feedFollowModel.DB.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("It didn't found the feed_follow")
	}
	return nil
}
