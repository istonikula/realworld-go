package domain

type Auth struct {
	Settings Security
}

func (a *Auth) NewToken(user UserId) string {
	// TODO Jwt
	return user.String()
}

func (a *Auth) EncryptPassword(plain string) string {
	// TODO encryptor
	return plain
}
