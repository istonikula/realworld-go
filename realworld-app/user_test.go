package main

import (
	"encoding/json"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/rest"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/rest/apitest"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var testUser = struct {
	Email    string
	Username string
	Password string
}{
	Email:    "foo@bar.com",
	Username: "foo",
	Password: "plain",
}

func TestUsers(t *testing.T) {
	t.Run("register", func(t *testing.T) {
		var db = db()
		defer deleteUsers(db)
		var client = apitest.Client{Router: router(db), Token: nil}

		r := client.Post("/api/users", rest.UserRegistration{
			Email:    testUser.Email,
			Username: testUser.Username,
			Password: testUser.Password,
		})

		assert.Equal(t, http.StatusCreated, r.Code)

		var act rest.UserResponse
		assert.NoError(t, json.Unmarshal(r.Body.Bytes(), &act))
		exp := rest.User{Email: testUser.Email, Username: testUser.Username, Token: "ignore", Bio: nil, Image: nil}
		assertUserIgnoreToken(t, exp, act.User)
	})
}

func deleteUsers(db *sqlx.DB) {
	db.MustExec("DELETE FROM users")
}

func assertUserIgnoreToken(t *testing.T, exp, act rest.User) {
	exp.Token = "ignore"
	act.Token = "ignore"
	assert.Equal(t, exp, act)
}
