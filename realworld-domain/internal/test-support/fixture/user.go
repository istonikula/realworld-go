package fixture

import (
	"fmt"

	"github.com/istonikula/realworld-go/realworld-domain/user"
)

type UserFactory struct {
	Auth user.Auth
}

func (f UserFactory) NewUser(username string, id ...user.Id) *user.User {
	userId := user.NewId()
	if len(id) == 1 {
		userId = id[0]
	}
	return &user.User{
		Id:       userId,
		Email:    fmt.Sprintf("%s@realwold.io", username),
		Token:    "",
		Username: username,
	}
}

func (f UserFactory) ValidRegistration(r *user.Registration) *user.ValidRegistration {
	id := user.NewId()

	return &user.ValidRegistration{
		Id:                id,
		Email:             r.Email,
		Username:          r.Username,
		Token:             f.Auth.NewToken(id),
		EncryptedPassword: f.Auth.EncryptPassword(r.Password),
	}
}

type TestUser user.User

func (u *TestUser) Registration() *user.Registration {
	return &user.Registration{Username: u.Username, Email: u.Email, Password: "plain"}
}
