package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/istonikula/realworld-go/realworld-app/internal/config"
	appDb "github.com/istonikula/realworld-go/realworld-app/internal/db"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/rest"
	domain "github.com/istonikula/realworld-go/realworld-domain"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg := readConfig()

	migrate(&cfg.DataSource)

	if err := router(db(&cfg.DataSource), cfg).Run(); err != nil {
		log.Fatal(err)
	}
}

func readConfig() *config.Config {
	var cfg config.Config
	if err := cleanenv.ReadConfig("config.yml", &cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}

func migrate(c *config.DataSource) {
	u, _ := url.Parse(fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.MigrateUser, c.MigratePassword, c.Host, c.Port, c.Name,
	))

	err := dbmate.New(u).Migrate()
	if err != nil {
		log.Fatal(err)
	}
}

func db(c *config.DataSource) *sqlx.DB {
	return sqlx.MustConnect(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Password, c.Name),
	)
}

func router(db *sqlx.DB, cfg *config.Config) *gin.Engine {
	auth := domain.Auth{Settings: domain.AuthSettings{TokenSecret: cfg.Auth.TokenSecret, TokenTTL: cfg.Auth.TokenTTL}}

	txMgr := &appDb.TxMgr{DB: db}

	r := gin.Default()
	rest.UserRoutes(r, &auth, txMgr)

	return r
}
