package domain

type ValidateUserService struct {
	Auth             Auth
	ExistsByUsername ExistsByUsername
	ExistsByEmail    ExistsByEmail
}

func (s *ValidateUserService) ValidateUser(r *UserRegistration) (*ValidUserRegistration, error) {
	if s.ExistsByEmail(r.Email) {
		return nil, EmailAlreadyTaken
	}

	if s.ExistsByUsername(r.Username) {
		return nil, UsernameAlreadyTaken
	}

	id := NewUserId()
	valid := &ValidUserRegistration{
		Id:                id,
		Email:             r.Email,
		Username:          r.Username,
		Token:             s.Auth.NewToken(id),
		EncryptedPassword: s.Auth.EncryptPassword(r.Password),
	}
	return valid, nil
}
