package rest

import (
	"errors"
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
		var err error
		if err = c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var act *domain.User
		err = txMgr.Write(func(tx *sqlx.Tx) error {
			repo := &appDb.UserRepo{Tx: tx}

			validateUserSrv := domain.ValidateUserService{
				Auth:             *auth,
				ExistsByUsername: repo.ExistsByUsername,
				ExistsByEmail:    repo.ExistsByEmail,
			}

			act, err = domain.RegisterUserUseCase{
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
			if errors.Is(err, domain.EmailAlreadyTaken) || errors.Is(err, domain.UsernameAlreadyTaken) {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusCreated, UserResponse{User{}.fromDomain(act)})
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
