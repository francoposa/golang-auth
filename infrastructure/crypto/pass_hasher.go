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
	// hash is our known argon2 hash in the standard encoded representation
	knownEncodedHashParts := strings.Split(hash, "$")

	// Get metadata
	var knownHashMemory uint32
	var knownEncodedHashIterations uint32
	var knownEncodedHashThreads uint8
	_, err := fmt.Sscanf(
		knownEncodedHashParts[3],
		"m=%d,t=%d,p=%d",
		&knownHashMemory,
		&knownEncodedHashIterations,
		&knownEncodedHashThreads,
	)
	if err != nil {
		return false, err
	}

	b64Salt := knownEncodedHashParts[4]
	knownDecodedHashSalt, err := base64.RawStdEncoding.DecodeString(b64Salt)
	if err != nil {
		return false, err
	}

	b64Hash := knownEncodedHashParts[5]
	knownDecodedHashKey, err := base64.RawStdEncoding.DecodeString(b64Hash)
	if err != nil {
		return false, err
	}

	keyLength := uint32(len(knownDecodedHashKey))

	candidateHash := argon2.IDKey(
		[]byte(password),
		knownDecodedHashSalt,
		knownEncodedHashIterations,
		knownHashMemory,
		knownEncodedHashThreads,
		keyLength,
	)

	if subtle.ConstantTimeCompare(knownDecodedHashKey, candidateHash) == 1 {
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
