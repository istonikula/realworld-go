package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/istonikula/realworld-go/realworld-app/internal/db"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/grpc/proto"
	domain "github.com/istonikula/realworld-go/realworld-domain"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errUnauthorized = status.Error(codes.Unauthenticated, "unauthorized")
)

func ResolveUser(auth *domain.Auth, txMgr *db.TxMgr) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		token := auth.ValidateToken(TokenFromContext(ctx))
		if token == nil {
			return handler(ctx, req)
		}

		err = txMgr.Read(func(tx *sqlx.Tx) error {
			repo := &db.UserRepo{Tx: tx}

			user, err := repo.FindById(token.Id)
			if err != nil {
				return err
			}

			if user != nil {
				ctx = NewContextWithUser(ctx, user)
			}
			return nil
		})

		if err != nil {
			slog.Info(fmt.Errorf("ResolveUser: %w", err).Error())
		}

		return handler(ctx, req)
	}
}

func RequireUser() grpc.UnaryServerInterceptor {
	notRequired := []string{
		proto.Users_Login_FullMethodName,
		proto.Users_RegisterUser_FullMethodName,
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

		if slices.Contains(notRequired, info.FullMethod) {
			return handler(ctx, req)
		}

		if UserFromContext(ctx) == nil {
			return nil, errUnauthorized
		}

		return handler(ctx, req)
	}
}

func HandleError() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err == nil {
			return
		}

		var regErr *domain.UserRegistrationError
		var vErrs validation.Errors
		switch {
		case errors.As(err, &regErr):
			err = status.Error(codes.AlreadyExists, err.Error())
		case errors.As(err, &vErrs):
			err = status.Error(codes.InvalidArgument, err.Error())
		}

		return
	}
}
