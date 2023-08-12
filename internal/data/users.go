package data

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func NewUser(name string) (*User, error) {
	return &User{
		Id:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}, nil
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(user *User) error {
	query := `INSERT INTO users (id,created_at,update_at,name,apikey)
	VALUES($1,$2,$3,$4,encode(sha256(random()::text::bytea), 'hex'))
	RETURNING *
    `
	args := []any{user.Id, user.CreatedAt, user.UpdatedAt, user.Name}

	return u.DB.QueryRow(query, args...).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.ApiKey)
}

func (u *UserModel) Get(apiKey string) (*User, error) {
	query := `SELECT * FROM users where apikey = $1`
	args := []any{apiKey}
	user := User{}
	err := u.DB.QueryRow(query, args...).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.ApiKey)
	if err != nil {
		return &user, err
	}
	return &user, nil
}
