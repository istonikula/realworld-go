package domain

type ValidateUserRegistration func(*UserRegistration) (*ValidUserRegistration, error)
type CreateUser func(*ValidUserRegistration) (*User, error)

type ExistsByUsername func(username string) bool
type ExistsByEmail func(email string) bool
