package server

import "github.com/istonikula/realworld-go/realworld-app/internal/http/grpc/proto"

func UserRoutes() proto.UsersServer {
	return &server{}
}

type server struct {
	proto.UnimplementedUsersServer
}
