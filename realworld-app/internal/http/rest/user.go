package rest

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	appDb "github.com/istonikula/realworld-go/realworld-app/internal/db"
	domain "github.com/istonikula/realworld-go/realworld-domain"
	"github.com/jmoiron/sqlx"
)

type UserRegistration struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	User User `json:"user"`
}

type User struct {
	Email    string  `json:"email"`
	Token    string  `json:"token"`
	Username string  `json:"username"`
	Bio      *string `json:"bio"`
	Image    *string `json:"image"`
}

func UserRoutes(router *gin.Engine, auth *domain.Auth, txMgr *appDb.TxMgr) {
	router.GET("/api/user", RequireUser(), func(c *gin.Context) {
		ctx := Context{c}
		ctx.JSON(http.StatusOK, UserResponse{User{}.fromDomain(ctx.User())})
	})

	router.POST("/api/users", func(c *gin.Context) {
		ctx := Context{c}

		var dto UserRegistration
		err := c.ShouldBindJSON(&dto)
		if err != nil {
			ctx.AbortWithError(&BindError{err})
			return
		}

		var u *domain.User
		err = txMgr.Write(func(tx *sqlx.Tx) error {
			repo := &appDb.UserRepo{Tx: tx}

			validateUserSrv := domain.ValidateUserService{
				Auth:             *auth,
				ExistsByUsername: repo.ExistsByUsername,
				ExistsByEmail:    repo.ExistsByEmail,
			}

			u, err = domain.RegisterUserUseCase{
				Validate:   validateUserSrv.ValidateUser,
				CreateUser: repo.Create,
			}.Run(&domain.UserRegistration{
				Username: dto.Username,
				Email:    dto.Email,
				Password: dto.Password,
			})
			return err
		})

		if err != nil {
			ctx.AbortWithError(err)
			return
		}

		c.JSON(http.StatusCreated, UserResponse{User{}.fromDomain(u)})
	})

	router.POST("/api/users/login", func(c *gin.Context) {
		ctx := Context{c}

		var dto Login
		err := ctx.ShouldBindJSON(&dto)
		if err != nil {
			ctx.AbortWithError(&BindError{err})
			return
		}

		var u *domain.User
		err = txMgr.Read(func(tx *sqlx.Tx) error {
			repo := &appDb.UserRepo{Tx: tx}

			u, err = domain.LoginUserUseCase{
				Auth:    *auth,
				GetUser: repo.FindByEmail,
			}.Run(&domain.Login{
				Email:    dto.Email,
				Password: dto.Password,
			})
			return err
		})

		if err != nil {
			slog.Info(fmt.Errorf("login: %w", err).Error())
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.JSON(http.StatusOK, UserResponse{User{}.fromDomain(u)})
	})
}

func (dto User) fromDomain(domain *domain.User) User {
	return User{
		Email:    domain.Email,
		Token:    domain.Token,
		Username: domain.Username,
		Bio:      domain.Bio,
		Image:    domain.Image,
	}
}
