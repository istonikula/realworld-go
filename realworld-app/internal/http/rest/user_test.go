package rest_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/istonikula/realworld-go/realworld-app/internal/boot"
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

func (u TestUser) WithUsername(s string) TestUser {
	u.Username = s
	return u
}
func (u TestUser) WithEmail(s string) TestUser {
	u.Email = s
	return u
}
func (u TestUser) WithPassword(s string) TestUser {
	u.Password = s
	return u
}
func (u TestUser) Reg() *rest.UserRegistration {
	addressable := rest.UserRegistration(u)
	return &addressable
}
func (u TestUser) Login() *rest.Login {
	return &rest.Login{Email: u.Email, Password: u.Password}
}

var testUser = TestUser{
	Username: "foo",
	Email:    "foo@bar.com",
	Password: "plain",
}

var userFactory = fixture.UserFactory{Auth: stub.UserStub.Auth}

func TestUsers(t *testing.T) {
	t.Run("register and login", func(t *testing.T) {
		db, cfg := setup()
		defer deleteUsers(db)
		client := apitest.Client{Router: rest.Router(cfg), Token: nil}

		r := client.Post("/api/users", testUser.Reg())
		require.Equal(t, http.StatusCreated, r.Code)
		registered := readBody[rest.UserResponse](t, r).User
		require.Equal(t, testUser.Email, registered.Email)
		require.Equal(t, testUser.Username, registered.Username)
		require.NotNil(t, registered.Token)

		r = client.Post("/api/users/login", testUser.Login())
		require.Equal(t, http.StatusOK, r.Code)
		loggedIn := readBody[rest.UserResponse](t, r).User
		require.Equal(t, rest.User{
			Email:    registered.Email,
			Token:    registered.Token,
			Username: registered.Username,
		}, loggedIn)

		r = client.Get("/api/user")
		require.Equal(t, http.StatusUnauthorized, r.Code)

		client.Token = &loggedIn.Token
		r = client.Get("/api/user")
		require.Equal(t, http.StatusOK, r.Code)
		require.Equal(t, rest.User{
			Email:    loggedIn.Email,
			Token:    loggedIn.Token,
			Username: loggedIn.Username,
		}, readBody[rest.UserResponse](t, r).User)
	})

	t.Run("register: already existing username", func(t *testing.T) {
		db, cfg := setup()
		defer deleteUsers(db)
		client := apitest.Client{Router: rest.Router(cfg), Token: nil}

		existing := userFactory.ValidRegistration(domain.UserRegistration(testUser))
		existing.Email = "unique." + testUser.Email
		saveUser(db, existing)

		r := client.Post("/api/users", rest.UserRegistration(testUser))

		require.Equal(t, http.StatusConflict, r.Code)
		require.Equal(t, "{\"error\":\"username already taken\"}", r.Body.String())
	})

	t.Run("register: already existing email", func(t *testing.T) {
		db, cfg := setup()
		defer deleteUsers(db)
		client := apitest.Client{Router: rest.Router(cfg), Token: nil}

		existing := userFactory.ValidRegistration(domain.UserRegistration(testUser))
		existing.Username = "unique"
		saveUser(db, existing)

		r := client.Post("/api/users", rest.UserRegistration(testUser))

		require.Equal(t, http.StatusConflict, r.Code)
		require.Equal(t, "{\"error\":\"email already taken\"}", r.Body.String())
	})

	t.Run("register: validation", func(t *testing.T) {
		tcs := map[string]struct {
			payload *rest.UserRegistration
			want    string
		}{
			"missing payload":  {nil, "email: cannot be blank; password: cannot be blank; username: cannot be blank."},
			"missing username": {testUser.WithUsername("").Reg(), "username: cannot be blank."},
			"missing email":    {testUser.WithEmail("").Reg(), "email: cannot be blank."},
			"invalid email":    {testUser.WithEmail("invalid").Reg(), "email: must be a valid email address."},
			"missing password": {testUser.WithPassword("").Reg(), "password: cannot be blank."},
		}
		for name, tc := range tcs {
			t.Run(name, func(t *testing.T) {
				db, cfg := setup()
				defer deleteUsers(db)
				client := apitest.Client{Router: rest.Router(cfg), Token: nil}

				r := client.Post("/api/users", tc.payload)
				require.Equal(t, http.StatusUnprocessableEntity, r.Code)
				require.Equal(t, fmt.Sprintf("{\"error\":\"%s\"}", tc.want), r.Body.String())
			})
		}
	})

	t.Run("register: unexpected exception", func(t *testing.T) {
		db, cfg := setup()
		defer deleteUsers(db)

		userRepo := func(tx *sqlx.Tx) appDb.UserRepoOps {
			return &apitest.MockUserRepo{
				UserRepo:          &appDb.UserRepo{Tx: tx},
				MockExistsByEmail: func(email string) (bool, error) { return false, errors.New("BOOM!") },
			}
		}

		client := apitest.Client{Router: rest.Router(cfg, boot.WithUserRepo(userRepo)), Token: nil}

		r := client.Post("/api/users", testUser.Reg())
		require.Equal(t, http.StatusInternalServerError, r.Code)
		require.Equal(t, "{\"error\":\"BOOM!\"}", r.Body.String())
	})

	t.Run("login: validation", func(t *testing.T) {
		tcs := map[string]struct {
			payload *rest.Login
			want    string
		}{
			"missing payload":  {nil, "email: cannot be blank; password: cannot be blank."},
			"missing email":    {testUser.WithEmail("").Login(), "email: cannot be blank."},
			"invalid email":    {testUser.WithEmail("invalid").Login(), "email: must be a valid email address."},
			"missing password": {testUser.WithPassword("").Login(), "password: cannot be blank."},
		}
		for name, tc := range tcs {
			t.Run(name, func(t *testing.T) {
				db, cfg := setup()
				defer deleteUsers(db)
				client := apitest.Client{Router: rest.Router(cfg), Token: nil}

				r := client.Post("/api/users/login", tc.payload)
				require.Equal(t, http.StatusUnprocessableEntity, r.Code)
				require.Equal(t, fmt.Sprintf("{\"error\":\"%s\"}", tc.want), r.Body.String())
			})
		}
	})

	t.Run("login: invalid password", func(t *testing.T) {
		db, cfg := setup()
		defer deleteUsers(db)
		client := apitest.Client{Router: rest.Router(cfg), Token: nil}

		saveUser(db, userFactory.ValidRegistration(domain.UserRegistration(testUser)))

		login := testUser.Login()
		login.Password = "invalid"
		r := client.Post("/api/users/login", login)

		require.Equal(t, http.StatusUnauthorized, r.Code)
	})

	t.Run("current user: invalid token", func(t *testing.T) {
		db, cfg := setup()
		defer deleteUsers(db)
		client := apitest.Client{Router: rest.Router(cfg), Token: nil}

		saveUser(db, userFactory.ValidRegistration(domain.UserRegistration(testUser)))

		token := "invalid"
		client.Token = &token
		r := client.Get("/api/user")
		require.Equal(t, http.StatusUnauthorized, r.Code)
	})

	t.Run("current user: user not found", func(t *testing.T) {
		db, cfg := setup()
		defer deleteUsers(db)

		token := userFactory.ValidRegistration(domain.UserRegistration(testUser)).Token
		client := apitest.Client{Router: rest.Router(cfg), Token: &token}

		r := client.Get("/api/user")
		require.Equal(t, http.StatusUnauthorized, r.Code)
	})
}

func setup() (*sqlx.DB, *config.Config) {
	cfg := boot.ReadConfig("../../../config.yml")
	return boot.MustConnect(&cfg.DataSource), cfg
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
