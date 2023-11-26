package user

type ValidateRegistration func(*Registration) (*ValidRegistration, error)
type CreateUser func(*ValidRegistration) (*User, error)

type ExistsByUsername func(username string) bool
type ExistsByEmail func(email string) bool
