package domain

import "github.com/google/uuid"

type Auth struct {
	Settings Security
}

func (a *Auth) NewToken(user UserId) string {
	// TODO Jwt
	return uuid.UUID(user).String()
}

func (a *Auth) EncryptPassword(plain string) string {
	// TODO encryptor
	return plain
}
