package server

import (
	"github.com/istonikula/realworld-go/realworld-app/internal/boot"
	"github.com/istonikula/realworld-go/realworld-app/internal/config"
	"github.com/istonikula/realworld-go/realworld-app/internal/db"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/grpc/proto"
	domain "github.com/istonikula/realworld-go/realworld-domain"
	"google.golang.org/grpc"
)

func Router(cfg config.Config, opt ...boot.RouterOption) *grpc.Server {
	opts := &boot.RouterOptions{
		UserRepo: db.UserRepoProvider,
	}
	for _, o := range opt {
		o(opts)
	}

	auth := &domain.Auth{Settings: domain.AuthSettings{TokenSecret: cfg.Auth.TokenSecret, TokenTTL: cfg.Auth.TokenTTL}}
	txMgr := &db.TxMgr{DB: boot.MustConnect(cfg.DataSource)}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(HandleError(), ResolveUser(auth, txMgr), RequireUser()))
	proto.RegisterUsersServer(s, UserRoutes(auth, txMgr, opts.UserRepo))

	return s
}
