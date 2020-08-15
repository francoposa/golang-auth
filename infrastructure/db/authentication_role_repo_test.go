package db

import (
	"golang-auth/usecases/repos"
	"golang-auth/usecases/resources"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPGAuthNRoleRepo(t *testing.T) {
	assertions := assert.New(t)

	sqlxDB, closeDB := SetUpDB(t)
	defer closeDB(t, sqlxDB)
	authNRoleRepo, _ := SetUpAuthNRoleRepo(t, sqlxDB)

	role := resources.NewAuthNRole("staff")

	t.Run("create role", func(t *testing.T) {
		createdRole, _ := authNRoleRepo.Create(role)
		assertAuthNRole(assertions, role, createdRole)
	})

	t.Run("create already existing role - error", func(t *testing.T) {
		alreadyCreatedRole, err := authNRoleRepo.Create(role)
		assertions.Nil(alreadyCreatedRole)
		assertions.IsType(repos.AuthNRoleNameAlreadyExistsError{}, err)
	})

	t.Run("get role", func(t *testing.T) {
		retrievedRole, err := authNRoleRepo.GetByName(role.RoleName)
		if err != nil {
			t.Error(err)
		}
		assertAuthNRole(assertions, role, retrievedRole)
	})

	t.Run("get nonexistent role - error", func(t *testing.T) {
		nonexistentRole, err := authNRoleRepo.GetByName("xxx")
		assertions.Nil(nonexistentRole, "expected nil struct, got: %q", nonexistentRole)
		assertions.IsType(repos.AuthNRoleNotFoundError{}, err)
	})
}

func assertAuthNRole(a *assert.Assertions, want, got *resources.AuthNRole) {
	a.Equal(
		want, got, "expected equivalent structs, want: %q, got: %q", want, got,
	)
}
