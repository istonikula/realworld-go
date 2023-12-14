package db

import (
	"fmt"
	domain "github.com/istonikula/realworld-go/realworld-domain"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	DB *sqlx.DB
}

var usersTbl = table("users")

func (r UserRepo) Create(reg *domain.ValidUserRegistration) (*domain.User, error) {
	q := usersTbl.insert(
		"id", "email", "token", "username", "password",
	) + " RETURNING id, email, token, username, bio, image"
	stmt, err := r.DB.Preparex(q)
	if err != nil {
		return nil, fmt.Errorf("UserRepo#Create: prepare failed: %w", err)
	}

	var user domain.User
	err = stmt.QueryRowx(
		reg.Id, reg.Email, reg.Token, reg.Username, reg.EncryptedPassword,
	).StructScan(&user)
	if err != nil {
		return nil, fmt.Errorf("UserRepo#Create: insert failed: %w", err)
	}
	return &user, nil
}

func (r UserRepo) ExistsByUsername(username string) (bool, error) {
	return usersTbl.queryIfExists(r.DB, "username=$1", username)
}

func (r UserRepo) ExistsByEmail(email string) (bool, error) {
	return usersTbl.queryIfExists(r.DB, "email = $1", email)
}
