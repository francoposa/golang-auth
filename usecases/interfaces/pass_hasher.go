package interfaces

type PassHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) (bool, error)
}
