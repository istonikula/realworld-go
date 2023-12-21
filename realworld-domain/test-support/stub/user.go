package stub

import (
	"errors"

	domain "github.com/istonikula/realworld-go/realworld-domain"
)

var UserStub = struct {
	Auth                 domain.Auth
	CreateUser           domain.CreateUser
	CreateUserError      domain.CreateUser
	UnexpectedCreateUser domain.CreateUser
	ValidateUser         func(func(domain.UserRegistration) *domain.ValidUserRegistration) domain.ValidateUserRegistration
	ValidateUserError    func(*domain.UserRegistrationError) domain.ValidateUserRegistration
}{
	Auth: domain.Auth{Settings: domain.AuthSettings{TokenSecret: "secret", TokenTTL: 1800}},

	CreateUser: func(r *domain.ValidUserRegistration) (*domain.User, error) {
		return &domain.User{Id: domain.NewUserId(), Email: r.Email, Token: r.Token, Username: r.Username}, nil
	},
	CreateUserError: func(*domain.ValidUserRegistration) (*domain.User, error) {
		return nil, errors.New("unexpected error")
	},
	UnexpectedCreateUser: func(*domain.ValidUserRegistration) (*domain.User, error) {
		panic("unexpected create user")
	},

	ValidateUser: func(f func(domain.UserRegistration) *domain.ValidUserRegistration) domain.ValidateUserRegistration {
		return func(r *domain.UserRegistration) (*domain.ValidUserRegistration, error) {
			return f(*r), nil
		}
	},
	ValidateUserError: func(e *domain.UserRegistrationError) domain.ValidateUserRegistration {
		return func(*domain.UserRegistration) (*domain.ValidUserRegistration, error) {
			return nil, e
		}
	},
}
