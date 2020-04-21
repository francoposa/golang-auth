package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2PassHasher struct {
	config Argon2Config
}

func NewDefaultArgon2PassHasher() Argon2PassHasher {
	return Argon2PassHasher{config: NewDefaultArgon2Config()}
}

func (a Argon2PassHasher) Hash(password string) (string, error) {
	salt, err := generateRandomSalt(a.config.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		a.config.iterations,
		a.config.memory,
		a.config.threads,
		a.config.keyLength,
	)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		a.config.memory,
		a.config.iterations,
		a.config.threads,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

func (a Argon2PassHasher) Verify(password, hash string) (bool, error) {

	knownHashparts := strings.Split(hash, "$")

	var memory uint32
	var iterations uint32
	var threads uint8
	_, err := fmt.Sscanf(
		knownHashparts[3],
		"m=%d,t=%d,p=%d",
		&memory,
		&iterations,
		&threads,
	)
	if err != nil {
		return false, err
	}

	knownHashSalt, err := base64.RawStdEncoding.DecodeString(knownHashparts[4])
	if err != nil {
		return false, err
	}

	decodedKnownHash, err := base64.RawStdEncoding.DecodeString(knownHashparts[5])
	if err != nil {
		return false, err
	}

	keyLength := uint32(len(decodedKnownHash))

	candidateHash := argon2.IDKey(
		[]byte(password),
		knownHashSalt,
		iterations,
		memory,
		threads,
		keyLength,
	)

	if subtle.ConstantTimeCompare(decodedKnownHash, candidateHash) == 1 {
		return true, nil
	}
	return false, nil

}

func generateRandomSalt(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
