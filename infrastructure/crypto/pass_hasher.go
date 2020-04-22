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

// Verify checks a candidate password against a known hash using the argon2 hashing algorithm
func (a Argon2PassHasher) Verify(password, hash string) (bool, error) {
	// hash is our known argon2 hash in the standard encoded representation
	knownEncodedHashParts := strings.Split(hash, "$")

	// Extract metadata from the known encoded argon2 hash
	var knownHashMemory uint32
	var knownHashIterations uint32
	var knownHashThreads uint8
	_, err := fmt.Sscanf(
		knownEncodedHashParts[3],
		"m=%d,t=%d,p=%d",
		&knownHashMemory,
		&knownHashIterations,
		&knownHashThreads,
	)
	if err != nil {
		return false, err
	}

	// Base64 decode the salt from the known encoded argon2 hash
	knownB64Salt := knownEncodedHashParts[4]
	knownSalt, err := base64.RawStdEncoding.DecodeString(knownB64Salt)
	if err != nil {
		return false, err
	}

	// Base64 decode the password section from the known encoded argon2 hash
	knownB64HashedPassword := knownEncodedHashParts[5]
	knownHashedPassword, err := base64.RawStdEncoding.DecodeString(knownB64HashedPassword)
	if err != nil {
		return false, err
	}

	keyLength := uint32(len(knownHashedPassword))

	candidateHashedPassword := argon2.IDKey(
		[]byte(password),
		knownSalt,
		knownHashIterations,
		knownHashMemory,
		knownHashThreads,
		keyLength,
	)

	if subtle.ConstantTimeCompare(knownHashedPassword, candidateHashedPassword) == 1 {
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
