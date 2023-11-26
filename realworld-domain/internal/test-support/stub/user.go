package stub

import (
	"errors"
	domain "github.com/istonikula/realworld-go/realworld-domain"
	"github.com/istonikula/realworld-go/realworld-domain/user"
)

var UserStub = struct {
	Auth                 user.Auth
	CreateUser           user.CreateUser
	CreateUserError      user.CreateUser
	UnexpectedCreateUser user.CreateUser
	ValidateUser         func(func(*user.Registration) *user.ValidRegistration) user.ValidateRegistration
	ValidateUserError    func(*user.RegistrationError) user.ValidateRegistration
}{
	Auth: user.Auth{Settings: domain.Security{TokenSecret: ""}},

	CreateUser: func(r *user.ValidRegistration) (*user.User, error) {
		return &user.User{Id: user.NewId(), Email: r.Email, Token: r.Token, Username: r.Username}, nil
	},
	CreateUserError: func(*user.ValidRegistration) (*user.User, error) {
		return nil, errors.New("unexpected error")
	},
	UnexpectedCreateUser: func(*user.ValidRegistration) (*user.User, error) {
		panic("unexpected create user")
	},

	ValidateUser: func(f func(*user.Registration) *user.ValidRegistration) user.ValidateRegistration {
		return func(r *user.Registration) (*user.ValidRegistration, error) {
			return f(r), nil
		}
	},
	ValidateUserError: func(e *user.RegistrationError) user.ValidateRegistration {
		return func(*user.Registration) (*user.ValidRegistration, error) {
			return nil, e
		}
	},
}
