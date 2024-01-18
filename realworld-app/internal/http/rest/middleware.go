package rest

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/istonikula/realworld-go/realworld-app/internal/db"
	domain "github.com/istonikula/realworld-go/realworld-domain"
	"github.com/jmoiron/sqlx"
)

func ResolveUser(auth *domain.Auth, txMgr *db.TxMgr) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := Context{c}

		token := auth.ValidateToken(ctx.Token())
		if token == nil {
			ctx.Next()
			return
		}

		err := txMgr.Read(func(tx *sqlx.Tx) error {
			repo := &db.UserRepo{Tx: tx}

			user, err := repo.FindById(token.Id)
			if err != nil {
				return err
			}

			if user != nil {
				ctx.SetUser(*user)
			}
			return nil
		})

		if err != nil {
			slog.Info(fmt.Errorf("ResolveUser: %w", err).Error())
		}

		ctx.Next()
	}
}

func RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := Context{c}

		if ctx.User() == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Next()
	}
}

type BindError struct {
	err error
}

func (b *BindError) Error() string {
	return b.err.Error()
}

func HandleLastError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}

		var errBind *BindError
		var errReg domain.UserRegistrationError
		var errV validation.Errors
		switch {
		case errors.As(err, &errBind):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.As(err, &errReg):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.As(err, &errV):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}
