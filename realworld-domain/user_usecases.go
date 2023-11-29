package domain

type UserRegistrationError struct {
	Kind string
}

func (e *UserRegistrationError) Error() string { return e.Kind }

var (
	EmailAlreadyTaken    = &UserRegistrationError{"email already taken"}
	UsernameAlreadyTaken = &UserRegistrationError{"username already taken"}
)

type RegisterUserUseCase struct {
	Validate   ValidateUserRegistration
	CreateUser CreateUser
}

func (u RegisterUserUseCase) Run(r *UserRegistration) (*User, error) {
	valid, err := u.Validate(r)
	if err != nil {
		return nil, err
	}

	return u.CreateUser(valid)
}
