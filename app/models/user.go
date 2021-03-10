package models

import (
	"apisim/app/db"
	"context"
	"database/sql"
)

const (
	createUserSQL     = `insert into users (username, first_name, last_name, email, password_hash, created_at) values ($1, $2, $3, $4, $5, $6) returning id`
	getUsersSQL       = `select id, username, first_name, last_name, email, password_hash, created_at, updated_at from users`
	getUserByID       = getUsersSQL + ` where id=$1`
	getUserByUsername = getUsersSQL + ` where username=$1`
	updateUserSQL     = `update users set (username, first_name, last_name, email, updated_at) = ($1, $2, $3, $4, $5) where id = $6`
)

type (
	User struct {
		SequentialIdentifier
		Username     string `json:"username"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Email        string `json:"email"`
		PasswordHash string `json:"-"`
		Timestamps
	}
)

// func (u *User) String() string {
// 	return fmt.Sprintf("User(%s)", u.Username)
// }

func (u *User) GetByID(
	ctx context.Context,
	db db.SQLOperations,
	id int64,
) (*User, error) {
	var user User
	row := db.QueryRowContext(ctx, getUserByID, id)

	err := u.scanRowIntoUser(row, &user)
	return &user, err
}

func (u *User) GetByUsername(
	ctx context.Context,
	db db.SQLOperations,
	username string,
) (*User, error) {
	var user User
	row := db.QueryRowContext(ctx, getUserByUsername, username)

	err := u.scanRowIntoUser(row, &user)
	return &user, err
}

func (u *User) Save(
	ctx context.Context,
	db db.SQLOperations,
) error {
	u.Timestamps.Touch()

	var err error
	if u.IsNew() {
		err := db.QueryRowContext(
			ctx,
			createUserSQL,
			u.Username,
			u.FirstName,
			u.LastName,
			u.Email,
			u.PasswordHash,
			u.Timestamps.CreatedAt,
		).Scan(&u.ID)
		return err
	}
	_, err = db.ExecContext(
		ctx,
		updateUserSQL,
		u.Username,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Timestamps.UpdatedAt,
		u.ID,
	)
	return err
}

func (q *User) scanRowIntoUser(
	row *sql.Row,
	user *User,
) error {
	return row.Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}
