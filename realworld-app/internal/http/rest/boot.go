package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/istonikula/realworld-go/realworld-app/internal/boot"
	"github.com/istonikula/realworld-go/realworld-app/internal/config"
	"github.com/istonikula/realworld-go/realworld-app/internal/db"
	domain "github.com/istonikula/realworld-go/realworld-domain"
)

func Router(cfg *config.Config, opts ...boot.RepoOpt) *gin.Engine {
	repos := &boot.Repos{
		User: db.UserRepoProvider,
	}
	for _, applyOpt := range opts {
		applyOpt(repos)
	}

	auth := &domain.Auth{Settings: domain.AuthSettings{TokenSecret: cfg.Auth.TokenSecret, TokenTTL: cfg.Auth.TokenTTL}}
	txMgr := &db.TxMgr{DB: boot.MustConnect(&cfg.DataSource)}

	r := gin.Default()
	r.Use(HandleLastError(), ResolveUser(auth, txMgr))
	UserRoutes(r, auth, txMgr, repos.User)

	return r
}
