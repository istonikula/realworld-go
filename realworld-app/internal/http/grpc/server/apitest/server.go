package apitest

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func Server(ctx context.Context, router *grpc.Server) (conn *grpc.ClientConn, cleanup func()) {
	lis := bufconn.Listen(1024 * 1024)

	go func() {
		if err := router.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	return conn, func() {
		err := lis.Close()
		if err != nil {
			log.Fatalf("error closing listener: %v", err)
		}
		router.Stop()
	}
}
