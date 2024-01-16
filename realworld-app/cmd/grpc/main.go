package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"

	"github.com/istonikula/realworld-go/realworld-app/internal/boot"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/grpc/server"
	"google.golang.org/grpc"
)

func main() {
	cfg := boot.ReadConfig("../../config.yml")

	boot.Migrate("../../db", &cfg.DataSource)

	if err := run(server.Router(cfg)); err != nil {
		log.Fatal(err)
	}
}

func run(s *grpc.Server) error {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("listening and serving gRPC on %v", l.Addr()))
	return s.Serve(l)
}
