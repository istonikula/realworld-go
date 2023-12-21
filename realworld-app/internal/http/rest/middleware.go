package rest

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/istonikula/realworld-go/realworld-app/internal/db"
	domain "github.com/istonikula/realworld-go/realworld-domain"
	"github.com/jmoiron/sqlx"
)

func ResolveUser(auth *domain.Auth, txMgr *db.TxMgr) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := Context{c}

		err := txMgr.Read(func(tx *sqlx.Tx) error {
			repo := &db.UserRepo{Tx: tx}

			token := auth.ValidateToken(ctx.Token())
			if token == nil {
				return nil
			}

			user, err := repo.FindById(token.Id)
			if err != nil {
				return err
			}

			ctx.SetUser(*user)
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
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
