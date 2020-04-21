package crypto

/*
Argon2 libraries across various programming languages produce hashes in a standard encoded representation, like:
$argon2id$v=19$m=65536,t=1,p=4$aB9riPbBtIDNEjfb+VoFGPY2NKov6LA60B7NzzY3Kec$LRT6Oq/QWn1E7rTjqYropIBtm8qXEB6gOITabt0Nurg

This is a "modular crypt format" (MHC), specifically the PHC string format as described here:
https://github.com/P-H-C/phc-string-format/blob/master/phc-sf-spec.md#argon2-encoding

The Python passlib docs offer a non-argon2 specific overview of Modular Crypt Format, with this noteable caveat:
"there’s no official specification document describing this format. Nor is there a central registry of
identifiers, or actual rules. The modular crypt format is more of an ad-hoc idea rather than a true standard.
https://passlib.readthedocs.io/en/stable/modular_crypt_format.html
*/

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
