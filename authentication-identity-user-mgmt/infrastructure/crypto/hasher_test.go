package crypto

import (
	"testing"
)

func TestArgon2PassHasher(t *testing.T) {
	passHasher := NewDefaultArgon2PassHasher()

	originalPassword := "password"
	originalHash, err := passHasher.Hash(originalPassword)
	if err != nil {
		t.Error(err.Error())
	}

	validPassword := "password"
	verified, err := passHasher.Verify(validPassword, originalHash)
	if err != nil {
		t.Error(err.Error())
	}

	if !verified {
		t.Errorf(
			"unexpected verification failure for candidate password `%s`; correct password `%s`",
			validPassword,
			originalPassword,
		)
	}

	invalidPassword := "Password"
	verified, err = passHasher.Verify(invalidPassword, originalHash)
	if err != nil {
		t.Error(err.Error())
	}

	if verified {
		t.Errorf(
			"unexpected verification for candidate password `%s`; correct password `%s`",
			invalidPassword,
			originalPassword,
		)
	}

}
