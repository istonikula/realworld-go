package db

import (
	"fmt"
	domain "github.com/istonikula/realworld-go/realworld-domain"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type UserRepo struct {
	DB *sqlx.DB
}

func (r UserRepo) Create(reg *domain.ValidUserRegistration) (*domain.User, error) {
	q := table("users").Insert(
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

func (r UserRepo) ExistsByUsername(username string) bool {
	return false
}

func (r UserRepo) ExistsByEmail(email string) bool {
	return false
}

type table string

func (t table) Insert(cols ...string) string {
	markerSql := "$1"
	if len(cols) > 1 {
		for i := range cols[1:] {
			markerSql += ", $" + strconv.Itoa(i+2)
		}
	}
	return "INSERT INTO " + string(t) + " (" + strings.Join(cols, ", ") + ") VALUES (" + markerSql + ")"
}
