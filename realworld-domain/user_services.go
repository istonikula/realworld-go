package domain

type ValidateUserService struct {
	Auth             Auth
	ExistsByUsername ExistsByUsername
	ExistsByEmail    ExistsByEmail
}

func (s *ValidateUserService) ValidateUser(r *UserRegistration) (*ValidUserRegistration, error) {
	if exists, err := s.ExistsByEmail(r.Email); err != nil {
		return nil, err
	} else if exists {
		return nil, EmailAlreadyTaken
	}

	if exists, err := s.ExistsByUsername(r.Username); err != nil {
		return nil, err
	} else if exists {
		return nil, UsernameAlreadyTaken
	}

	id := NewUserId()

	token, err := s.Auth.NewToken(id)
	if err != nil {
		return nil, err
	}

	valid := &ValidUserRegistration{
		Id:                id,
		Email:             r.Email,
		Username:          r.Username,
		Token:             token,
		EncryptedPassword: s.Auth.EncryptPassword(r.Password),
	}
	return valid, nil
}
