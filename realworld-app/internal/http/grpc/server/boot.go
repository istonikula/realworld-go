package server

import (
	"github.com/istonikula/realworld-go/realworld-app/internal/boot"
	"github.com/istonikula/realworld-go/realworld-app/internal/config"
	"github.com/istonikula/realworld-go/realworld-app/internal/db"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/grpc/proto"
	domain "github.com/istonikula/realworld-go/realworld-domain"
	"google.golang.org/grpc"
)

func Router(cfg *config.Config, opts ...boot.RepoOpt) *grpc.Server {
	repos := &boot.Repos{
		User: db.UserRepoProvider,
	}
	for _, applyOpt := range opts {
		applyOpt(repos)
	}

	auth := &domain.Auth{Settings: domain.AuthSettings{TokenSecret: cfg.Auth.TokenSecret, TokenTTL: cfg.Auth.TokenTTL}}
	txMgr := &db.TxMgr{DB: boot.MustConnect(&cfg.DataSource)}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(ResolveUser(auth, txMgr), RequireUser()))
	proto.RegisterUsersServer(s, UserRoutes(auth, txMgr, repos.User))

	return s
}
