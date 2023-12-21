package rest

import (
	"strings"

	"github.com/gin-gonic/gin"
	domain "github.com/istonikula/realworld-go/realworld-domain"
)

type Context struct{ *gin.Context }

const userKey = "user"

func (c Context) User() *domain.User {
	maybeUser, exists := c.Get(userKey)
	if !exists {
		return nil
	}

	user, ok := maybeUser.(domain.User)
	if !ok {
		return nil
	}

	return &user
}

func (c Context) SetUser(user domain.User) {
	c.Set(userKey, user)
}

func (c Context) Token() string {
	const prefix = "Token "
	auth := c.Request.Header.Get("Authorization")
	return strings.TrimPrefix(auth, prefix)
}
