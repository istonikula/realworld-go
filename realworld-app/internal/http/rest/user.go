package rest

import (
	"errors"
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
	router.GET("/api/user", ResolveUser(auth, txMgr), RequireUser(), func(c *gin.Context) {
		ctx := Context{c}
		ctx.JSON(http.StatusOK, UserResponse{User{}.fromDomain(ctx.User())})
	})

	router.POST("/api/users", func(c *gin.Context) {
		var dto UserRegistration
		err := c.ShouldBindJSON(&dto)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			var regErr *domain.UserRegistrationError
			if errors.As(err, &regErr) {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusCreated, UserResponse{User{}.fromDomain(u)})
	})

	router.POST("/api/users/login", func(c *gin.Context) {
		var dto Login
		err := c.ShouldBindJSON(&dto)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			c.Status(http.StatusUnauthorized)
			return
		}

		c.JSON(http.StatusOK, UserResponse{User{}.fromDomain(u)})
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
