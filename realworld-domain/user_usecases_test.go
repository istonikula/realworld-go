package domain_test

import (
	rwt "github.com/istonikula/realworld-go/realworld-testing"
	"testing"

	"github.com/istonikula/realworld-go/realworld-domain"
	"github.com/istonikula/realworld-go/realworld-domain/internal/test-support/fixture"
	"github.com/istonikula/realworld-go/realworld-domain/internal/test-support/stub"
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
		rwt.Ok(t, err)
		rwt.Equals(t, jane.Email, act.Email)
		rwt.Equals(t, jane.Username, act.Username)
	})

	t.Run("email already taken", func(t *testing.T) {
		_, err := domain.RegisterUserUseCase{
			Validate:   stub.UserStub.ValidateUserError(domain.EmailAlreadyTaken),
			CreateUser: stub.UserStub.UnexpectedCreateUser,
		}.Run(reg)
		rwt.Assert(t, err != nil, "error expected")
		rwt.Equals(t, domain.EmailAlreadyTaken, err)
	})

	t.Run("username already taken", func(t *testing.T) {
		_, err := domain.RegisterUserUseCase{
			Validate:   stub.UserStub.ValidateUserError(domain.UsernameAlreadyTaken),
			CreateUser: stub.UserStub.UnexpectedCreateUser,
		}.Run(reg)
		rwt.Assert(t, err != nil, "error expected")
		rwt.Equals(t, domain.UsernameAlreadyTaken, err)
	})

	t.Run("create failure", func(t *testing.T) {
		_, err := domain.RegisterUserUseCase{
			Validate:   stub.UserStub.ValidateUser(userFactory.ValidRegistration),
			CreateUser: stub.UserStub.CreateUserError,
		}.Run(reg)
		rwt.Assert(t, err != nil, "error expected")
		rwt.Equals(t, "unexpected error", err.Error())
	})
}