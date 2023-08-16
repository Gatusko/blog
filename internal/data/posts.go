package data

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Post struct {
	Id          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	FeedId      uuid.UUID `json:"feed_id"`
}

func NewPost(feedId uuid.UUID, title string, url string, description string) *Post {
	return &Post{Id: uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       title,
		Url:         url,
		Description: description,
		PublishedAt: time.Now(),
		FeedId:      feedId}
}

type PostModel struct {
	DB *sql.DB
}

func (pm *PostModel) InsertPost(post *Post) error {
	query := `INSERT INTO posts (id, created_at ,updated_at ,title ,url ,description , published_at , feed_id)
	values ($1,$2,$3,$4,$5,$6,$7,$8)
	returning *`
	args := []any{post.Id, post.CreatedAt, post.UpdatedAt, post.Title, post.Url, post.Description, post.PublishedAt, post.FeedId}
	return pm.DB.QueryRow(query, args...).Scan(&post.Id, &post.CreatedAt, &post.Title, &post.UpdatedAt, &post.Url, &post.Description, &post.PublishedAt, &post.FeedId)
}

func (pm *PostModel) GetPostByUser(userId uuid.UUID, size int) ([]Post, error) {
	query := `SELECT p.id, p.created_at , p.updated_at , p.title , p.url ,p.description , p.published_at , p.feed_id 
			  FROM posts as p
			  left join feeds as f
			  on p.feed_id=f.id
			  where f.user_id=$1
			  order by p.created_at desc
			  limit $2`
	rows, err := pm.DB.Query(query, userId, size)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := []Post{}
	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.CreatedAt, &post.UpdatedAt, &post.Title, &post.Url, &post.Description, &post.PublishedAt, &post.FeedId)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, err
}
