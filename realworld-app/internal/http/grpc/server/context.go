package server

import (
	"context"
	"strings"

	domain "github.com/istonikula/realworld-go/realworld-domain"
	"google.golang.org/grpc/metadata"
)

type userKey struct{}

func UserFromContext(ctx context.Context) *domain.User {
	user, _ := ctx.Value(userKey{}).(*domain.User)
	return user
}

func NewContextWithUser(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

func TokenFromContext(ctx context.Context) string {
	const prefix = "Token "

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	auth := md["authorization"]
	if len(auth) < 1 {
		return ""
	}

	return strings.TrimPrefix(auth[0], prefix)
}
