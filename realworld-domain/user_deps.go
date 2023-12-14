package domain

type ValidateUserRegistration func(*UserRegistration) (*ValidUserRegistration, error)
type CreateUser func(*ValidUserRegistration) (*User, error)

type ExistsByUsername func(username string) (bool, error)
type ExistsByEmail func(email string) (bool, error)
