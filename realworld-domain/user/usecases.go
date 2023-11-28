package user

type RegistrationError struct {
	Kind string
}

func (e *RegistrationError) Error() string { return e.Kind }

var (
	EmailAlreadyTaken    = &RegistrationError{"email already taken"}
	UsernameAlreadyTaken = &RegistrationError{"username already taken"}
)

type RegisterUseCase struct {
	Validate   ValidateRegistration
	CreateUser CreateUser
}

func (u RegisterUseCase) Run(r *Registration) (*User, error) {
	valid, err := u.Validate(r)
	if err != nil {
		return nil, err
	}

	return u.CreateUser(valid)
}
