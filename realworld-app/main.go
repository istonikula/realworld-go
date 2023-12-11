package main

import (
	"fmt"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/gin-gonic/gin"
	domain "github.com/istonikula/realworld-go/realworld-domain"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func migrate() {
	u, _ := url.Parse("postgres://postgres:secret@127.0.0.1:5432/realworld?sslmode=disable")
	err := dbmate.New(u).Migrate()
	if err != nil {
		log.Fatal(err)
	}
}

func db() *sqlx.DB {
	return sqlx.MustConnect("postgres", "user=realworld password=secret dbname=realworld sslmode=disable")
}

func router(db *sqlx.DB) *gin.Engine {
	auth := domain.Auth{Settings: domain.Security{TokenSecret: "TODO token"}}
	userRepo := &UserRepo{db}

	r := gin.Default()

	r.POST("/api/users", func(c *gin.Context) {
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

	return r
}

func main() {
	migrate()

	if err := router(db()).Run(); err != nil {
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
	db *sqlx.DB
}

func (r UserRepo) Create(reg *domain.ValidUserRegistration) (*domain.User, error) {
	q := Table("users").Insert(
		"id", "email", "token", "username", "password",
	) + " RETURNING id, email, token, username, bio, image"
	stmt, err := r.db.Preparex(q)
	if err != nil {
		return nil, fmt.Errorf("UserRepo#Create: prepare failed: %w", err)
	}

	var user domain.User
	err = stmt.QueryRowx(
		reg.Id, reg.Email, reg.Token, reg.Username, reg.EncryptedPassword,
	).StructScan(&user)
	if err != nil {
		return nil, fmt.Errorf("UserRepo#Create: insert failed: %w", err)
	}
	return &user, nil
}

func (r UserRepo) ExistsByUsername(username string) bool {
	return false
}

func (r UserRepo) ExistsByEmail(email string) bool {
	return false
}

type Table string

func (t Table) Insert(cols ...string) string {
	markerSql := "$1"
	if len(cols) > 1 {
		for i := range cols[1:] {
			markerSql += ", $" + strconv.Itoa(i+2)
		}
	}
	return "INSERT INTO " + string(t) + " (" + strings.Join(cols, ", ") + ") VALUES (" + markerSql + ")"
}
