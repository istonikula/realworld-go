package main

import (
	"encoding/json"
	appDb "github.com/istonikula/realworld-go/realworld-app/internal/db"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/rest"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/rest/apitest"
	domain "github.com/istonikula/realworld-go/realworld-domain"
	"github.com/istonikula/realworld-go/realworld-domain/test-support/fixture"
	"github.com/istonikula/realworld-go/realworld-domain/test-support/stub"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type TestUser rest.UserRegistration

var testUser = TestUser{
	Email:    "foo@bar.com",
	Username: "foo",
	Password: "plain",
}

func (u TestUser) User(token string, bio *string, image *string) rest.User {
	return rest.User{Email: u.Email, Username: u.Username, Token: token, Bio: bio, Image: image}
}

// TODO: real auth needed here
var userFactory = fixture.UserFactory{Auth: stub.UserStub.Auth}

func TestUsers(t *testing.T) {
	t.Run("register", func(t *testing.T) {
		db := setup()
		defer deleteUsers(db)
		var client = apitest.Client{Router: router(db), Token: nil}

		r := client.Post("/api/users", rest.UserRegistration(testUser))

		assert.Equal(t, http.StatusCreated, r.Code)

		var act rest.UserResponse
		assert.NoError(t, json.Unmarshal(r.Body.Bytes(), &act))
		exp := testUser.User("ignore", nil, nil)
		assertUserIgnoreToken(t, exp, act.User)
	})

	t.Run("cannot register already existing username", func(t *testing.T) {
		db := setup()
		defer deleteUsers(db)
		var client = apitest.Client{Router: router(db), Token: nil}

		var existing = domain.UserRegistration(testUser)
		existing.Email = "unique." + testUser.Email
		saveUser(db, &existing)

		r := client.Post("/api/users", rest.UserRegistration(testUser))

		assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
		assert.Equal(t, "{\"error\":\"username already taken\"}", r.Body.String())
	})

	t.Run("cannot register already existing email", func(t *testing.T) {
		db := setup()
		defer deleteUsers(db)
		var client = apitest.Client{Router: router(db), Token: nil}

		var existing = domain.UserRegistration(testUser)
		existing.Username = "unique"
		saveUser(db, &existing)

		r := client.Post("/api/users", rest.UserRegistration(testUser))

		assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
		assert.Equal(t, "{\"error\":\"email already taken\"}", r.Body.String())
	})
}

func setup() *sqlx.DB {
	cfg := readConfig()
	return db(&cfg.DataSource)
}

func saveUser(db *sqlx.DB, user *domain.UserRegistration) {
	txMgr := appDb.TxMgr{DB: db}
	_ = txMgr.Write(func(tx *sqlx.Tx) error {
		var repo = appDb.UserRepo{Tx: tx}
		_, _ = repo.Create(userFactory.ValidRegistration(user))
		return nil
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
