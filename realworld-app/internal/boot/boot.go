package boot

import (
	"fmt"
	"log"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/istonikula/realworld-go/realworld-app/internal/config"
	"github.com/istonikula/realworld-go/realworld-app/internal/db"
	"github.com/jmoiron/sqlx"

	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	_ "github.com/lib/pq"
)

func ReadConfig(path string) config.Config {
	var cfg config.Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}

func Migrate(path string, c config.DataSource) {
	u, _ := url.Parse(fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.MigrateUser, c.MigratePassword, c.Host, c.Port, c.Name,
	))

	db := dbmate.New(u)
	db.MigrationsDir = []string{path + "/migrations"}
	db.SchemaFile = path + "/schema.sql"

	err := db.Migrate()
	if err != nil {
		log.Fatal(err)
	}
}

func MustConnect(c config.DataSource) *sqlx.DB {
	return sqlx.MustConnect(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Password, c.Name),
	)
}

type RouterOptions struct {
	UserRepo db.NewUserRepo
}

type RouterOption func(*RouterOptions)

func UserRepo(u db.NewUserRepo) RouterOption {
	return func(o *RouterOptions) {
		o.UserRepo = u
	}
}
