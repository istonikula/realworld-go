package main

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	domain "github.com/istonikula/realworld-go/realworld-domain"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	connStr := "user=realworld password=secret dbname=realworld sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// NOTE sql.Open does not create a connection to the database, it only validates the arguments provided
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	auth := domain.Auth{Settings: domain.Security{TokenSecret: "TODO token"}}
	userRepo := &UserRepo{db}

	router := gin.Default()

	router.POST("/api/users", func(c *gin.Context) {
		var dto UserRegistration
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validateUserSrv := domain.ValidateUserService{
			Auth:             auth,
			ExistsByUsername: userRepo.ExistsByUsername,
			ExistsByEmail:    userRepo.ExistsByEmail,
		}

		act, err := domain.RegisterUserUseCase{
			Validate:   validateUserSrv.ValidateUser,
			CreateUser: userRepo.Create,
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

	err = router.Run()
	if err != nil {
		log.Fatal(err)
	}
}

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

func (dto User) fromDomain(domain *domain.User) User {
	return User{
		Email:    domain.Email,
		Token:    domain.Token,
		Username: domain.Username,
		Bio:      domain.Bio,
		Image:    domain.Image,
	}
}

type UserRepo struct {
	db *sql.DB
}

func (r UserRepo) Create(user *domain.ValidUserRegistration) (*domain.User, error) {
	return nil, errors.New("TODO repo")
}

func (r UserRepo) ExistsByUsername(username string) bool {
	return false
}

func (r UserRepo) ExistsByEmail(email string) bool {
	return false
}
