package user

type RegisterUseCase struct {
	CreateUser CreateUser
}

func (u *RegisterUseCase) Run(r *Registration) User {
	return u.CreateUser(r)
}
