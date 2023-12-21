package db

import (
	"database/sql"
	"errors"
	"fmt"

	domain "github.com/istonikula/realworld-go/realworld-domain"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	Tx *sqlx.Tx
}

var tbl = table("users")

const selectCols = "id, email, token, username, bio, image"

func (r UserRepo) Create(reg *domain.ValidUserRegistration) (*domain.User, error) {
	q := tbl.insert(
		"id", "email", "token", "username", "password",
	) + " RETURNING " + selectCols
	stmt, err := r.Tx.Preparex(q)
	if err != nil {
		return nil, fmt.Errorf("UserRepo#Create prepare: %w", err)
	}

	var user domain.User
	err = stmt.QueryRowx(
		reg.Id, reg.Email, reg.Token, reg.Username, reg.EncryptedPassword,
	).StructScan(&user)
	if err != nil {
		return nil, fmt.Errorf("UserRepo#Create insert: %w", err)
	}
	return &user, nil
}

func (r UserRepo) ExistsByUsername(username string) (bool, error) {
	exists, err := tbl.queryIfExists(r.Tx, "username = $1", username)
	if err != nil {
		return false, fmt.Errorf("UserRepo#ExistsByUsername: %w", err)
	}
	return exists, nil
}

func (r UserRepo) ExistsByEmail(email string) (bool, error) {
	exists, err := tbl.queryIfExists(r.Tx, "email = $1", email)
	if err != nil {
		return false, fmt.Errorf("UserRepo#ExistsByEmail: %w", err)
	}
	return exists, nil
}

func (r UserRepo) FindById(id domain.UserId) (*domain.User, error) {
	stmt, err := r.Tx.Preparex(fmt.Sprintf("SELECT %v FROM %v WHERE id = $1", selectCols, tbl))
	if err != nil {
		return nil, fmt.Errorf("UserRepo#FindById: %w", err)
	}

	var user domain.User
	err = stmt.QueryRowx(id).StructScan(&user)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
