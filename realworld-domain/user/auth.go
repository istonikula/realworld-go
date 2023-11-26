package user

import domain "github.com/istonikula/realworld-go/realworld-domain"

type Auth struct {
	Settings domain.Security
}

func (a *Auth) NewToken(user Id) string {
	// TODO Jwt
	return user.Value.String()
}

func (a *Auth) EncryptPassword(plain string) string {
	// TODO encryptor
	return plain
}
