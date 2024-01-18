package proto

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (r *UserRegistration) Validate() error {
	return v.ValidateStruct(r,
		v.Field(&r.Username, v.Required),
		v.Field(&r.Email, v.Required, is.Email),
		v.Field(&r.Password, v.Required),
	)
}

func (l *LoginRequest) Validate() error {
	return v.ValidateStruct(l,
		v.Field(&l.Email, v.Required, is.Email),
		v.Field(&l.Password, v.Required),
	)
}
