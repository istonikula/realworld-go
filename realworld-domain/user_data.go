package domain

import "github.com/google/uuid"

type UserId struct{ uuid.UUID }

func NewUserId() UserId {
	return UserId{uuid.New()}
}

type User struct {
	Id       UserId
	Email    string
	Token    string
	Username string
	Bio      *string
	Image    *string
}

type UserRegistration struct {
	Username string
	Email    string
	Password string
}

type ValidUserRegistration struct {
	Id                UserId
	Email             string
	Token             string
	Username          string
	EncryptedPassword string
}

type Login struct {
	Email    string
	Password string
}

type UserAndPassword struct {
	User
	EncryptedPassword string
}
