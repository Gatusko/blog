package data

import "database/sql"

type Models struct {
	Users       UserModel
	Feeds       FeedModel
	FeedFollows FeedFollowModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:       UserModel{DB: db},
		Feeds:       FeedModel{DB: db},
		FeedFollows: FeedFollowModel{DB: db},
	}
}
