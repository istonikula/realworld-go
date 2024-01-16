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
	"google.golang.org/grpc"
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
		ctx := setup()
		defer ctx.teardown()

		client := proto.NewUsersClient(ctx.conn)

		r, _ := client.RegisterUser(ctx.Context, testUser.Reg())
		registered := r.User
		require.Equal(t, testUser.Email, registered.Email)
		require.Equal(t, testUser.Username, registered.Username)
		require.NotNil(t, registered.Token)

		r, _ = client.Login(ctx.Context, testUser.Login())
		loggedIn := r.User
		require.Equal(t, registered, loggedIn)
	})
}

type testCtx struct {
	context.Context
	db       *sqlx.DB
	cfg      *config.Config
	conn     *grpc.ClientConn
	teardown func()
}

func setup() *testCtx {
	ctx := context.Background()
	cfg := boot.ReadConfig("../../../../config.yml")
	db := boot.MustConnect(&cfg.DataSource)
	conn, cleanup := apitest.Server(ctx, server.Router(cfg))
	return &testCtx{
		Context: ctx,
		db:      db,
		cfg:     cfg,
		conn:    conn,
		teardown: func() {
			cleanup()
			deleteUsers(db)
		},
	}
}

// NOTE: same in rest
func deleteUsers(db *sqlx.DB) {
	db.MustExec("DELETE FROM users")
}
