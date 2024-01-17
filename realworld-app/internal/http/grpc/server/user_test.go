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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TestUser struct {
	Username string
	Email    string
	Password string
}

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

		r, _ := client.RegisterUser(ctx, testUser.Reg())
		registered := r.User
		require.Equal(t, testUser.Email, registered.Email)
		require.Equal(t, testUser.Username, registered.Username)
		require.NotNil(t, registered.Token)

		r, _ = client.Login(ctx, testUser.Login())
		loggedIn := r.User
		require.Equal(t, registered, loggedIn)

		_, err := client.CurrentUser(ctx.Context, &emptypb.Empty{})
		require.Equal(t, codes.Unauthenticated, status.Code(err))

		r, _ = client.CurrentUser(server.NewContextWithToken(ctx.Context, loggedIn.Token), &emptypb.Empty{})
		current := r.User
		require.Equal(t, loggedIn, current)
	})

	t.Run("register: validation", func(t *testing.T) {
		tcs := map[string]struct {
			payload     *proto.UserRegistration
			wantCode    codes.Code
			wantMessage string
		}{
			"missing payload":  {nil, codes.Internal, "grpc: error while marshaling: proto: Marshal called with nil"},
			"missing username": {testUser.WithUsername("").Reg(), codes.InvalidArgument, "username: cannot be blank."},
			"missing email":    {testUser.WithEmail("").Reg(), codes.InvalidArgument, "email: cannot be blank."},
			"invalid email":    {testUser.WithEmail("invalid").Reg(), codes.InvalidArgument, "email: must be a valid email address."},
			"missing password": {testUser.WithPassword("").Reg(), codes.InvalidArgument, "password: cannot be blank."},
		}
		for name, tc := range tcs {
			t.Run(name, func(t *testing.T) {
				ctx := setup()
				defer ctx.teardown()

				client := proto.NewUsersClient(ctx.conn)

				_, err := client.RegisterUser(ctx, tc.payload)
				require.Equal(t, tc.wantCode.String(), status.Code(err).String())
				require.Equal(t, tc.wantMessage, status.Convert(err).Message())
			})
		}
	})
}

type testCtx struct {
	context.Context
	db       *sqlx.DB
	cfg      config.Config
	conn     *grpc.ClientConn
	teardown func()
}

func setup() *testCtx {
	ctx := context.Background()
	cfg := boot.ReadConfig("../../../../config.yml")
	db := boot.MustConnect(cfg.DataSource)
	conn, cleanup := apitest.Server(ctx, server.Router(cfg))
	return &testCtx{
		Context: ctx,
		db:      db,
		cfg:     cfg,
		conn:    conn,
		teardown: func() {
			cleanup()
			db.MustExec("DELETE FROM users")
		},
	}
}
