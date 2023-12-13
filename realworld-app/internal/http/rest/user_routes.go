package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/istonikula/realworld-go/realworld-app/internal/db"
	domain "github.com/istonikula/realworld-go/realworld-domain"
	"net/http"
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

func UserRoutes(router *gin.Engine, auth *domain.Auth, repo *db.UserRepo) {
	router.POST("/api/users", func(c *gin.Context) {
		var dto UserRegistration
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validateUserSrv := domain.ValidateUserService{
			Auth:             *auth,
			ExistsByUsername: repo.ExistsByUsername,
			ExistsByEmail:    repo.ExistsByEmail,
		}

		act, err := domain.RegisterUserUseCase{
			Validate:   validateUserSrv.ValidateUser,
			CreateUser: repo.Create,
		}.Run(&domain.UserRegistration{
			Username: dto.Username,
			Email:    dto.Email,
			Password: dto.Password,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
