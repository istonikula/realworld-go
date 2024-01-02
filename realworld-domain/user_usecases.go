package domain

type UserRegistrationError struct {
	Kind string
}

func (e *UserRegistrationError) Error() string { return e.Kind }

var (
	EmailAlreadyTaken    = &UserRegistrationError{"email already taken"}
	UsernameAlreadyTaken = &UserRegistrationError{"username already taken"}
)

type UserLoginError struct {
	Kind string
}

func (e *UserLoginError) Error() string { return e.Kind }

var (
	UserNotFound   = &UserLoginError{"user not found"}
	BadCredentials = &UserLoginError{"bad credentials"}
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

type LoginUserUseCase struct {
	Auth    Auth
	GetUser GetUserByEmail
}

func (u LoginUserUseCase) Run(l *Login) (*User, error) {
	found, err := u.GetUser(l.Email)
	if err != nil {
		return nil, err
	}

	if found == nil {
		return nil, UserNotFound
	}

	if ok := u.Auth.CheckPassword(l.Password, found.EncryptedPassword); !ok {
		return nil, BadCredentials
	}

	return &found.User, nil
}
