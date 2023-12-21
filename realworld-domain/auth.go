package domain

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrUnauthorized = errors.New("unauthorized")

type AuthSettings struct {
	TokenSecret string
	TokenTTL    int64
}

type Auth struct{ Settings AuthSettings }
type Token struct{ Id UserId }

func (a *Auth) NewToken(user UserId) (string, error) {
	iat := time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": user.String(),
		"iat": iat,
		"eat": iat + a.Settings.TokenTTL,
	})

	return token.SignedString([]byte(a.Settings.TokenSecret))
}

func (a *Auth) ValidateToken(str string) *Token {
	if str == "" {
		return nil
	}

	token, err := a.parse(str)
	if err != nil {
		slog.Info(fmt.Errorf("Auth#ValidateToken parse: %v", err).Error())
		return nil
	}
	if !token.Valid {
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	idStr, err := claims.GetSubject()
	if err != nil {
		slog.Info(fmt.Errorf("Auth#ValidateToken get subject: %w", err).Error())
		return nil
	}

	idUuid, err := uuid.Parse(idStr)
	if err != nil {
		return nil
	}

	return &Token{Id: UserId{idUuid}}
}

func (a *Auth) EncryptPassword(plain string) string {
	// TODO encryptor
	return plain
}

func (a *Auth) parse(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(parsed *jwt.Token) (any, error) {
		if _, ok := parsed.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Info(fmt.Sprintf("unexpected signing method %v", parsed.Header["alg"]))
			return nil, ErrUnauthorized
		}

		return []byte(a.Settings.TokenSecret), nil
	})
}
