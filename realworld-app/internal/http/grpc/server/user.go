package server

import (
	"context"
	"fmt"
	"log/slog"

	appDb "github.com/istonikula/realworld-go/realworld-app/internal/db"
	"github.com/istonikula/realworld-go/realworld-app/internal/http/grpc/proto"
	domain "github.com/istonikula/realworld-go/realworld-domain"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/emptypb"
)

func UserRoutes(auth *domain.Auth, txMgr *appDb.TxMgr, userRepo appDb.NewUserRepo) proto.UsersServer {
	return &server{
		auth:     auth,
		txMgr:    txMgr,
		userRepo: userRepo,
	}
}

type server struct {
	proto.UnimplementedUsersServer
	auth     *domain.Auth
	txMgr    *appDb.TxMgr
	userRepo appDb.NewUserRepo
}

func (s *server) CurrentUser(ctx context.Context, _ *emptypb.Empty) (*proto.UserResponse, error) {
	return UserResponse.fromDomain(UserFromContext(ctx)), nil
}

func (s *server) RegisterUser(ctx context.Context, dto *proto.UserRegistration) (*proto.UserResponse, error) {
	err := dto.Validate()
	if err != nil {
		return nil, err
	}

	var u *domain.User
	err = s.txMgr.Write(func(tx *sqlx.Tx) error {
		repo := s.userRepo(tx)

		validateUserSrv := domain.ValidateUserService{
			Auth:             *s.auth,
			ExistsByUsername: repo.ExistsByUsername,
			ExistsByEmail:    repo.ExistsByEmail,
		}

		u, err = domain.RegisterUserUseCase{
			Validate:   validateUserSrv.ValidateUser,
			CreateUser: repo.Create,
		}.Run(domain.UserRegistration{
			Username: dto.Username,
			Email:    dto.Email,
			Password: dto.Password,
		})
		return err
	})

	if err != nil {
		return nil, err
	}

	return UserResponse.fromDomain(u), nil
}

func (s *server) Login(ctx context.Context, dto *proto.LoginRequest) (*proto.UserResponse, error) {
	err := dto.Validate()
	if err != nil {
		return nil, err
	}

	var u *domain.User
	err = s.txMgr.Read(func(tx *sqlx.Tx) error {
		repo := &appDb.UserRepo{Tx: tx}

		u, err = domain.LoginUserUseCase{
			Auth:    s.auth,
			GetUser: repo.FindByEmail,
		}.Run(domain.Login{
			Email:    dto.Email,
			Password: dto.Password,
		})
		return err
	})

	if err != nil {
		slog.Info(fmt.Errorf("login: %w", err).Error())
		return nil, err
	}

	return UserResponse.fromDomain(u), nil
}

var UserResponse = userResponse{}

type userResponse struct{}

func (userResponse) fromDomain(domain *domain.User) *proto.UserResponse {
	return &proto.UserResponse{User: &proto.User{
		Email:    domain.Email,
		Token:    domain.Token,
		Username: domain.Username,
		Bio:      domain.Bio,
		Image:    domain.Image,
	}}
}
