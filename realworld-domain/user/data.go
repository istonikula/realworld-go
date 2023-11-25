package user

import "github.com/google/uuid"

type Id struct{ Value uuid.UUID }

func NewId() Id {
	return Id{Value: uuid.New()}
}

type User struct {
	Id       Id
	Email    string
	Token    string
	Username string
}

type Registration struct {
	Username string
	Email    string
	Password string
}
