package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	domain "github.com/istonikula/realworld-go/realworld-domain"
	"github.com/istonikula/realworld-go/realworld-domain/test-support/fixture"
	"github.com/istonikula/realworld-go/realworld-domain/test-support/stub"
)

var userFactory = fixture.UserFactory{Auth: stub.UserStub.Auth}
var jane = fixture.TestUser(*userFactory.NewUser("jane"))

func TestRegisterUserUseCase(t *testing.T) {
	reg := jane.Registration()

	t.Run("register", func(t *testing.T) {
		act, err := domain.RegisterUserUseCase{
			Validate:   stub.UserStub.ValidateUser(userFactory.ValidRegistration),
			CreateUser: stub.UserStub.CreateUser,
		}.Run(reg)
		require.NoError(t, err)
		require.Equal(t, jane.Email, act.Email)
		require.Equal(t, jane.Username, act.Username)
	})

	t.Run("email already taken", func(t *testing.T) {
		_, err := domain.RegisterUserUseCase{
			Validate:   stub.UserStub.ValidateUserError(domain.EmailAlreadyTaken),
			CreateUser: stub.UserStub.UnexpectedCreateUser,
		}.Run(reg)
		require.Equal(t, domain.EmailAlreadyTaken, err)
	})

	t.Run("username already taken", func(t *testing.T) {
		_, err := domain.RegisterUserUseCase{
			Validate:   stub.UserStub.ValidateUserError(domain.UsernameAlreadyTaken),
			CreateUser: stub.UserStub.UnexpectedCreateUser,
		}.Run(reg)
		require.Equal(t, domain.UsernameAlreadyTaken, err)
	})

	t.Run("create failure", func(t *testing.T) {
		_, err := domain.RegisterUserUseCase{
			Validate:   stub.UserStub.ValidateUser(userFactory.ValidRegistration),
			CreateUser: stub.UserStub.CreateUserError,
		}.Run(reg)
		require.Equal(t, "unexpected error", err.Error())
	})
}

func TestLoginUserUseCase(t *testing.T) {
	reg := jane.Registration()

	t.Run("login", func(t *testing.T) {
		exp := userFactory.UserAndPassword(reg)

		act, err := domain.LoginUserUseCase{
			Auth: stub.UserStub.Auth,
			GetUser: func(email string) (*domain.UserAndPassword, error) {
				return &exp, nil
			},
		}.Run(domain.Login{
			Email:    reg.Email,
			Password: reg.Password,
		})
		require.NoError(t, err)
		require.Equal(t, &exp.User, act)
	})

	t.Run("not found", func(t *testing.T) {
		act, err := domain.LoginUserUseCase{
			Auth: stub.UserStub.Auth,
			GetUser: func(email string) (*domain.UserAndPassword, error) {
				return nil, nil
			},
		}.Run(domain.Login{
			Email:    reg.Email,
			Password: reg.Password,
		})
		require.Equal(t, domain.UserNotFound, err)
		require.Nil(t, act)
	})

	t.Run("bad credentials", func(t *testing.T) {
		act, err := domain.LoginUserUseCase{
			Auth: stub.UserStub.Auth,
			GetUser: func(email string) (*domain.UserAndPassword, error) {
				u := userFactory.UserAndPassword(reg)
				return &u, nil
			},
		}.Run(domain.Login{
			Email:    reg.Email,
			Password: "invalid",
		})
		require.Equal(t, domain.BadCredentials, err)
		require.Nil(t, act)
	})
}
