package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"

	"github.com/istonikula/realworld-go/realworld-app/internal/boot"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/grpc/server"
)

func main() {
	cfg := boot.ReadConfig("../../config.yml")

	boot.Migrate("../../db", &cfg.DataSource)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	slog.Info(fmt.Sprintf("server listening at %v", lis.Addr()))

	if err = server.Router(cfg).Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
