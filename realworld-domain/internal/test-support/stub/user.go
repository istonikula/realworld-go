package stub

import (
	domain "github.com/istonikula/realworld-go/realworld-domain"
	"github.com/istonikula/realworld-go/realworld-domain/user"
)

var UserStub = struct {
	Auth       domain.Auth
	CreateUser user.CreateUser
}{
	Auth: domain.Auth{Settings: domain.Security{TokenSecret: ""}},

	// TODO token
	CreateUser: func(r *user.Registration) user.User {
		return user.User{Id: user.NewId(), Email: r.Email, Token: "", Username: r.Username}
	},
}
