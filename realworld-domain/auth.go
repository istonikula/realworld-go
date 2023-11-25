package domain

import "github.com/istonikula/realworld-go/realworld-domain/user"

type Auth struct {
	Settings Security
}

func (a *Auth) NewToken(user user.Id) string {
	// TODO Jwt
	return user.Value.String()
}
