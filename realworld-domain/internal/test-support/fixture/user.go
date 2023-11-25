package fixture

import (
	"fmt"

	domain "github.com/istonikula/realworld-go/realworld-domain"
	"github.com/istonikula/realworld-go/realworld-domain/user"
)

type UserFactory struct {
	Auth domain.Auth
}

func (f UserFactory) NewUser(username string, id ...user.Id) user.User {
	userId := user.NewId()
	if len(id) == 1 {
		userId = id[0]
	}
	return user.User{
		Id:       userId,
		Email:    fmt.Sprintf("%s@realwold.io", username),
		Token:    "",
		Username: username,
	}
}

type TestUser user.User

func (u TestUser) Registration() user.Registration {
	return user.Registration{Username: u.Username, Email: u.Email, Password: "plain"}
}
