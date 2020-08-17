package db

import (
	"golang-auth/usecases/repos"
	"golang-auth/usecases/resources"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPGResourceRepo(t *testing.T) {
	assertions := assert.New(t)

	sqlxDB, closeDB := SetUpDB(t)
	defer closeDB(t, sqlxDB)
	resourceRepo, _ := SetUpAuthZResourceRepo(t, sqlxDB)

	resource := resources.NewAuthZResourceType("payment_method", "ExampleCom User Payment entity")

	t.Run("create resource", func(t *testing.T) {
		createdResource, _ := resourceRepo.Create(resource)
		assertResource(assertions, resource, createdResource)
	})

	t.Run("create already existing resource - error", func(t *testing.T) {
		alreadyCreatedResource, err := resourceRepo.Create(resource)
		assertions.Nil(alreadyCreatedResource)
		assertions.IsType(&repos.ResourceNameAlreadyExistsError{}, err)
	})

	t.Run("get resource", func(t *testing.T) {
		retrievedResource, err := resourceRepo.Get(resource.Name)
		if err != nil {
			t.Error(err)
		}
		assertResource(assertions, resource, retrievedResource)
	})

	t.Run("get nonexistent resource - error", func(t *testing.T) {
		nonexistentResource, err := resourceRepo.Get("xxx")
		assertions.Nil(nonexistentResource, "expected nil struct, got: %q", nonexistentResource)
		assertions.IsType(&repos.ResourceNameNotFoundError{}, err)
	})
}

func assertResource(a *assert.Assertions, want, got *resources.AuthZResourceType) {
	a.Equal(
		want, got, "expected equivalent structs, want: %q, got: %q", want, got,
	)
}
