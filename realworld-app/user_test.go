package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/istonikula/realworld-go/realworld-app/internal/config"
	appDb "github.com/istonikula/realworld-go/realworld-app/internal/db"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/rest"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/rest/apitest"
	domain "github.com/istonikula/realworld-go/realworld-domain"
	"github.com/istonikula/realworld-go/realworld-domain/test-support/fixture"
	"github.com/istonikula/realworld-go/realworld-domain/test-support/stub"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

type TestUser rest.UserRegistration

var testUser = TestUser{
	Email:    "foo@bar.com",
	Username: "foo",
	Password: "plain",
}

var userFactory = fixture.UserFactory{Auth: stub.UserStub.Auth}

func TestUsers(t *testing.T) {
	t.Run("register", func(t *testing.T) {
		db, cfg := setup()
		defer deleteUsers(db)
		client := apitest.Client{Router: router(db, cfg), Token: nil}

		r := client.Post("/api/users", rest.UserRegistration(testUser))

		require.Equal(t, http.StatusCreated, r.Code)

		act := readBody[rest.UserResponse](t, r).User
		require.Equal(t, testUser.Email, act.Email)
		require.Equal(t, testUser.Username, act.Username)
	})

	t.Run("cannot register already existing username", func(t *testing.T) {
		db, cfg := setup()
		defer deleteUsers(db)
		client := apitest.Client{Router: router(db, cfg), Token: nil}

		existing := userFactory.ValidRegistration(domain.UserRegistration(testUser))
		existing.Email = "unique." + testUser.Email
		saveUser(db, existing)

		r := client.Post("/api/users", rest.UserRegistration(testUser))

		require.Equal(t, http.StatusUnprocessableEntity, r.Code)
		require.Equal(t, "{\"error\":\"username already taken\"}", r.Body.String())
	})

	t.Run("cannot register already existing email", func(t *testing.T) {
		db, cfg := setup()
		defer deleteUsers(db)
		client := apitest.Client{Router: router(db, cfg), Token: nil}

		existing := userFactory.ValidRegistration(domain.UserRegistration(testUser))
		existing.Username = "unique"
		saveUser(db, existing)

		r := client.Post("/api/users", rest.UserRegistration(testUser))

		require.Equal(t, http.StatusUnprocessableEntity, r.Code)
		require.Equal(t, "{\"error\":\"email already taken\"}", r.Body.String())
	})

	t.Run("current user is resolved from token", func(t *testing.T) {
		db, cfg := setup()
		defer deleteUsers(db)

		registered := userFactory.ValidRegistration(domain.UserRegistration(testUser))
		saveUser(db, registered)

		client := apitest.Client{Router: router(db, cfg), Token: &registered.Token}
		r := client.Get("/api/user")
		require.Equal(t, http.StatusOK, r.Code)
		require.Equal(t, rest.User{
			Email:    registered.Email,
			Token:    registered.Token,
			Username: registered.Username,
		}, readBody[rest.UserResponse](t, r).User)
	})
}

func setup() (*sqlx.DB, *config.Config) {
	cfg := readConfig()
	return db(&cfg.DataSource), cfg
}

func saveUser(db *sqlx.DB, user *domain.ValidUserRegistration) {
	txMgr := appDb.TxMgr{DB: db}
	_ = txMgr.Write(func(tx *sqlx.Tx) error {
		var repo = appDb.UserRepo{Tx: tx}
		_, _ = repo.Create(user)
		return nil
	})
}

func deleteUsers(db *sqlx.DB) {
	db.MustExec("DELETE FROM users")
}

func readBody[B any](t *testing.T, r *httptest.ResponseRecorder) *B {
	var body B
	require.NoError(t, json.Unmarshal(r.Body.Bytes(), &body))
	return &body
}
