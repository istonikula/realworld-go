package fixture

import (
	"fmt"

	domain "github.com/istonikula/realworld-go/realworld-domain"
)

type UserFactory struct {
	Auth domain.Auth
}

func (f UserFactory) NewUser(username string, id ...domain.UserId) *domain.User {
	userId := domain.NewUserId()
	if len(id) == 1 {
		userId = id[0]
	}
	return &domain.User{
		Id:       userId,
		Email:    fmt.Sprintf("%s@realwold.io", username),
		Token:    "",
		Username: username,
	}
}

func (f UserFactory) ValidRegistration(r domain.UserRegistration) *domain.ValidUserRegistration {
	id := domain.NewUserId()

	token, _ := f.Auth.NewToken(id)

	return &domain.ValidUserRegistration{
		Id:                id,
		Email:             r.Email,
		Username:          r.Username,
		Token:             token,
		EncryptedPassword: f.Auth.EncryptPassword(r.Password),
	}
}

type TestUser domain.User

func (u *TestUser) Registration() *domain.UserRegistration {
	return &domain.UserRegistration{Username: u.Username, Email: u.Email, Password: "plain"}
}
