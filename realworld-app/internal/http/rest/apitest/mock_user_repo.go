package apitest

import (
	"github.com/istonikula/realworld-go/realworld-app/internal/db"
	domain "github.com/istonikula/realworld-go/realworld-domain"
)

var _ db.UserRepoOps = (*MockUserRepo)(nil)

type MockUserRepo struct {
	*db.UserRepo

	MockCreate           func(domain.ValidUserRegistration) (*domain.User, error)
	MockExistsByUsername func(username string) (bool, error)
	MockExistsByEmail    func(email string) (bool, error)
	MockFindById         func(id domain.UserId) (*domain.User, error)
	MockFindByEmail      func(email string) (*domain.UserAndPassword, error)
}

func (r *MockUserRepo) Create(reg domain.ValidUserRegistration) (*domain.User, error) {
	if r.MockCreate != nil {
		return r.MockCreate(reg)
	}
	return r.UserRepo.Create(reg)
}

func (r *MockUserRepo) ExistsByUsername(username string) (bool, error) {
	if r.MockExistsByUsername != nil {
		return r.MockExistsByUsername(username)
	}
	return r.UserRepo.ExistsByUsername(username)
}

func (r *MockUserRepo) ExistsByEmail(email string) (bool, error) {
	if r.MockExistsByEmail != nil {
		return r.MockExistsByEmail(email)
	}
	return r.UserRepo.ExistsByEmail(email)
}

func (r *MockUserRepo) FindById(id domain.UserId) (*domain.User, error) {
	if r.MockFindById != nil {
		return r.MockFindById(id)
	}
	return r.UserRepo.FindById(id)
}

func (r *MockUserRepo) FindByEmail(email string) (*domain.UserAndPassword, error) {
	if r.MockFindByEmail != nil {
		return r.MockFindByEmail(email)
	}
	return r.UserRepo.FindByEmail(email)
}
