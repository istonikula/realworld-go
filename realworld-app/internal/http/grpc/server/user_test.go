package server_test

import (
	"context"
	"testing"

	"github.com/istonikula/realworld-go/realworld-app/internal/boot"
	"github.com/istonikula/realworld-go/realworld-app/internal/config"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/grpc/proto"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/grpc/server"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/grpc/server/apitest"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

type TestUser struct {
	Username string
	Email    string
	Password string
}

func (u TestUser) Reg() *proto.UserRegistration {
	return &proto.UserRegistration{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
}
func (u TestUser) Login() *proto.LoginRequest {
	return &proto.LoginRequest{
		Email:    u.Email,
		Password: u.Password,
	}
}

var testUser = TestUser{
	Username: "foo",
	Email:    "foo@bar.com",
	Password: "plain",
}

func TestUsers(t *testing.T) {
	t.Run("register and login", func(t *testing.T) {
		ctx := context.Background()
		db, cfg := setup()
		defer deleteUsers(db)

		conn, cleanup := apitest.Server(ctx, server.Router(cfg))
		defer cleanup()
		client := proto.NewUsersClient(conn)

		r, _ := client.RegisterUser(ctx, testUser.Reg())
		registered := r.User
		require.Equal(t, testUser.Email, registered.Email)
		require.Equal(t, testUser.Username, registered.Username)
		require.NotNil(t, registered.Token)

		r, _ = client.Login(ctx, testUser.Login())
		loggedIn := r.User
		require.Equal(t, registered, loggedIn)
	})
}

// NOTE: same in rest (different path)
func setup() (*sqlx.DB, *config.Config) {
	cfg := boot.ReadConfig("../../../../config.yml")
	return boot.MustConnect(&cfg.DataSource), cfg
}

// NOTE: same in rest
func deleteUsers(db *sqlx.DB) {
	db.MustExec("DELETE FROM users")
}
