package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/istonikula/realworld-go/realworld-app/internal/boot"
	"github.com/istonikula/realworld-go/realworld-app/internal/config"
	"github.com/istonikula/realworld-go/realworld-app/internal/db"
	domain "github.com/istonikula/realworld-go/realworld-domain"
)

func Router(cfg config.Config, opt ...boot.RouterOption) *gin.Engine {
	opts := &boot.RouterOptions{
		UserRepo: db.UserRepoProvider,
	}
	for _, o := range opt {
		o(opts)
	}

	auth := &domain.Auth{Settings: domain.AuthSettings{TokenSecret: cfg.Auth.TokenSecret, TokenTTL: cfg.Auth.TokenTTL}}
	txMgr := &db.TxMgr{DB: boot.MustConnect(cfg.DataSource)}

	r := gin.Default()
	r.Use(HandleLastError(), ResolveUser(auth, txMgr))
	UserRoutes(r, auth, txMgr, opts.UserRepo)

	return r
}
