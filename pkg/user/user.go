package user

import (
	"context"
	"database/sql"
	"errors"
	"redditclone/pkg/jwt"
	"redditclone/pkg/support"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	password string
}

func NewUser(username, password string) *User {
	return &User{
		Username: username,
		password: password,
	}
}

var (
	Unauthorized    = errors.New("unauthorized")
	NotFound        = errors.New("user not found")
	InvalidPassword = errors.New("invalid password")
	AlreadyExists   = errors.New("already exists")
)

func FromContext(ctx context.Context) (*User, error) {
	usr, ok := ctx.Value("user").(*User)
	if !ok {
		return nil, Unauthorized
	}
	return usr, nil
}

// --------------

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

var (
	ValidationUserNameRequired = errors.New("username is required")
	ValidationPasswordRequired = errors.New("password is required")
)

func (r *Repository) Create(u *User) (*User, error) {
	u.ID = support.GenerateID(24)
	if err := r.Validate(u); err != nil {
		return nil, err
	}
	var username string
	row := r.DB.QueryRow("SELECT username FROM users WHERE username = ?", u.Username)
	err := row.Scan(&username)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if username != "" {
		return nil, AlreadyExists
	}
	_, err = r.DB.Exec("INSERT INTO users(`id`, `username`, `password`) VALUES (?, ?, ?)", u.ID, u.Username, u.password)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *Repository) Authorize(in *User) (*User, error) {
	if err := r.Validate(in); err != nil {
		return nil, err
	}
	usr := &User{}
	row := r.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", in.Username)
	err := row.Scan(&usr.ID, &usr.Username, &usr.password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NotFound
		}
		return nil, err
	}
	if usr.password != in.password {
		return nil, InvalidPassword
	}
	return usr, nil
}

func (r *Repository) Validate(u *User) error {
	if u.Username == "" {
		return ValidationUserNameRequired
	}
	if u.password == "" {
		return ValidationPasswordRequired
	}
	return nil
}

func (r *Repository) GetByToken(tokenString string) (*User, error) {
	userInfo, err := jwt.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	usr := &User{}
	row := r.DB.QueryRow("SELECT id, username, password FROM users WHERE id=? LIMIT 1", userInfo.ID)
	err = row.Scan(&usr.ID, &usr.Username, &usr.password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NotFound
		}
		return nil, err
	}
	return usr, nil
}
