package crypto

type Argon2Config struct {
	saltLength uint32
	iterations uint32
	memory     uint32
	threads    uint8
	keyLength  uint32
}

func NewDefaultArgon2Config() Argon2Config {
	return Argon2Config{
		saltLength: 32,
		iterations: 1,
		memory:     64 * 1024,
		threads:    4,
		keyLength:  32,
	}
}
