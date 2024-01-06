package boot

import (
	"fmt"
	"log"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/istonikula/realworld-go/realworld-app/internal/config"
	"github.com/istonikula/realworld-go/realworld-app/internal/db"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/rest"
	domain "github.com/istonikula/realworld-go/realworld-domain"
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

func Migrate(c *config.DataSource) {
	u, _ := url.Parse(fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.MigrateUser, c.MigratePassword, c.Host, c.Port, c.Name,
	))

	err := dbmate.New(u).Migrate()
	if err != nil {
		log.Fatal(err)
	}
}

type repos struct{ User db.NewUserRepo }
type repoOpt func(r *repos)

func WithUserRepo(u db.NewUserRepo) repoOpt {
	return func(r *repos) {
		r.User = u
	}
}

func Router(cfg *config.Config, opts ...repoOpt) *gin.Engine {
	repos := &repos{
		User: db.UserRepoProvider,
	}
	for _, applyOpt := range opts {
		applyOpt(repos)
	}

	auth := &domain.Auth{Settings: domain.AuthSettings{TokenSecret: cfg.Auth.TokenSecret, TokenTTL: cfg.Auth.TokenTTL}}
	txMgr := &db.TxMgr{DB: MustConnect(&cfg.DataSource)}

	r := gin.Default()
	r.Use(rest.HandleLastError(), rest.ResolveUser(auth, txMgr))
	rest.UserRoutes(r, auth, txMgr, repos.User)

	return r
}

func MustConnect(c *config.DataSource) *sqlx.DB {
	return sqlx.MustConnect(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Password, c.Name),
	)
}
