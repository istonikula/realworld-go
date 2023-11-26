package user

type ValidateService struct {
	Auth             Auth
	ExistsByUsername ExistsByUsername
	ExistsByEmail    ExistsByEmail
}

func (s *ValidateService) Validate(r *Registration) (*ValidRegistration, error) {
	if s.ExistsByEmail(r.Email) {
		return nil, EmailAlreadyTaken
	}

	if s.ExistsByUsername(r.Username) {
		return nil, UsernameAlreadyTaken
	}

	id := NewId()
	valid := &ValidRegistration{
		Id:                id,
		Email:             r.Email,
		Username:          r.Username,
		Token:             s.Auth.NewToken(id),
		EncryptedPassword: s.Auth.EncryptPassword(r.Password),
	}
	return valid, nil
}
