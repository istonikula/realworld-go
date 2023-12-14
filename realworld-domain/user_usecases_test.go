package domain_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/istonikula/realworld-go/realworld-domain"
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
		assert.NoError(t, err)
		assert.Equal(t, jane.Email, act.Email)
		assert.Equal(t, jane.Username, act.Username)
	})

	t.Run("email already taken", func(t *testing.T) {
		_, err := domain.RegisterUserUseCase{
			Validate:   stub.UserStub.ValidateUserError(domain.EmailAlreadyTaken),
			CreateUser: stub.UserStub.UnexpectedCreateUser,
		}.Run(reg)
		assert.Equal(t, domain.EmailAlreadyTaken, err)
	})

	t.Run("username already taken", func(t *testing.T) {
		_, err := domain.RegisterUserUseCase{
			Validate:   stub.UserStub.ValidateUserError(domain.UsernameAlreadyTaken),
			CreateUser: stub.UserStub.UnexpectedCreateUser,
		}.Run(reg)
		assert.Equal(t, domain.UsernameAlreadyTaken, err)
	})

	t.Run("create failure", func(t *testing.T) {
		_, err := domain.RegisterUserUseCase{
			Validate:   stub.UserStub.ValidateUser(userFactory.ValidRegistration),
			CreateUser: stub.UserStub.CreateUserError,
		}.Run(reg)
		assert.Equal(t, "unexpected error", err.Error())
	})
}
