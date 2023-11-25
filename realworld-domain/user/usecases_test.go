package user_test

import (
	"testing"

	"github.com/istonikula/realworld-go/realworld-domain/internal/test-support/fixture"
	"github.com/istonikula/realworld-go/realworld-domain/internal/test-support/stub"
	rwt "github.com/istonikula/realworld-go/realworld-testing"

	. "github.com/istonikula/realworld-go/realworld-domain/user"
)

var userFactory = fixture.UserFactory{Auth: stub.UserStub.Auth}
var jane = fixture.TestUser(userFactory.NewUser("jane"))

func TestUserRegister(t *testing.T) {
	reg := jane.Registration()
	uc := RegisterUseCase{stub.UserStub.CreateUser}
	act := uc.Run(&reg)

	rwt.Equals(t, act.Email, jane.Email)
	rwt.Equals(t, act.Username, jane.Username)
}
