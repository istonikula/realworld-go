package server

import (
	"github.com/istonikula/realworld-go/realworld-app/internal/boot"
	"github.com/istonikula/realworld-go/realworld-app/internal/config"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/grpc/proto"
	"google.golang.org/grpc"
)

func Router(cfg *config.Config, opts ...boot.RepoOpt) *grpc.Server {
	s := grpc.NewServer()
	proto.RegisterUsersServer(s, UserRoutes())
	return s
}
