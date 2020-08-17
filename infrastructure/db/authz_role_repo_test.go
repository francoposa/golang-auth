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
	authNRoleRepo, _ := SetUpAuthZRoleRepo(t, sqlxDB)

	role := resources.NewAuthZRole("staff")

	t.Run("create role", func(t *testing.T) {
		createdRole, _ := authNRoleRepo.Create(role)
		assertAuthZRole(assertions, role, createdRole)
	})

	t.Run("create already existing role - error", func(t *testing.T) {
		alreadyCreatedRole, err := authNRoleRepo.Create(role)
		assertions.Nil(alreadyCreatedRole)
		assertions.IsType(repos.AuthZRoleNameAlreadyExistsError{}, err)
	})

	t.Run("get role", func(t *testing.T) {
		retrievedRole, err := authNRoleRepo.GetByName(role.Name)
		if err != nil {
			t.Error(err)
		}
		assertAuthZRole(assertions, role, retrievedRole)
	})

	t.Run("get nonexistent role - error", func(t *testing.T) {
		nonexistentRole, err := authNRoleRepo.GetByName("xxx")
		assertions.Nil(nonexistentRole, "expected nil struct, got: %q", nonexistentRole)
		assertions.IsType(repos.AuthZRoleNameNotFoundError{}, err)
	})
}

func assertAuthZRole(a *assert.Assertions, want, got *resources.AuthZRole) {
	a.Equal(
		want, got, "expected equivalent structs, want: %q, got: %q", want, got,
	)
}
