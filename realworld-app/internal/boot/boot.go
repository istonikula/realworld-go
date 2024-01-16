package boot

import (
	"fmt"
	"log"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/istonikula/realworld-go/realworld-app/internal/config"
	"github.com/istonikula/realworld-go/realworld-app/internal/db"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ReadConfig(path ...string) *config.Config {
	p := "config.yml"
	if len(path) > 0 {
		p = path[0]
	}

	var cfg config.Config
	if err := cleanenv.ReadConfig(p, &cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}

func Migrate(path string, c *config.DataSource) {
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

type Repos struct{ User db.NewUserRepo }
type RepoOpt func(r *Repos)

func WithUserRepo(u db.NewUserRepo) RepoOpt {
	return func(r *Repos) {
		r.User = u
	}
}

func MustConnect(c *config.DataSource) *sqlx.DB {
	return sqlx.MustConnect(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Password, c.Name),
	)
}
